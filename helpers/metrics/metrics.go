package metrics

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	// ControllerSummary metrics of controllers call
	ControllerSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "web",
			Subsystem: "controllers",
			Name:      "duration_seconds",
			Help:      "is controller calls metrics",
			MaxAge:    time.Minute,
		},
		[]string{"controller", "status"},
	)
	// ExternalSummary is external calls metrics
	ExternalSummary = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Namespace: "web",
			Subsystem: "external",
			Name:      "duration_seconds",
			Help:      "is external calls metrics",
			MaxAge:    time.Minute,
		},
		[]string{"external", "path", "status"},
	)
)

func init() {
	prometheus.MustRegister(ControllerSummary)
	prometheus.MustRegister(ExternalSummary)
}
