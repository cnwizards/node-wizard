package controller

import (
	"context"

	"github.com/pehlicd/node-wizard/pkg/utils"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/cache"
)

func Setup() *Controller {
	// Get dynamic kubernetes client.
	var config, err = utils.GetClusterConfig()
	if err != nil {
		log.Panicf("Error getting cluster config: %v", err)
	}
	dClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Panicf("Error getting dynamic client: %v", err)
	}

	// Define resources to watch.
	var resource = Resource{
		Group:    "",
		Version:  "v1",
		Resource: "nodes",
	}

	var eventHandler = cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			err := OnAdd(obj)
			if err != nil {
				log.Errorf("Error adding object: %s", err)
				runtime.HandleError(err)
			}
		},
		UpdateFunc: func(oldObj, obj interface{}) {
			err := OnUpdate(oldObj, obj)
			if err != nil {
				log.Errorf("Error updating object: %s", err)
				runtime.HandleError(err)
			}
		},
		DeleteFunc: func(obj interface{}) {
			err := OnDelete(obj)
			if err != nil {
				log.Errorf("Error deleting object: %s", err)
				runtime.HandleError(err)
			}
		},
	}

	stopCtx := context.Background()

	// Create a new controller.
	controller, err := NewController(dClient, stopCtx, cache.Indexers{}, resource, eventHandler)
	if err != nil {
		log.Panicf("Error creating controller: %s", err)
	}

	return controller
}
