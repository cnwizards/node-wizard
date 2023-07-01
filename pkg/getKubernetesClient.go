package pkg

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetKubernetesClient() (*kubernetes.Clientset, error) {
	var config, err = GetClusterConfig()
	if err != nil {
		log.Fatalf("Error getting cluster config: %v", err)
		return nil, err
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Debugf("Error getting kubernetes client: %v", err)
		return nil, err
	}
	return client, nil
}

func GetClusterConfig() (*rest.Config, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Panicf("Error getting cluster config: %v", err)
		return nil, err
	}
	return config, nil
}
