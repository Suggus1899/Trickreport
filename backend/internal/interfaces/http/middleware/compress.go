package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

const (
	// minCompressSize is the minimum response size (in bytes) worth
	// compressing. Smaller responses add overhead without meaningful savings.
	minCompressSize = 1024
)

// gzipResponseWriter wraps an http.ResponseWriter to gzip-encode the response
// body when appropriate.
type gzipResponseWriter struct {
	http.ResponseWriter
	gz          *gzip.Writer
	statusCode  int
	wroteHeader bool
	skipGzip    bool
	buf         []byte
}

func newGzipResponseWriter(w http.ResponseWriter) *gzipResponseWriter {
	return &gzipResponseWriter{
		ResponseWriter: w,
		gz:             gzip.NewWriter(io.Discard), // placeholder; replaced on first write
		statusCode:     http.StatusOK,
	}
}

func (g *gzipResponseWriter) WriteHeader(code int) {
	g.statusCode = code
	g.wroteHeader = true
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	if g.skipGzip {
		if !g.wroteHeader {
			g.ResponseWriter.WriteHeader(g.statusCode)
			g.wroteHeader = true
		}
		return g.ResponseWriter.Write(data)
	}

	// Buffer the first chunk so we can decide whether to compress based on
	// total size.
	g.buf = append(g.buf, data...)

	if len(g.buf) < minCompressSize {
		return len(data), nil
	}

	// Threshold reached — flush as gzip.
	if !g.wroteHeader {
		g.ResponseWriter.Header().Set("Content-Encoding", "gzip")
		g.ResponseWriter.Header().Set("Vary", "Accept-Encoding")
		g.ResponseWriter.WriteHeader(g.statusCode)
		g.wroteHeader = true
	}
	g.gz.Reset(g.ResponseWriter)
	if _, err := g.gz.Write(g.buf); err != nil {
		return len(data), err
	}
	g.buf = nil
	return len(data), nil
}

// Flush flushes any buffered content and closes the gzip writer.
func (g *gzipResponseWriter) Flush() {
	if g.skipGzip || g.buf == nil {
		if !g.wroteHeader {
			g.ResponseWriter.WriteHeader(g.statusCode)
			g.wroteHeader = true
		}
		if g.buf != nil {
			g.ResponseWriter.Write(g.buf)
			g.buf = nil
		}
		return
	}

	// Buffered content is below the threshold — write uncompressed.
	if !g.wroteHeader {
		g.ResponseWriter.WriteHeader(g.statusCode)
		g.wroteHeader = true
	}
	g.ResponseWriter.Write(g.buf)
	g.buf = nil
}

// Compress returns middleware that gzip-compresses responses for clients that
// accept gzip encoding. It skips compression for SSE (text/event-stream),
// WebSocket upgrades, and responses smaller than 1KB.
//
// TODO(wire): wire Compress into server.go's middleware chain.
func Compress() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !acceptsGzip(r) || isSSEOrWebSocket(r) {
				next.ServeHTTP(w, r)
				return
			}

			gw := newGzipResponseWriter(w)
			defer gw.Flush()

			next.ServeHTTP(gw, r)
		})
	}
}

func acceptsGzip(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
}

func isSSEOrWebSocket(r *http.Request) bool {
	accept := strings.ToLower(r.Header.Get("Accept"))
	if strings.Contains(accept, "text/event-stream") {
		return true
	}
	if strings.EqualFold(r.Header.Get("Upgrade"), "websocket") {
		return true
	}
	return false
}
