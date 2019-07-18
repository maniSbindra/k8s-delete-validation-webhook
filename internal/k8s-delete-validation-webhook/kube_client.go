package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func KubeClientSet(inCluster bool) (*kubernetes.Clientset, error) {

	var config *rest.Config

	if inCluster {
		c, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		config = c
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}
