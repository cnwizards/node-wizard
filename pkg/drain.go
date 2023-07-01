package pkg

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/drain"
)

func DrainNode() error {
	helper := BuildDrainHelper()
	log.Debugf("Drain helper: %+v", helper)
	node := v1.Node{ObjectMeta: metav1.ObjectMeta{Name: "minikube"}}
	log.Debugf("Node: %+v", node)

	err := drain.RunCordonOrUncordon(helper, &node, true)
	if err != nil {
		log.Errorf("Error cordoning node %s: %v", node.Name, err)
		return err
	}
	log.Infof("Node %s cordoned", node.Name)

	err = drain.RunNodeDrain(helper, node.Name)
	if err != nil {
		log.Errorf("Error draining node %s: %v", node.Name, err)
		return err
	}
	log.Infof("Node %s drained", node.Name)

	return nil
}

func BuildDrainHelper() *drain.Helper {
	clientset, err := GetKubernetesClient()
	if err != nil {
		panic(err.Error())
	}

	return &drain.Helper{
		Ctx:                 context.Background(),
		Client:              clientset,
		Force:               true,
		GracePeriodSeconds:  0,
		IgnoreAllDaemonSets: true,
		DeleteEmptyDirData:  true,
		Timeout:             100 * time.Second,
		Out:                 log.StandardLogger().Out,
		ErrOut:              log.StandardLogger().Out,
	}
}
