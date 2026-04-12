package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// TaskDuration tracks how long each analyzer task takes.
	TaskDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "safesurf_task_duration_seconds",
			Help:    "Duration of each analyzer task in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"task"},
	)

	// TaskErrors counts errors per task.
	TaskErrors = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "safesurf_task_errors_total",
			Help: "Total number of errors returned by each analyzer task.",
		},
		[]string{"task"},
	)

	// CacheHits counts cache hits keyed by cache-key prefix (task name).
	CacheHits = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "safesurf_cache_hits_total",
			Help: "Total number of cache hits per task.",
		},
		[]string{"task"},
	)

	// CacheMisses counts cache misses keyed by cache-key prefix (task name).
	CacheMisses = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "safesurf_cache_misses_total",
			Help: "Total number of cache misses per task.",
		},
		[]string{"task"},
	)

	// HTTPRequests counts inbound HTTP requests by method, path, and status code.
	HTTPRequests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "safesurf_http_requests_total",
			Help: "Total number of HTTP requests received.",
		},
		[]string{"method", "path", "status"},
	)

	// HTTPDuration tracks inbound request latency by method and path.
	HTTPDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "safesurf_http_request_duration_seconds",
			Help:    "HTTP request latency in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// RiskScore tracks the distribution of computed risk scores (0–100).
	RiskScore = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "safesurf_risk_score",
			Help:    "Distribution of URL risk scores (0–100).",
			Buckets: prometheus.LinearBuckets(0, 10, 11),
		},
	)

	// TrustScore tracks the distribution of computed trust scores (0–100).
	TrustScore = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "safesurf_trust_score",
			Help:    "Distribution of URL trust scores (0–100).",
			Buckets: prometheus.LinearBuckets(0, 10, 11),
		},
	)
)
