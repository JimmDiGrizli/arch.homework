package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

var (
	metricsRequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "app_request_count",
			Help: "Application Request Count.",
		},
		[]string{"method", "endpoint"},
	)
)

var (
	metricsRequestLatency = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "app_request_latency_seconds",
			Help: "Application Request Latency.",
			// Buckets: prometheus.LinearBuckets(*normMean-5**normDomain, .5**normDomain, 20),
		},
		[]string{"method", "endpoint", "http_status"},
	)
)

func init() {

	prometheus.MustRegister(metricsRequestCount)
	prometheus.MustRegister(metricsRequestLatency)
	prometheus.MustRegister(collectors.NewBuildInfoCollector())
}
