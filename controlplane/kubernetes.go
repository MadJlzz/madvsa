package main

import "context"

type KubernetesService struct {
}

func NewKubernetesService() *KubernetesService {
	return &KubernetesService{}
}

func (k *KubernetesService) Scan(ctx context.Context, scanner string, image string) error {
	//TODO implement me
	panic("implement me")
}
