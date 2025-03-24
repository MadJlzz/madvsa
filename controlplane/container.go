package main

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

type ImageScanner interface {
	Scan(ctx context.Context, scanner string, image string) error
}

type ContainerService struct {
	cli         *client.Client
	scannerOpts map[string]*container.Config
}

func NewContainerService(socketPath string) *ContainerService {
	conn, err := client.NewClientWithOpts(client.WithHost(socketPath), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	// Dirty as fuck, will clean that up later.
	sOpts := map[string]*container.Config{
		"trivy": {
			Cmd:   strslice.StrSlice{"--image"},
			Image: "madvsa/trivy:latest",
		},
	}
	return &ContainerService{
		cli:         conn,
		scannerOpts: sOpts,
	}
}

func (cs *ContainerService) Scan(ctx context.Context, scanner string, image string) error {
	// Get the scanner options from it's name, and run that.
	cCfg := cs.scannerOpts[scanner]
	cCfg.Cmd = append(cCfg.Cmd, image)

	// Maybe it's a good idea to generate an UUID as an execution ID here for observability.
	resp, err := cs.cli.ContainerCreate(ctx, cCfg, nil, nil, nil, "")
	if err != nil {
		return err
	}

	if err = cs.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}
	return nil
}
