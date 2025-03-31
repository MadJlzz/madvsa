package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/MadJlzz/madvsa/scanner/internal/pkg/storage"
	"github.com/MadJlzz/madvsa/scanner/internal/pkg/vuln"
	"net/url"
	"time"
)

const (
	defaultImage     = "docker.io/library/alpine"
	imageUsage       = "the image to scan"
	executionIdUsage = "the reference to use for the scan operation"
	defaultOutput    = "file:///tmp"
	outputUsage      = "the path where to store the file at. (e.g. file:///path/to/scan or s3://path/to/bucket, ...)"
)

type args struct {
	executionId *string
	image       *string
	output      *string
}

type Command struct {
	args    args
	scanner *vuln.Scanner
}

func NewCommand(scanner *vuln.Scanner) *Command {
	return &Command{scanner: scanner, args: args{
		executionId: flag.String("id", time.Now().Format(time.RFC3339), executionIdUsage),
		image:       flag.String("image", defaultImage, imageUsage),
		output:      flag.String("output", defaultOutput, outputUsage),
	}}
}

func (c *Command) Execute(ctx context.Context) error {
	flag.Parse()

	u, err := url.Parse(*c.args.output)
	if err != nil {
		return fmt.Errorf("failed to parse output: %s\n", err)
	}

	storerFactory, err := storage.NewStorerFactory(ctx, u)
	if err != nil {
		return fmt.Errorf("failed to init storer factory: %s\n", err)
	}

	b, err := c.scanner.Scan(ctx, *c.args.image)
	if err != nil {
		return fmt.Errorf("failed to scan: %s\n", err)
	}

	err = storerFactory.Store(ctx, b, u)
	if err != nil {
		return fmt.Errorf("failed to store: %s\n", err)
	}

	return nil
}
