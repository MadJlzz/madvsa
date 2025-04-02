package grype

import (
	"context"
	"os/exec"
)

const BinaryName = "grype"

func Cmd(ctx context.Context, image string, extraArgs ...string) *exec.Cmd {
	args := append([]string{image}, extraArgs...)
	return exec.CommandContext(ctx, BinaryName, args...)
}
