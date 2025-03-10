package main

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/strslice"
	"github.com/docker/docker/client"
)

type containerService struct {
	cli *client.Client
}

func NewContainerService(socketPath string) *containerService {
	conn, err := client.NewClientWithOpts(client.WithHost(socketPath), client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	return &containerService{
		cli: conn,
	}
}

func (cs *containerService) ScanContainer(ctx context.Context, image string) error {

	resp, err := cs.cli.ContainerCreate(ctx, &container.Config{
		Cmd:   strslice.StrSlice{"--image", image},
		Image: "madvsa-trivy:latest",
	}, nil, nil, nil, "testdju")

	if err != nil {
		panic(err)
	}

	if err = cs.cli.ContainerStart(ctx, resp.ID, container.StartOptions{}); err != nil {
		panic(err)
	}

	return nil
}
