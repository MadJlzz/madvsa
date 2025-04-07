package trivy

import (
	"context"
	"os/exec"
)

const BinaryName = "trivy"

func Cmd(ctx context.Context, image string) *exec.Cmd {
	args := append([]string{"image", image})
	return exec.CommandContext(ctx, BinaryName, args...)
}
