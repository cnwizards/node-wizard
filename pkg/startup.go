package pkg

import (
	"context"
	"os"
	"time"

	"github.com/pehlicd/node-wizard/pkg/controller"
	"github.com/pehlicd/node-wizard/pkg/logger"
	"github.com/pehlicd/node-wizard/pkg/utils"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

const (
	LeaderElectionLockName  = "node-wizard-lock"
	LeaderElectionNamespace = "node-wizard"
)

func Run() {
	logger.SetupLogger()
	client, err := utils.GetKubernetesClient()
	if err != nil {
		log.Fatalf("Error getting kubernetes client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctr := controller.Setup()

	id := os.Getenv("POD_NAME")

	lock := NewLock(LeaderElectionLockName, id, LeaderElectionNamespace, client)

	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
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
