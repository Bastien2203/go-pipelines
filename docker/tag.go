package docker

import (
	"context"
	"os"
	"os/exec"
)

func Tag(ctx context.Context, baseTag, fullTag string) error {
	cmd := exec.CommandContext(ctx, "docker", "tag", baseTag, fullTag)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
