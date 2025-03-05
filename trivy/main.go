package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var image string

func init() {
	const (
		defaultImage = "docker.io/library/alpine"
		usage        = "the image to scan"
	)
	flag.StringVar(&image, "image", defaultImage, usage)
	flag.StringVar(&image, "i", defaultImage, usage+" (shorthand)")
}

type TrivyScanner struct {
	binaryPath string
}

func NewTrivyScanner() (*TrivyScanner, error) {
	path, err := exec.LookPath("trivy")
	if err != nil {
		return nil, fmt.Errorf("look path: %w", err)
	}
	return &TrivyScanner{
		binaryPath: path,
	}, err
}

func (t *TrivyScanner) Scan(ctx context.Context, image string) error {
	tCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(tCtx, t.binaryPath, "image", image)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func main() {
	flag.Parse()
	fmt.Println(image)

	s, err := NewTrivyScanner()
	if err != nil {
		log.Fatalf("new trivy scanner: %s\n", err)
	}

	err = s.Scan(context.Background(), image)
	if err != nil {
		log.Printf("scan: %s\n", err)
	}
}
