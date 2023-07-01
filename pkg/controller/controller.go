package controller

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
)

const (
	informerResyncPeriod = time.Second * 10
)

type Resource struct {
	Group    string
	Version  string
	Resource string
}

type Controller struct {
	Client            DynamicClient
	StopCtx           context.Context
	Indexers          cache.Indexers
	Informer          cache.SharedIndexInformer
	EventHandler      cache.ResourceEventHandlerFuncs
	StopControllerCh  chan struct{} // Channel to receive stop signal.
	ControllerStateCh chan struct{} // Channel to signal the controller has stopped.
}

type DynamicClient interface {
	dynamic.Interface
}

func NewController(client DynamicClient, stopCtx context.Context, indexers cache.Indexers,
	resource Resource, eventHandler cache.ResourceEventHandlerFuncs) (*Controller, error) {

	controlFactory := dynamicinformer.NewDynamicSharedInformerFactory(
		client, informerResyncPeriod,
	)
	informer := controlFactory.ForResource(schema.GroupVersionResource{
		Group:    resource.Group,
		Version:  resource.Version,
		Resource: resource.Resource,
	}).Informer()
	_, err := informer.AddEventHandler(eventHandler)
	if err != nil {
		log.Errorf("Error adding event handler: %s", err)
		return nil, err
	}
	log.Debugf("Added event handler for resource successfully: %s", resource.Resource)

	return &Controller{
		Client:            client,
		StopCtx:           stopCtx,
		Indexers:          indexers,
		Informer:          informer,
		EventHandler:      eventHandler,
		StopControllerCh:  make(chan struct{}),
		ControllerStateCh: make(chan struct{}),
	}, nil
}

func (c *Controller) Run() {
	defer func() {
		close(c.ControllerStateCh)
		log.Info("Controller stopped")
	}()

	// Start the informer.
	go c.Informer.Run(c.StopCtx.Done())
	if !cache.WaitForCacheSync(c.StopCtx.Done(), c.Informer.HasSynced) {
		log.Panic("Timed out waiting for cache to sync")
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
