package config

import (
	"os"
	"os/exec"
	"strings"
)

func pass(args ...string) *exec.Cmd {
	cmd := exec.Command("pass", args...)
	cmd.Stderr = os.Stderr

	return cmd
}

func getSecret(name string) (string, error) {
	cmd := pass("show", name)

	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}
