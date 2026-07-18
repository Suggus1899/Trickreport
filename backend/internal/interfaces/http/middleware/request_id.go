package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

const requestIDKey contextKey = "request_id"

// RequestID ensures a request ID exists in the request context. If chi's
// RequestID middleware has already set one (via the X-Request-ID header or a
// generated value), that value is reused. Otherwise a new UUID is generated
// and stored in the context, and the X-Request-ID response header is set so
// clients can correlate logs with responses.
//
// Use RequestIDFromContext to retrieve the request ID in handlers and log
// calls so that every log line carries the request identifier.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := chiMiddleware.GetReqID(r.Context())
		if reqID == "" {
			reqID = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), requestIDKey, reqID)
		w.Header().Set("X-Request-ID", reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext returns the request ID stored in the context, or an
// empty string when no request ID is present (e.g. outside an HTTP request).
// It first checks the requestIDKey set by the RequestID middleware, then
// falls back to chi's request ID (set by chiMiddleware.RequestID) so it works
// whether or not the dedicated RequestID middleware is wired.
func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok && v != "" {
		return v
	}
	if reqID := chiMiddleware.GetReqID(ctx); reqID != "" {
		return reqID
	}
	return ""
}
