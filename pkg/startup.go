package pkg

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/cnwizards/node-wizard/pkg/controller"
	"github.com/cnwizards/node-wizard/pkg/logger"
	"github.com/cnwizards/node-wizard/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

const (
	LeaderElectionLockName  string = "node-wizard-lock"
	LeaderElectionNamespace string = "node-wizard"
	MetricsLabel            string = "node_name"
)

var drainCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "node_wizard_drain_count",
		Help: "Number of times the drain function has been called",
	},
	[]string{MetricsLabel},
)

func DrainCounter(w http.ResponseWriter, _ *http.Request) {
	if w.Header().Get(MetricsLabel) == "" {
		log.Errorf("DrainCounter called without node name in header.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Debugf("DrainCounter called for node: %s", w.Header().Get(MetricsLabel))
	drainCounter.WithLabelValues(w.Header().Get(MetricsLabel)).Inc()
}

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
				go func() {
					prometheus.MustRegister(drainCounter)
					http.HandleFunc("/drain", DrainCounter)
					http.Handle("/metrics", promhttp.Handler())
					http.ListenAndServe(":8989", nil)
				}()
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
