package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	Requests              prometheus.Counter
	ResponseTimeHistogram *prometheus.HistogramVec
}

func NewMetrics(namespace string) (*Metrics, http.Handler) {
	reg := prometheus.NewRegistry()
	buckets := []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10}

	m := &Metrics{
		Requests: prometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "http_requests_total",
			Help:      "Number of HTTP requests.",
		}),
		ResponseTimeHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "http_server_request_duration_seconds",
			Help:      "Histogram of response time for handler in seconds",
			Buckets:   buckets,
		}, []string{"route", "method", "status_code"}),
	}

	reg.MustRegister(m.Requests, m.ResponseTimeHistogram)

	return m, promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
}
