package main

import (
	"context"
	"fmt"
	"github.com/google/uuid"
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

func generatePodConfiguration(requestId string, scanner string, image string) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: requestId,
			//Labels: map[string]string{
			//	"scanner":    scanner,
			//	"created-by": "madvsa/controlplane",
			//},
		},
		//
		//env:
		//- name: GOOGLE_APPLICATION_CREDENTIALS
		//valueFrom:
		//secretKeyRef:
		//name:
		//key:

		Spec: corev1.PodSpec{
			Volumes: []corev1.Volume{
				{Name: "gcp-credentials", VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{SecretName: "madvsa-gcs-credentials"},
				}},
			},
			Containers: []corev1.Container{
				{
					Name:  "main",
					Image: "madjlzz/madvsa-trivy:latest",
					Args:  []string{"-image", image, "-output", fmt.Sprintf("gcs://testdju/%s.trivy.txt", requestId)},
					VolumeMounts: []corev1.VolumeMount{
						{Name: "gcp-credentials", MountPath: "/gcp/credentials", ReadOnly: true},
					},
					Env: []corev1.EnvVar{
						{Name: "GOOGLE_APPLICATION_CREDENTIALS", Value: "/gcp/credentials/key.json"},
					},
					//Env: []corev1.EnvVar{
					//	{Name: "GOOGLE_APPLICATION_CREDENTIALS", Value: "/gcp/crendetials", ValueFrom: &corev1.EnvVarSource{
					//		SecretKeyRef: &corev1.SecretKeySelector{
					//			LocalObjectReference: corev1.LocalObjectReference{Name: "madvsa-gcs-credentials"},
					//			Key:                  "key.json",
					//		},
					//	}},
					//},
				},
			},
			RestartPolicy: corev1.RestartPolicyNever,
		},
	}
}

func (k *KubernetesService) Scan(ctx context.Context, scanner string, image string) error {
	id, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("cannot create uuid for pod: %w", err)
	}

	pod := generatePodConfiguration(id.String(), scanner, image)
	// TODO: let's see how we could give a user friendly value for that namespace.
	// 	Maybe config?
	createdPod, err := k.cli.CoreV1().Pods("default").Create(ctx, pod, metav1.CreateOptions{})
	if err != nil {
		panic(fmt.Sprintf("failed to create pod: %v", err))
	}
	fmt.Printf("Pod %s created\n", createdPod.Name)

	return nil
}
