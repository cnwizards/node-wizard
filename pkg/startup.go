package pkg

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"time"

	"github.com/cnwizards/node-wizard/pkg/controller"
	"github.com/cnwizards/node-wizard/pkg/logger"
	"github.com/cnwizards/node-wizard/pkg/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

const (
	LeaderElectionLockName string = "node-wizard-lock"
)

func Run() {
	logger.SetupLogger()
	client, err := utils.GetKubernetesClient()
	if err != nil {
		log.Fatalf("Error getting kubernetes client: %v", err)
	}
	// Get the namespace of the pod
	LeaderElectionNamespace := os.Getenv("POD_NAMESPACE")
	if LeaderElectionNamespace == "" {
		log.Fatalf("Error getting namespace of the pod")
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		metricPort := os.Getenv("METRIC_PORT")
		if metricPort == "" {
			metricPort = "8989"
		}
		http.ListenAndServe(":"+metricPort, nil)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctr := controller.Setup()

	id := os.Getenv("POD_NAME")
	shaPodName := sha256.Sum256([]byte(id))
	id = id + "-" + hex.EncodeToString(shaPodName[:])[:8]

	lock := NewLock(LeaderElectionLockName, id, LeaderElectionNamespace, client)

	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second,
		RenewDeadline:   15 * time.Second,
		RetryPeriod:     5 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(_ context.Context) {
				log.Infof("Started leading.")
				//Run the controller
				go ctr.Run()
			},
			OnStoppedLeading: func() {
				log.Infof("No longer leading, transitioning to follower state.")
				//Stop the controller
				ctr.Stop()
			},
			OnNewLeader: func(currentId string) {
				if currentId == id {
					log.Debugf("Same leader as before.")
					return
				}
				log.Infof("New leader is %s", currentId)
			},
		},
	})
}

func NewLock(lockname, podname, namespace string, client *clientset.Clientset) *resourcelock.LeaseLock {
	return &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      lockname,
			Namespace: namespace,
		},
		Client: client.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: podname,
		},
	}
}
