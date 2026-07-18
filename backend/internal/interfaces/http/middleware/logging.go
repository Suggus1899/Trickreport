package middleware

import (
	"context"
	"net/http"
	"time"

	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const loggerKey contextKey = "logger"

// LoggerFromContext returns the request-scoped zerolog logger. When no
// logger is present in the context (e.g. outside an HTTP request), it falls
// back to the package default logger.
//
// It returns a pointer so that pointer-receiver methods (e.g. Error, Warn)
// can be called on it — a value returned from a context lookup is not
// addressable and therefore cannot invoke pointer-receiver methods.
func LoggerFromContext(ctx context.Context) *zerolog.Logger {
	if l, ok := ctx.Value(loggerKey).(*zerolog.Logger); ok {
		return l
	}
	return &log.Logger
}

// contextWithLogger stores a zerolog logger in the request context.
func contextWithLogger(ctx context.Context, l zerolog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, &l)
}

// statusRecorder wraps http.ResponseWriter to capture the status code for
// structured request logging.
type statusRecorder struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (r *statusRecorder) WriteHeader(code int) {
	if r.wroteHeader {
		return
	}
	r.status = code
	r.wroteHeader = true
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wroteHeader {
		r.status = http.StatusOK
		r.wroteHeader = true
	}
	return r.ResponseWriter.Write(b)
}

// RequestLogger is structured request-logging middleware. It reads the chi
// request ID from the context, builds a zerolog logger carrying that ID as a
// field, stores it in the request context (so handlers can retrieve it via
// LoggerFromContext), and logs method, path, status, and duration for every
// request.
//
// It replaces chi's default middleware.Logger so that request logs include
// the request_id field and use zerolog's structured output.
func RequestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		reqID := chiMiddleware.GetReqID(r.Context())
		reqLogger := log.With().Str("request_id", reqID).Logger()

		ctx := contextWithLogger(r.Context(), reqLogger)

		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r.WithContext(ctx))

		reqLogger.Info().
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", rec.status).
			Int64("duration_ms", time.Since(start).Milliseconds()).
			Msg("request completed")
	})
}
