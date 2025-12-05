package git

import (
	"context"
	"os/exec"
	"strings"
)

func LatestTag(ctx context.Context, dir string) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "describe", "--tags", "--abbrev=0")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	result := strings.TrimSpace(string(out))
	return result, err
}
