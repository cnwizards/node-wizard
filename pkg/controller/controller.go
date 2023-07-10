package controller

import (
	"context"
	"time"

	"github.com/cnwizards/node-wizard/pkg/utils"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/tools/cache"
)

const (
	informerResyncPeriod = time.Second * 10
)

type Controller struct {
	Informer          cache.SharedIndexInformer
	StopCtx           context.Context
	StopControllerCh  chan struct{}
	ControllerStateCh chan struct{}
}

func Setup() *Controller {

	// Initialize the Kubernetes clientset and other necessary variables
	clientset, err := utils.GetKubernetesClient()
	if err != nil {
		panic(err.Error())
	}

	controlFactory := informers.NewSharedInformerFactory(clientset, informerResyncPeriod)
	controller := controlFactory.Core().V1().Nodes().Informer()
	defer runtime.HandleCrash()
	_, err = controller.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			err := OnAdd(obj)
			if err != nil {
				runtime.HandleError(err)
			}
		},
		DeleteFunc: func(obj interface{}) {
			err := OnDelete(obj)
			if err != nil {
				runtime.HandleError(err)
			}
		},
		UpdateFunc: func(_, obj interface{}) {
			node := obj.(*corev1.Node)
			err := OnUpdate(node)
			if err != nil {
				runtime.HandleError(err)
			}
		},
	})
	if err != nil {
		log.Errorf("An error happened while adding the event handlers. See detailed: %v", err)
		runtime.HandleError(err)
	}

	return &Controller{
		Informer:          controller,
		StopCtx:           context.Background(),
		StopControllerCh:  make(chan struct{}),
		ControllerStateCh: make(chan struct{}),
	}
}

func (c *Controller) Run() {
	defer func() {
		close(c.ControllerStateCh)
		log.Infof("Controller stopped unexpectedly.") // Might not be unexpected :D
	}()
	defer runtime.HandleCrash()

	go c.Informer.Run(c.StopCtx.Done())
	if !cache.WaitForCacheSync(c.StopCtx.Done(), c.Informer.HasSynced) {
		log.Errorf("Timed out waiting for caches to sync")
	}

	// Wait for a signal to stop the controller.
	select {
	case <-c.StopCtx.Done():
		return
	case <-c.StopControllerCh:
		return
	}
}

// Stop stops the controller.
func (c *Controller) Stop() {
	select {
	case <-c.ControllerStateCh:
		// Controller has already stopped.
		log.Debugf("Controller has already stopped.")
	case c.StopControllerCh <- struct{}{}:
		// Signal the controller to stop.
		log.Debugf("Signaled controller to stop.")
		<-c.ControllerStateCh
	}
}
