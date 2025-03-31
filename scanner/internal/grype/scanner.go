package grype

import (
	"context"
	"os/exec"
)

const BinaryName = "grype"

func Cmd(ctx context.Context, extraArg ...string) *exec.Cmd {
	return exec.CommandContext(ctx, BinaryName, extraArg...)
}
