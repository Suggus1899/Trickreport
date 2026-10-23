package middleware

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestLimitByIP_ForwardedHeadersCannotRotateBucket(t *testing.T) {
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP", "True-Client-IP"} {
		t.Run(header, func(t *testing.T) {
			limiter := NewRateLimiter(0, 1, time.Minute)
			handler := limiter.LimitByIP(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNoContent)
			}))
			for i, want := range []int{http.StatusNoContent, http.StatusTooManyRequests} {
				req := httptest.NewRequest(http.MethodGet, "/", nil)
				req.RemoteAddr = "198.51.100.7:4000"
				req.Header.Set(header, fmt.Sprintf("203.0.113.%d", i+1))
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, req)
				if rec.Code != want {
					t.Fatalf("request %d status=%d, want %d", i+1, rec.Code, want)
				}
			}
		})
	}
}
