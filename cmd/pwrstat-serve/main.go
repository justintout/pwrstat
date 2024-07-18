package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/justintout/pwrstat"
)

var (
	noroot bool
	path   string
	host   string
	port   int
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		log.Fatal(err)
	}
}

func run(ctx context.Context) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	flag.BoolVar(&noroot, "noroot", false, "execute pwrstat without elevating to root")
	flag.StringVar(&path, "path", pwrstat.DefaultPath, "path to the pwrstat executable")
	flag.StringVar(&host, "host", "0.0.0.0", "host for server to listen on")
	flag.IntVar(&port, "port", 7977, "port for server to listen on")
	flag.Parse()

	// TODO: fork a child that handles setuid to separate HTTP server
	// from a the process running as root. it's safe as-is but it makes
	// me itchy to keep escalated privs hanging around
	if !noroot {
		err := syscall.Setuid(0)
		if err != nil {
			return fmt.Errorf("failed to setuid 0: %w", err)
		}
	}

	svr := pwrstat.NewServer(pwrstat.ServerConfig{
		Host:   host,
		Port:   port,
		Path:   path,
		NoRoot: noroot,
	})

	go func() {
		log.Printf("listening on %s\n", svr.Addr)
		if err := svr.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("error listening and serving: %s\n", err)
		}
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		sdCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := svr.Shutdown(sdCtx); err != nil {
			log.Printf("error shutting down server: %v", err)
		}
	}()
	wg.Wait()
	return nil
}
