package http

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/trickreport/backend/internal/config"
)

func TestServer_ForwardedHeadersDoNotRewritePeer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv := New(ctx, &config.Config{Env: "development", JWTSecret: "test-secret", JWTExpHours: 8}, nil)
	srv.router.Get("/test-peer", func(w http.ResponseWriter, r *http.Request) {
		//nolint:gosec // G705: test fixture echoes RemoteAddr; no HTML context
		_, _ = w.Write([]byte(r.RemoteAddr))
	})
	for _, header := range []string{"X-Forwarded-For", "X-Real-IP", "True-Client-IP"} {
		t.Run(header, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test-peer", nil)
			req.RemoteAddr = "198.51.100.7:4000"
			req.Header.Set(header, "203.0.113.9")
			rec := httptest.NewRecorder()
			srv.router.ServeHTTP(rec, req)
			if rec.Code != http.StatusOK || rec.Body.String() != "198.51.100.7:4000" {
				t.Fatalf("status=%d peer=%q, want connected peer", rec.Code, rec.Body.String())
			}
		})
	}
}

func TestServer_PublicQuotaDoesNotBlockProtectedAuthentication(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv := New(ctx, &config.Config{Env: "development", JWTSecret: "test-secret", JWTExpHours: 8}, nil)

	limited := false
	for range 1000 {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/csrf", nil)
		rec := httptest.NewRecorder()
		srv.router.ServeHTTP(rec, req)
		if rec.Code == http.StatusTooManyRequests {
			limited = true
			break
		}
		if rec.Code != http.StatusNoContent {
			t.Fatalf("public request: status=%d", rec.Code)
		}
	}
	if !limited {
		t.Fatal("public auth routes must retain a bounded request quota")
	}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/ws", nil)
	rec := httptest.NewRecorder()
	srv.router.ServeHTTP(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("protected authentication: status=%d, want 401 independent of the public quota", rec.Code)
	}
}
