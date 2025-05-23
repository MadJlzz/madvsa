package grype

import (
	"context"
	"os/exec"
)

const BinaryName = "grype"

func Cmd(ctx context.Context, image string) *exec.Cmd {
	args := append([]string{image, "-q"})
	return exec.CommandContext(ctx, BinaryName, args...)
}
