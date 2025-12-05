package git

import (
	"context"
	"os"
	"os/exec"
)

func Clone(ctx context.Context, branch, url, dir string) error {
	cmd := exec.CommandContext(ctx, "git", "clone", "-b", branch, url, dir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
