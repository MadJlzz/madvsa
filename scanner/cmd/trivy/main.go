package main

import (
	"context"
	"github.com/MadJlzz/madvsa/scanner/internal/pkg/cmd"
	"github.com/MadJlzz/madvsa/scanner/internal/trivy"
	"log"
)

func main() {
	ctx := context.Background()

	s, err := trivy.New()
	if err != nil {
		log.Fatalf("new trivy scanner: %s\n", err)
	}

	trivyCmd := cmd.NewCommand(s)
	if err = trivyCmd.Execute(ctx); err != nil {
		log.Fatalf("trivy scanner: %s\n", err)
	}
}
