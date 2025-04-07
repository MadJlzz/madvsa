package main

import (
	"context"
	"github.com/MadJlzz/madvsa/scanner/internal/pkg/cmd"
	"github.com/MadJlzz/madvsa/scanner/internal/pkg/vuln"
	"github.com/MadJlzz/madvsa/scanner/internal/trivy"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	s, err := vuln.NewScanner(l, trivy.BinaryName, trivy.Cmd)
	if err != nil {
		l.Error("new trivy scanner", "err", err)
		os.Exit(1)
	}

	trivyCmd := cmd.NewCommand(l, s)
	if err = trivyCmd.Execute(ctx); err != nil {
		l.Error("trivy scanner", "err", err)
		os.Exit(1)
	}
}
