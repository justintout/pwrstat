package pwrstat

import (
	"fmt"
	"os/exec"
)

var DefaultPath = "/usr/sbin/pwrstat"

func Status(path string) (StatusResult, error) {
	o, err := getStatus(path)
	if err != nil {
		return StatusResult{}, err
	}
	return parseStatus(o)
}

func getStatus(path string) (string, error) {
	c := exec.Command(path, "-status")
	o, err := c.Output()
	if err != nil {
		return "", fmt.Errorf("failed to execute 'pwrstat -status': %w", err)
	}
	return string(o), nil
}
