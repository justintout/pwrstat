package pwrstat

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
)

type ServerConfig struct {
	Host string
	Port int

	Path   string
	NoRoot bool
}

func NewServer(cfg ServerConfig) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", handleStatus(cfg.Path))
	mux.Handle("/healthz", handleHealtz())
	return &http.Server{
		Addr:    net.JoinHostPort(cfg.Host, strconv.Itoa(cfg.Port)),
		Handler: mux,
	}
}

func handleStatus(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s, err := Status(path)
		if err != nil {
			log.Printf("error getting status: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Printf("error encoding status: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func handleHealtz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
