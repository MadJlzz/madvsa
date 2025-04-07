package main

import (
	"context"
	"github.com/MadJlzz/madvsa/scanner/internal/grype"
	"github.com/MadJlzz/madvsa/scanner/internal/pkg/cmd"
	"github.com/MadJlzz/madvsa/scanner/internal/pkg/vuln"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()

	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{}))

	s, err := vuln.NewScanner(l, grype.BinaryName, grype.Cmd)
	if err != nil {
		l.Error("new grype scanner", "err", err)
		os.Exit(1)
	}

	grypeCmd := cmd.NewCommand(l, s)
	if err = grypeCmd.Execute(ctx); err != nil {
		l.Error("grype scanner", "err", err)
		os.Exit(1)
	}
}
