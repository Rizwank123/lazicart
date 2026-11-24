package middleware

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// HTTPRequestsTotal counts total HTTP requests
	HTTPRequestsTotal *prometheus.CounterVec

	// HTTPRequestDuration measures request duration
	HTTPRequestDuration *prometheus.HistogramVec

	// HTTPRequestsActive tracks currently active requests
	HTTPRequestsActive *prometheus.GaugeVec
)

// InitHTTPMetrics initializes Prometheus metrics for HTTP middleware
// Call this once during application startup, before registering with Prometheus
func InitHTTPMetrics() {
	HTTPRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"service", "method", "path", "status"},
	)

	HTTPRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service", "method", "path", "status"},
	)

	HTTPRequestsActive = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "http_requests_active",
			Help: "Number of HTTP requests currently being processed",
		},
		[]string{"service"},
	)
}

// RegisterHTTPMetrics registers HTTP metrics with Prometheus
// Call this after InitHTTPMetrics and before starting the HTTP server
func RegisterHTTPMetrics() {
	if HTTPRequestsTotal != nil {
		prometheus.MustRegister(HTTPRequestsTotal)
	}
	if HTTPRequestDuration != nil {
		prometheus.MustRegister(HTTPRequestDuration)
	}
	if HTTPRequestsActive != nil {
		prometheus.MustRegister(HTTPRequestsActive)
	}
}

// MetricsMiddleware collects HTTP metrics for Prometheus
type MetricsMiddleware struct {
	serviceName string
}

// NewMetricsMiddleware creates a new metrics middleware
// serviceName identifies the service in metrics (e.g., "auth-service")
func NewMetricsMiddleware(serviceName string) *MetricsMiddleware {
	// Initialize metrics if not already done
	if HTTPRequestsTotal == nil {
		InitHTTPMetrics()
		RegisterHTTPMetrics()
	}

	return &MetricsMiddleware{
		serviceName: serviceName,
	}
}

// Collect is the middleware handler that collects HTTP metrics
func (m *MetricsMiddleware) Collect(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		start := time.Now()

		// Increment active requests
		if HTTPRequestsActive != nil {
			HTTPRequestsActive.WithLabelValues(m.serviceName).Inc()
			defer HTTPRequestsActive.WithLabelValues(m.serviceName).Dec()
		}

		// Execute the request
		err := next(c)

		// Calculate duration
		duration := time.Since(start).Seconds()

		// Get status code
		status := strconv.Itoa(c.Response().Status)
		method := c.Request().Method
		path := c.Path() // Use route path, not URL path

		// Record metrics
		if HTTPRequestsTotal != nil {
			HTTPRequestsTotal.WithLabelValues(m.serviceName, method, path, status).Inc()
		}

		if HTTPRequestDuration != nil {
			HTTPRequestDuration.WithLabelValues(m.serviceName, method, path, status).Observe(duration)
		}

		return err
	}
}
