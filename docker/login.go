package docker

import (
	"context"
	"os"
	"os/exec"
	"strings"
)

func Login(ctx context.Context, registryUrl, username, passwordEnv string) error {
	pass := os.Getenv(passwordEnv)
	cmd := exec.CommandContext(ctx, "docker", "login", registryUrl, "-u", username, "--password-stdin")
	cmd.Stdin = strings.NewReader(pass)

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
