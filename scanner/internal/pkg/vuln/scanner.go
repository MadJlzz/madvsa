package vuln

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"os/exec"
	"time"
)

type ScannerCmdFn func(ctx context.Context, image string) *exec.Cmd

type Scanner struct {
	logger     *slog.Logger
	binaryPath string
	cmd        ScannerCmdFn
}

func NewScanner(logger *slog.Logger, binaryName string, fn ScannerCmdFn) (*Scanner, error) {
	path, err := exec.LookPath(binaryName)
	if err != nil {
		return nil, fmt.Errorf("look path: %w", err)
	}
	return &Scanner{
		logger:     logger,
		binaryPath: path,
		cmd:        fn,
	}, err
}

func (s *Scanner) Scan(ctx context.Context, image string) (*bytes.Buffer, error) {
	tCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var buf bytes.Buffer
	cmd := s.cmd(tCtx, image)
	cmd.Stdout = &buf

	s.logger.Info("scanning image and storing results in buffer", "scanner", s.binaryPath, "image", image, "args", cmd.Args)

	err := cmd.Run()
	return &buf, err
}
