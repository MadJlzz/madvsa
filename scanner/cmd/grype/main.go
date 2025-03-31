package main

import (
	"context"
	"github.com/MadJlzz/madvsa/scanner/internal/grype"
	"github.com/MadJlzz/madvsa/scanner/internal/pkg/cmd"
	"log"
)

func main() {
	ctx := context.Background()

	s, err := grype.New()
	if err != nil {
		log.Fatalf("new grype scanner: %s\n", err)
	}

	grypeCmd := cmd.NewCommand(s)
	if err = grypeCmd.Execute(ctx); err != nil {
		log.Fatalf("grype scanner: %s\n", err)
	}
}
