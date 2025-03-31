package vuln

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"
)

type ScannerCmdFn func(ctx context.Context, extraArg ...string) *exec.Cmd

type Scanner struct {
	binaryPath string
	cmd        ScannerCmdFn
}

func NewScanner(binaryName string, cmd ScannerCmdFn) (*Scanner, error) {
	path, err := exec.LookPath(binaryName)
	if err != nil {
		return nil, fmt.Errorf("look path: %w", err)
	}
	return &Scanner{
		binaryPath: path,
		cmd:        cmd,
	}, err
}

func (s *Scanner) Scan(ctx context.Context, image string) (*bytes.Buffer, error) {
	tCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var buf bytes.Buffer

	cmd := s.cmd(tCtx, image)
	cmd.Stdout = &buf
	//cmd.Stderr = os.Stderr

	err := cmd.Run()
	return &buf, err
}
