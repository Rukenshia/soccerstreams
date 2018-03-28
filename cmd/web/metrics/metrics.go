package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTPRequests represents the number of HTTP Requests made split by response status
	HTTPRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "http_requests_total",
		Namespace: "web",
		Help:      "number of http requests",
	}, []string{"status"})

	// ImageNotFound represents failed requests to one of our magic image services
	ImageNotFound = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name:      "image_not_found",
		Namespace: "web",
		Help:      "number of image not found errors",
	}, []string{"type"})
)

// Register registers all metrics
func Register() {
	prometheus.MustRegister(HTTPRequests)
	prometheus.MustRegister(ImageNotFound)
}

// Serve serves the metrics endpoint
func Serve() error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":9000", nil)
}
