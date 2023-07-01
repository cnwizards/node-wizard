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

func OnUpdate(_, obj interface{}) error {
	//TODO: Add logic to check conditions.
	log.Debugf("Received update event: %v", obj.(metav1.Object).GetName())
	node := obj.(*corev1.Node)
	for _, condition := range node.Status.Conditions {
		log.Debugf("Node condition: %+v", condition)
		if condition.Status != corev1.ConditionTrue {
			log.Debugf("Node condition status %s is not true", condition.Type)
			continue
		} else {
			log.Debugf("Node condition status %s is true", condition.Type)
			if condition.Type == corev1.NodeReady {
				log.Debugf("Node condition type is %s", condition.Type)
				log.Infof("Node %s is not ready for reason %s so draining the node.", node.Name, condition.Reason)
				err := utils.DrainNode()
				if err != nil {
					log.Errorf("Error draining node: %v", err)
				}
			}
		}
	}
	return nil
}
