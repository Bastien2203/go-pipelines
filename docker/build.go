package docker

import (
	"context"
	"os"
	"os/exec"
)

func Build(ctx context.Context, baseTag, dir string) error {
	cmd := exec.CommandContext(ctx, "docker", "build", "-t", baseTag, ".")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
