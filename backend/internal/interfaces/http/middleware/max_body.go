package middleware

import (
	"net/http"

	"github.com/trickreport/backend/internal/interfaces/http/response"
)

// MaxBodySize returns middleware that limits the size of the request body.
// If the declared Content-Length exceeds maxBytes, a 413 Payload Too Large
// response is sent immediately. The body is also wrapped with
// http.MaxBytesReader so streaming bodies that exceed the limit are caught
// during read.
//
// TODO(wire): wire MaxBodySize into server.go's middleware chain.
func MaxBodySize(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxBytes {
				response.Error(w, http.StatusRequestEntityTooLarge, "request body too large")
				return
			}
			if r.Body != nil && maxBytes > 0 {
				r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			}
			next.ServeHTTP(w, r)
		})
	}
}
