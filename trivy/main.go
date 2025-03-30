package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var executionId string
var image string
var output string

func init() {
	const (
		defaultImage     = "docker.io/library/alpine"
		imageUsage       = "the image to scan"
		executionIdUsage = "the reference to use for the scan operation"
		defaultOutput    = "file:///tmp"
		outputUsage      = "the path where to store the file at. (e.g. file:///path/to/scan or s3://path/to/bucket, ...)"
	)
	defaultExecutionId := time.Now().Format(time.RFC3339)
	flag.StringVar(&image, "image", defaultImage, imageUsage)
	flag.StringVar(&image, "i", defaultImage, imageUsage+" (shorthand)")
	flag.StringVar(&executionId, "id", defaultExecutionId, executionIdUsage)
	flag.StringVar(&output, "output", defaultOutput, outputUsage+" (shorthand)")
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

func (t *TrivyScanner) Scan(ctx context.Context, image string) (*bytes.Buffer, error) {
	tCtx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	var buf bytes.Buffer

	cmd := exec.CommandContext(tCtx, t.binaryPath, "image", image)
	cmd.Stdout = &buf
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	return &buf, err
}

func main() {
	flag.Parse()
	fmt.Println(image)

	s, err := NewTrivyScanner()
	if err != nil {
		log.Fatalf("new trivy scanner: %s\n", err)
	}

	b, err := s.Scan(context.Background(), image)
	if err != nil {
		log.Fatalf("failed to scan: %s\n", err)
	}

	sd, err := NewStorageDestination(output, executionId)
	if err != nil {
		log.Fatalf("wrong storage destination: %s\n", err)
	}

	err = Store(b, sd)
	if err != nil {
		log.Fatalf("failed to store: %s\n", err)
	}

	fmt.Println(b.String())
}
