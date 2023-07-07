package utils

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/drain"
)

type Node struct {
	Info *corev1.Node
}

func (n Node) DrainNode() error {
	helper := BuildDrainHelper()
	log.Debugf("Drain helper: %+v", helper)
	node := corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: n.Info.Name}}
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

func (n Node) UncordonNode() error {
	clientset, err := GetKubernetesClient()
	if err != nil {
		return err
	}

	n.Info.Spec.Unschedulable = false
	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), n.Info, metav1.UpdateOptions{})
	if err != nil {
		log.Errorf("Error uncordoning node %s: %v", n.Info.Name, err)
		return err
	}

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
		Out:                 log.StandardLogger().Writer(),
		ErrOut:              log.StandardLogger().Writer(),
	}
}
