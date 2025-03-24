package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

type ImageScanner interface {
	Scan(ctx context.Context, scanner string, image string) error
}

type Scanner string

type ScannerConfigFn func(img string) *container.Config

const (
	Trivy Scanner = "trivy"
)

var scannersCfg = map[Scanner]ScannerConfigFn{
	Trivy: TrivyScannerConfiguration,
}

func TrivyScannerConfiguration(img string) *container.Config {
	return &container.Config{
		Cmd:   strslice.StrSlice{"--image", img},
		Image: "madvsa/trivy:latest",
	}
}

type ContainerService struct {
	cli *client.Client
}

func NewContainerService(socketPath string) *ContainerService {
	conn, err := client.NewClientWithOpts(client.WithHost(socketPath), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &ContainerService{
		cli: conn,
	}
}

func (cs *ContainerService) Scan(ctx context.Context, scanner string, image string) error {
	// Get the scanner options from it's name, and run that.
	ctnCfgFn, ok := scannersCfg[Scanner(scanner)]
	if !ok {
		return fmt.Errorf("unknown scanner %s", scanner)
	}

	// Maybe it's a good idea to generate an UUID as an execution ID here for observability.
	resp, err := cs.cli.ContainerCreate(ctx, ctnCfgFn(image), nil, nil, nil, "")
	if err != nil {
		return err
	}

	if err = cs.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}
	return nil
}
