package cmd

import (
	"context"
	"flag"
	"fmt"
	"github.com/MadJlzz/madvsa/scanner/internal/pkg/storage"
	"github.com/MadJlzz/madvsa/scanner/internal/pkg/vuln"
	"log/slog"
	"net/url"
)

const (
	defaultImage  = "docker.io/library/alpine"
	imageUsage    = "the image to scan"
	defaultOutput = "file:///tmp"
	outputUsage   = "the path where to store the file at. (e.g. file:///path/to/scan or s3://path/to/bucket, ...)"
)

type args struct {
	image  *string
	output *string
}

type Command struct {
	args    *args
	logger  *slog.Logger
	scanner *vuln.Scanner
}

func NewCommand(logger *slog.Logger, scanner *vuln.Scanner) *Command {
	// TODO: something is fully broken here. The values passed from the flag library does not work.
	return &Command{
		scanner: scanner,
		logger:  logger,
		args: &args{
			image:  flag.String("image", defaultImage, imageUsage),
			output: flag.String("output", defaultOutput, outputUsage),
		}}
}

func (c *Command) Execute(ctx context.Context) error {
	flag.Parse()

	c.logger.Info("OK?", "parsed", flag.Parsed(), "args", *c.args)
	c.logger.Info("hello", "output", *c.args.output, "img", *c.args.image)

	u, err := url.Parse(*c.args.output)
	if err != nil {
		return fmt.Errorf("failed to parse output: %s\n", err)
	}

	c.logger.Info("test", "url", u)

	storerFactory, err := storage.NewStorerFactory(ctx, c.logger, u)
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
