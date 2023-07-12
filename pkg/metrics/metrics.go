package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const (
	MetricsLabel string = "node_name"
)

var drainCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "node_wizard_drain_count",
		Help: "Number of times the drain function has been called",
	},
	[]string{MetricsLabel},
)

var uncordonCounter = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "node_wizard_uncordon_count",
		Help: "Number of times the uncordon function has been called",
	},
	[]string{MetricsLabel},
)

func init() {
	prometheus.MustRegister(drainCounter)
	prometheus.MustRegister(uncordonCounter)
}

func IncrementDrainCounter(nodeName string) {
	log.Debugf("DrainCounter called for node: %s", nodeName)
	drainCounter.WithLabelValues(nodeName).Inc()
}

func IncrementUncordonCounter(nodeName string) {
	log.Debugf("UncordonCounter called for node: %s", nodeName)
	uncordonCounter.WithLabelValues(nodeName).Inc()
}
