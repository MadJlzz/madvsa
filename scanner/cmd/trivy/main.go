package main

import (
	"context"
	"flag"
	"github.com/MadJlzz/madvsa/scanner/internal/pkg/storage"
	"github.com/MadJlzz/madvsa/scanner/internal/trivy"
	"log"
	"net/url"
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

func main() {
	flag.Parse()

	ctx := context.Background()

	u, err := url.Parse(output)
	if err != nil {
		log.Fatalf("failed to parse output: %s\n", err)
	}

	s, err := trivy.New()
	if err != nil {
		log.Fatalf("new trivy scanner: %s\n", err)
	}

	storerFactory, err := storage.NewStorerFactory(ctx, u)
	if err != nil {
		log.Fatalf("failed to init storer factory: %s\n", err)
	}

	b, err := s.Scan(context.Background(), image)
	if err != nil {
		log.Fatalf("failed to scan: %s\n", err)
	}

	err = storerFactory.Store(ctx, b, u)
	if err != nil {
		log.Fatalf("failed to store: %s\n", err)
	}
}
