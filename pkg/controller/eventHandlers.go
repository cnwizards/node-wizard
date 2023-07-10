package controller

import (
	"net/http"

	"github.com/cnwizards/node-wizard/pkg/utils"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	IgnoreLabel string = "node-wizard/ignore"
)

func OnAdd(obj interface{}) error {
	log.Debugf("Received add event: %v", obj.(metav1.Object).GetName())
	return nil
}

func OnDelete(obj interface{}) error {
	log.Debugf("Received delete event: %v", obj.(metav1.Object).GetName())
	return nil
}

func OnUpdate(node *corev1.Node) error {
	nodeStruct := utils.Node{Info: node}
	for _, condition := range node.Status.Conditions {
		log.Debugf("Node %s condition: %+v", node.Name, condition)

		// If the node has the label ignore it.
		if node.Labels[IgnoreLabel] == "true" {
			log.Debugf("Node %s is ignored.", node.Name)
			return nil
		}

		if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue && !node.Spec.Unschedulable {
			// If the node is ready, and schedulable: do nothing
			log.Debugf("Node %s is ready and schedulable, skipping.", node.Name)
			return nil
		} else if condition.Type == corev1.NodeReady && condition.Status != corev1.ConditionTrue && !node.Spec.Unschedulable {
			// If the node is notReady, and schedulable: drain it and evict sts pods
			log.Infof("Node %s is not ready and schedulable, draining the node.", node.Name)
			err := nodeStruct.DrainNode()
			if err != nil {
				log.Errorf("Error draining the node: %v", err)
			}
			IncrementDrainMetric(node.Name)
		} else if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue && node.Spec.Unschedulable {
			// If the node is ready and cordened and do not have special label(it can be on maintenance mode): uncordon it
			err := nodeStruct.UncordonNode()
			if err != nil {
				log.Errorf("Error uncordoning the node: %v", err)
			}
			log.Infof("Uncordoning the node %s.", node.Name)
		}
	}

	return nil
}

func IncrementDrainMetric(nodeName string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8989/drain", nil)
	if err != nil {
		log.Errorf("Error creating request for incrementing drain metric: %v", err)
	}
	req.Header.Set("node_name", nodeName)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error incrementing drain metric: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Errorf("Error incrementing drain metric: %v", resp.Status)
	}
}
