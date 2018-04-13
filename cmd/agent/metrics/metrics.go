package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// PostsIngested represents the number of reddit posts handled by the agent
	PostsIngested = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "posts_ingested",
		Namespace: "agent",
		Help:      "number of reddit posts handled by agent",
	}, []string{"parsing"})

	// PostsDeleted represents the number of posts that have been deleted after they were parsed
	PostsDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "posts_deleted",
		Namespace: "agent",
		Help:      "number of posts that have been deleted after they were parsed",
	})

	// PostsPolling represents the amount of posts being polled by agent at the moment
	PostsPolling = prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "posts_polling",
		Namespace: "agent",
		Help:      "amount of posts being polled by agent at the moment",
	})

	// CommentsIngested represents the number of reddit comments handled by the agent
	CommentsIngested = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "comments_ingested",
		Namespace: "agent",
		Help:      "number of reddit comments handled by the agent",
	}, []string{"parsing"})

	// CommentsChanged represents the number of reddit comments that were changed after they were parsed.
	// This might contain the same comment more than once as they can be edited any number of times
	CommentsChanged = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "comments_changed",
		Namespace: "agent",
		Help:      "number of reddit comments that were changed after they were parsed. might contain the same comment more than once as they can be edited any number of times",
	})

	// CommentsDeleted represents the number of reddit comments that were deleted after they were parsed
	CommentsDeleted = prometheus.NewCounter(prometheus.CounterOpts{
		Name:      "comments_deleted",
		Namespace: "agent",
		Help:      "number of reddit comments that were deleted after they were parsed",
	})

	// GrawEvents represents any kind of graw interaction being handled
	GrawEvents = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "graw_events",
		Namespace: "agent",
		Help:      "number of graw events",
	}, []string{"type"})
)

// Register registers all metrics
func Register() {
	prometheus.MustRegister(PostsIngested)
	prometheus.MustRegister(PostsDeleted)
	prometheus.MustRegister(PostsPolling)
	prometheus.MustRegister(CommentsIngested)
	prometheus.MustRegister(CommentsChanged)
	prometheus.MustRegister(CommentsDeleted)
	prometheus.MustRegister(GrawEvents)
}

// Serve serves the metrics endpoint
func Serve() error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":9000", nil)
}
