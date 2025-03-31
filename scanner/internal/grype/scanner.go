package grype

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type Scanner struct {
	binaryPath string
}

func New() (*Scanner, error) {
	path, err := exec.LookPath("grype")
	if err != nil {
		return nil, fmt.Errorf("look path: %w", err)
	}
	return &Scanner{
		binaryPath: path,
	}, err
}

func (s *Scanner) Scan(ctx context.Context, image string) (*bytes.Buffer, error) {
	tCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var buf bytes.Buffer

	cmd := exec.CommandContext(tCtx, s.binaryPath, image)
	cmd.Stdout = &buf
	// TODO: if we set Stderr, is doesnt work anymore
	//cmd.Stderr = os.Stderr

	err := cmd.Run()
	return &buf, err
}
