package main

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

type ImageScanner interface {
	Scan(ctx context.Context, image string) error
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

func (cs *ContainerService) Scan(ctx context.Context, image string) error {
	resp, err := cs.cli.ContainerCreate(ctx, &container.Config{
		Cmd:   strslice.StrSlice{"--image", image},
		Image: "madvsa-trivy:latest",
	}, nil, nil, nil, "testdju")
	if err != nil {
		return err
	}

	if err = cs.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		return err
	}
	return nil
}
