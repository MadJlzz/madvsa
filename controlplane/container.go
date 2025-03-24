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

func GetScannerConfiguration(scanner string, img string) (*container.Config, error) {
	// Still dirty, I do not have a better idea yet.
	switch scanner {
	case "trivy":
		return &container.Config{
			Cmd:   strslice.StrSlice{"--image", img},
			Image: "madvsa/trivy:latest",
		}, nil
	default:
		return nil, fmt.Errorf("unknown scanner %s", scanner)
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
	ctnCfg, err := GetScannerConfiguration(scanner, image)
	if err != nil {
		return err
	}

	// Maybe it's a good idea to generate an UUID as an execution ID here for observability.
	resp, err := cs.cli.ContainerCreate(ctx, ctnCfg, nil, nil, nil, "")
	if err != nil {
		return err
	}

	if err = cs.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}
	return nil
}
