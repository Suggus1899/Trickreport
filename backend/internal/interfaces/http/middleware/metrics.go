package middleware

import (
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics collected for every HTTP request.
var (
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests handled.",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpResponseStatus = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_response_status_total",
			Help: "HTTP response status code distribution.",
		},
		[]string{"status"},
	)
)

// Metrics is Prometheus instrumentation middleware. It records the request
// count, duration histogram, and status-code distribution for every request.
func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := chiMiddleware.NewWrapResponseWriter(w, r.ProtoMajor)
		next.ServeHTTP(rec, r)

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(responseStatus(rec.Status()))
		path := "unmatched"
		if route := chi.RouteContext(r.Context()); route != nil && route.RoutePattern() != "" {
			path = route.RoutePattern()
		}

		httpRequestsTotal.WithLabelValues(r.Method, path, status).Inc()
		httpRequestDuration.WithLabelValues(r.Method, path).Observe(duration)
		httpResponseStatus.WithLabelValues(status).Inc()
	})
}

func responseStatus(status int) int {
	if status == 0 {
		return http.StatusOK
	}
	return status
}
