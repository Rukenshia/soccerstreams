package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// MatchthreadsDeleted represents the number of matchthreads deleted from the database
	MatchthreadsDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "matchthreads_deleted",
		Namespace: "sweeper",
		Help:      "number of matchthreads deleted",
	})
)

// Register registers all metrics
func Register() {
	prometheus.MustRegister(MatchthreadsDeleted)
}

// Serve serves the metrics endpoint
func Serve() error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":9000", nil)
}
