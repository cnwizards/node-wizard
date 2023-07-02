package controller

import (
	"github.com/pehlicd/node-wizard/pkg/utils"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
		log.Debugf("Node condition: %+v", condition)
		// If the node is ready, and schedulable: do nothing
		if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue && !node.Spec.Unschedulable {
			log.Debugf("Node %s is ready", node.Name)
			// If the node is notReady, and schedulable: drain it and evict sts pods
			// TODO: Add a function to evict sts pods
		} else if condition.Type == corev1.NodeReady && condition.Status != corev1.ConditionTrue && !node.Spec.Unschedulable {
			log.Debugf("Node %s is not ready", node.Name)
			err := nodeStruct.DrainNode()
			if err != nil {
				log.Errorf("Error draining node: %v", err)
			}
			// If the node is ready and cordened and do not have special label(it can be on maintenance mode): uncordon it
		} else if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue && node.Spec.Unschedulable && node.Labels["node-wizard/ignore"] != "true" {
			//only call uncordon function
			log.Debugf("hello %s", node.Name)
		} else if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue && node.Spec.Unschedulable && node.Labels["node-wizard/ignore"] == "true" {
			//do nothing
			log.Debugf("Node %s is ready and ignored", node.Name)
		} else {
			log.Debugf("Node %s is not ready and ignored", node.Name)
		}
	}

	return nil
}
