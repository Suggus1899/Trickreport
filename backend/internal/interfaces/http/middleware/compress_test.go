package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestCompress_LargeResponse_DecodesCorrectly is a regression test: the
// original implementation never finalized the gzip stream (no Close/Flush),
// so any response over the compression threshold came back as a truncated,
// unreadable gzip blob to every client that sent Accept-Encoding: gzip —
// which is every browser, on every JSON response over ~1KB.
func TestCompress_LargeResponse_DecodesCorrectly(t *testing.T) {
	body := strings.Repeat("a", minCompressSize+500)

	handler := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") != "gzip" {
		t.Fatalf("Content-Encoding = %q, want gzip", rec.Header().Get("Content-Encoding"))
	}

	gz, err := gzip.NewReader(rec.Body)
	if err != nil {
		t.Fatalf("response is not valid gzip: %v", err)
	}
	defer gz.Close()
	decoded, err := io.ReadAll(gz)
	if err != nil {
		t.Fatalf("failed to read gzip stream to completion (truncated stream): %v", err)
	}
	if string(decoded) != body {
		t.Errorf("decoded body does not match original (len %d vs %d)", len(decoded), len(body))
	}
}

// TestCompress_SmallResponse_Uncompressed verifies responses under the
// threshold are sent as-is, without a broken/empty Content-Encoding header.
func TestCompress_SmallResponse_Uncompressed(t *testing.T) {
	body := "short response"

	handler := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(body))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Accept-Encoding", "gzip")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Header().Get("Content-Encoding") == "gzip" {
		t.Error("small response should not be gzip-encoded")
	}
	if rec.Body.String() != body {
		t.Errorf("body = %q, want %q", rec.Body.String(), body)
	}
}

// TestCompress_NoAcceptEncoding_PassesThrough verifies clients that don't
// send Accept-Encoding: gzip get the handler's own status/body untouched.
func TestCompress_NoAcceptEncoding_PassesThrough(t *testing.T) {
	handler := Compress()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte("hello"))
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusCreated)
	}
	if rec.Body.String() != "hello" {
		t.Errorf("body = %q, want %q", rec.Body.String(), "hello")
	}
}
