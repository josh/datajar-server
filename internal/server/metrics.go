package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var MetricsRegistry = prometheus.NewRegistry()

var RequestsTotal = promauto.With(MetricsRegistry).NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Tracks the number of HTTP requests.",
	}, []string{"method", "code"},
)

var RequestDuration = promauto.With(MetricsRegistry).NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Tracks the latencies for HTTP requests.",
		Buckets: prometheus.ExponentialBuckets(0.1, 1.5, 5),
	},
	[]string{"method", "code"},
)

var MetricsHandler = promhttp.HandlerFor(
	MetricsRegistry,
	promhttp.HandlerOpts{},
)
