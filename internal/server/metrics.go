package server

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var MetricsRegistry = prometheus.NewRegistry()

var ReadsTotal = promauto.With(MetricsRegistry).NewCounterVec(
	prometheus.CounterOpts{
		Name: "datajar_reads_total",
		Help: "Tracks the number of Data Jar reads.",
	}, []string{"hostname", "ip", "path"},
)

var WritesTotal = promauto.With(MetricsRegistry).NewCounterVec(
	prometheus.CounterOpts{
		Name: "datajar_writes_total",
		Help: "Tracks the number of Data Jar writes.",
	}, []string{"hostname", "ip", "path"},
)

var UnauthorizedTotal = promauto.With(MetricsRegistry).NewCounterVec(
	prometheus.CounterOpts{
		Name: "datajar_unauthorized_total",
		Help: "Tracks the number of unauthorized Data Jar requests.",
	}, []string{"hostname", "ip", "path"},
)

var MetricsHandler = promhttp.HandlerFor(
	MetricsRegistry,
	promhttp.HandlerOpts{},
)
