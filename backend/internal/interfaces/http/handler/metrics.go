package handler

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// MetricsHandler exposes the Prometheus metrics endpoint at /metrics.
//
// TODO: Wire this handler into the router in server.go (do NOT edit server.go
// from this task — register it with r.Handle("/metrics", handler.MetricsHandler())).
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
