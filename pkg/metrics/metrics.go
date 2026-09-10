package metrics

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/maadiii/taskmanager/internal/app/port"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	RequestsTotal  *prometheus.CounterVec
	RequestLatency *prometheus.HistogramVec
	TasksCount     prometheus.Gauge
	registry       *prometheus.Registry
}

func New(execer port.Execer) (*Metrics, error) {
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "route", "status"},
	)
	requestLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "request_latency_histogram",
			Help: "HTTP request latency in seconds.",
		},
		[]string{"method", "route"},
	)
	tasksCount := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "tasks_count",
		Help: "Current number of tasks.",
	})
	var count int64
	if err := execer.QueryRow(context.Background(), "SELECT COUNT(*) FROM tasks").Scan(&count); err != nil {
		return nil, err
	}
	tasksCount.Set(float64(count))
	registry := prometheus.NewRegistry()
	if err := registry.Register(requestsTotal); err != nil {
		return nil, err
	}
	if err := registry.Register(requestLatency); err != nil {
		return nil, err
	}
	if err := registry.Register(tasksCount); err != nil {
		return nil, err
	}

	return &Metrics{
		RequestsTotal:  requestsTotal,
		RequestLatency: requestLatency,
		TasksCount:     tasksCount,
		registry:       registry,
	}, nil
}

func (m *Metrics) IncTasks() {
	m.TasksCount.Inc()
}

func (m *Metrics) DecTasks() {
	m.TasksCount.Dec()
}

func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		route := c.FullPath()
		if route == "" {
			route = "unknown"
		}
		method := c.Request.Method
		m.RequestsTotal.WithLabelValues(method, route, strconv.Itoa(c.Writer.Status())).Inc()
		m.RequestLatency.WithLabelValues(method, route).Observe(time.Since(start).Seconds())
	}
}

func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
}
