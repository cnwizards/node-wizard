package pkg

import (
	"context"
	"os"
	"time"

	"google.golang.org/appengine/log"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/drain"
)

func DrainNode() error {
	helper := BuildDrainHelper()
	node := v1.Node{ObjectMeta: metav1.ObjectMeta{Name: "minikube"}}
	err := drain.RunCordonOrUncordon(helper, &node, true)
	if err != nil {
		log.Errorf(helper.Ctx, "Error cordoning node %s: %v", node.Name, err)
		return err
	}
	log.Infof(helper.Ctx, "Node %s cordoned", node.Name)

	err = drain.RunNodeDrain(helper, node.Name)
	if err != nil {
		log.Errorf(helper.Ctx, "Error draining node %s: %v", node.Name, err)
	}
	log.Infof(helper.Ctx, "Node %s drained", node.Name)

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
		Timeout:             100 * time.Second,
		Out:                 os.Stdout,
		ErrOut:              os.Stderr,
	}
}
