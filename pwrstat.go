package pwrstat

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

var DefaultPath = "/usr/sbin/pwrstat"

func Status(path string, root bool) (StatusResult, error) {
	o, err := getStatus(path, root)
	if err != nil {
		return StatusResult{}, err
	}
	return parseStatus(o)
}

func getStatus(path string, root bool) (string, error) {
	var uid int
	if root {
		uid = syscall.Getuid()
		err := syscall.Setuid(0)
		if err != nil {
			return "", fmt.Errorf("failed to setuid 0: %w", err)
		}
	}
	c := exec.Command(path, "-status")
	o, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute 'pwrstat -status': %w", err)
	}
	if root {
		err := syscall.Setuid(uid)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to step down to original uid: %v", err)
		}
	}
	return string(o), nil
}
