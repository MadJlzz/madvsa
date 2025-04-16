package main

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
	"github.com/google/uuid"
)

type ImageScanner interface {
	Scan(ctx context.Context, scanner string, image string) error
}

type Scanner string

// ScannerConfigFn -> TODO: review this later, it might be probably easier to mount the configuration file of the scanner
// than passing all the args available.
type ScannerConfigFn func(requestId string, image string, cfg ScannersConfigurations) *container.Config

const (
	Trivy Scanner = "trivy"
	Grype Scanner = "grype"
)

var scannersCfg = map[Scanner]ScannerConfigFn{
	Trivy: TrivyScannerConfiguration,
	Grype: GrypeScannerConfiguration,
}

func TrivyScannerConfiguration(requestId string, image string, cfg ScannersConfigurations) *container.Config {
	// TODO: should pass --output to tell the cmd where to save the file.
	return &container.Config{
		Cmd:   strslice.StrSlice{"-image", image, "-output", fmt.Sprintf("file:///%s.trivy.txt", requestId)},
		Image: cfg.Trivy.Image,
	}
}

func GrypeScannerConfiguration(requestId string, image string, cfg ScannersConfigurations) *container.Config {
	return &container.Config{
		Cmd:   strslice.StrSlice{"-image", image, "-output", fmt.Sprintf("file:///%s.grype.txt", requestId)},
		Image: cfg.Grype.Image,
	}
}

type ContainerService struct {
	cli *client.Client
	cfg ScannersConfigurations
}

func NewContainerService(socketPath string, cfg ScannersConfigurations) *ContainerService {
	conn, err := client.NewClientWithOpts(client.WithHost(socketPath), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &ContainerService{
		cli: conn,
		cfg: cfg,
	}
}

func (cs *ContainerService) Scan(ctx context.Context, scanner string, image string) error {
	// Get the scanner options from it's name, and run that.
	ctnCfgFn, ok := scannersCfg[Scanner(scanner)]
	if !ok {
		return fmt.Errorf("unknown scanner %s", scanner)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("cannot create uuid for container: %w", err)
	}

	// Maybe it's a good idea to generate a UUID as an execution ID here for observability.
	resp, err := cs.cli.ContainerCreate(ctx, ctnCfgFn(id.String(), image, cs.cfg), nil, nil, nil, id.String())
	if err != nil {
		return err
	}

	if err = cs.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}
	return nil
}
