package trivy

import (
	"context"
	"os/exec"
)

const BinaryName = "trivy"

func Cmd(ctx context.Context, extraArg ...string) *exec.Cmd {
	args := append([]string{"image"}, extraArg...)
	return exec.CommandContext(ctx, BinaryName, args...)
}
