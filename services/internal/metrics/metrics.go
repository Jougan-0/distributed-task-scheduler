package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	TotalTasksProcessed = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "tasks_processed_total",
			Help: "Total number of tasks processed by workers",
		},
		[]string{"status"},
	)

	TaskProcessingTime = prometheus.NewHistogram(
		prometheus.HistogramOpts{
			Name:    "task_processing_time_seconds",
			Help:    "Histogram of task processing time",
			Buckets: prometheus.ExponentialBuckets(0.1, 2, 10),
		},
	)

	TaskRetries = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "task_retries_total",
			Help: "Total number of task retries by type",
		},
		[]string{"task_type"},
	)

	PendingTasksGauge = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "pending_tasks_total",
			Help: "Current number of pending tasks",
		},
	)
)

func RegisterMetrics() {
	prometheus.MustRegister(TotalTasksProcessed)
	prometheus.MustRegister(TaskProcessingTime)
	prometheus.MustRegister(TaskRetries)
	prometheus.MustRegister(PendingTasksGauge)
}

func StartMetricsServer() {
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		http.ListenAndServe(":8080", nil)
	}()
}
