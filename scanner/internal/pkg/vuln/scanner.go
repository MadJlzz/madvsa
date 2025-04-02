package vuln

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type ScannerCmdFn func(ctx context.Context, image string, extraArgs ...string) *exec.Cmd

type Scanner struct {
	binaryPath string
	cmd        ScannerCmdFn
}

func NewScanner(binaryName string, fn ScannerCmdFn) (*Scanner, error) {
	path, err := exec.LookPath(binaryName)
	if err != nil {
		return nil, fmt.Errorf("look path: %w", err)
	}
	return &Scanner{
		binaryPath: path,
		cmd:        fn,
	}, err
}

func (s *Scanner) Scan(ctx context.Context, image string, extraArg ...string) (*bytes.Buffer, error) {
	tCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var buf bytes.Buffer

	cmd := s.cmd(tCtx, image, extraArg...)
	cmd.Stdout = &buf
	//cmd.Stderr = os.Stderr

	err := cmd.Run()
	return &buf, err
}
