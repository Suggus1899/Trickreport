package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func TestLoggerFromContext_Default(t *testing.T) {
	l := LoggerFromContext(context.Background())
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestLoggerFromContext_WithLogger(t *testing.T) {
	custom := zerolog.Nop()
	ctx := contextWithLogger(context.Background(), custom)
	l := LoggerFromContext(ctx)
	if l == nil {
		t.Fatal("expected non-nil logger")
	}
	// The pointer stored should not be the default package logger.
	// We can't compare directly, but we can verify it's non-nil and
	// that calling a method on it doesn't panic.
	_ = l.Info()
}

func TestStatusRecorder_WriteHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	sr := &statusRecorder{ResponseWriter: rec, status: http.StatusOK}

	sr.WriteHeader(http.StatusCreated)
	if sr.status != http.StatusCreated {
		t.Errorf("status = %d, want %d", sr.status, http.StatusCreated)
	}
	if !sr.wroteHeader {
		t.Error("wroteHeader should be true")
	}
	if rec.Code != http.StatusCreated {
		t.Errorf("underlying code = %d, want %d", rec.Code, http.StatusCreated)
	}

	// Second call should be ignored
	sr.WriteHeader(http.StatusInternalServerError)
	if sr.status != http.StatusCreated {
		t.Errorf("status = %d, want %d (should not change)", sr.status, http.StatusCreated)
	}
}

func TestStatusRecorder_Write_NoHeaderYet(t *testing.T) {
	rec := httptest.NewRecorder()
	sr := &statusRecorder{ResponseWriter: rec, status: http.StatusOK}

	n, err := sr.Write([]byte("hello"))
	if err != nil {
		t.Fatalf("write error: %v", err)
	}
	if n != 5 {
		t.Errorf("wrote %d bytes, want 5", n)
	}
	if sr.status != http.StatusOK {
		t.Errorf("status = %d, want %d", sr.status, http.StatusOK)
	}
	if !sr.wroteHeader {
		t.Error("wroteHeader should be true after Write")
	}
}

func TestRequestLogger_Basic(t *testing.T) {
	r := chi.NewRouter()
	r.Use(RequestLogger)

	called := false
	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		called = true
		l := LoggerFromContext(r.Context())
		if l == nil {
			t.Error("expected logger in context")
		}
		w.WriteHeader(http.StatusAccepted)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if !called {
		t.Error("handler not called")
	}
	if rec.Code != http.StatusAccepted {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusAccepted)
	}
}

func TestRequestLogger_DefaultStatus(t *testing.T) {
	r := chi.NewRouter()
	r.Use(RequestLogger)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		// Don't call WriteHeader; just write body.
		_, _ = w.Write([]byte("ok"))
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}

func TestRequestLogger_WithReqID(t *testing.T) {
	r := chi.NewRouter()
	r.Use(chiMiddleware.RequestID)
	r.Use(RequestLogger)

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		reqID := chiMiddleware.GetReqID(r.Context())
		if reqID == "" {
			t.Error("expected non-empty request ID in context")
		}
		w.WriteHeader(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
}
