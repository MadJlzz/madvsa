package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type KubernetesService struct {
	cli *kubernetes.Clientset
}

func NewKubernetesService() *KubernetesService {
	// Use the in-cluster config (assumes service account token is mounted)
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
		//return nil, fmt.Errorf("failed to load in-cluster config: %w", err)
	}
	// Create clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
		//return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}
	return &KubernetesService{
		cli: clientSet,
	}
}

func (k *KubernetesService) Scan(ctx context.Context, scanner string, image string) error {
	// Define pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "custom-command-pod",
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "main",
					Image:   "busybox",
					Command: []string{"/bin/sh", "-c", "echo 'Hello World!'"},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
	// Create the pod in the default namespace
	createdPod, err := k.cli.CoreV1().Pods("default").Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		panic(fmt.Sprintf("failed to create pod: %v", err))
	}
	fmt.Printf("Pod %s created\n", createdPod.Name)

	return nil
}
