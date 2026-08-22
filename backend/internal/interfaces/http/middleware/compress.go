package middleware

import (
	"compress/gzip"
	"net/http"
	"strconv"
	"strings"
)

const (
	// minCompressSize is the minimum response size (in bytes) worth
	// compressing. Smaller responses add overhead without meaningful savings.
	minCompressSize = 1024
)

// gzipResponseWriter wraps an http.ResponseWriter to gzip-encode the response
// body when appropriate.
//
// The compression decision depends on the total body size, which isn't known
// until the handler has written it, so the response is buffered until it
// either crosses minCompressSize (switch to gzip) or the handler finishes
// (send it uncompressed). WriteHeader therefore only records the status —
// the real header is not sent until that decision is made, since
// Content-Encoding has to go out with it.
type gzipResponseWriter struct {
	http.ResponseWriter
	gz         *gzip.Writer
	statusCode int
	headerSent bool // the underlying writer's header has been written
	gzipping   bool // the gzip stream has been started
	buf        []byte
}

func newGzipResponseWriter(w http.ResponseWriter) *gzipResponseWriter {
	return &gzipResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (g *gzipResponseWriter) WriteHeader(code int) {
	g.statusCode = code
}

// sendHeader writes the recorded status to the underlying writer exactly once.
func (g *gzipResponseWriter) sendHeader() {
	if g.headerSent {
		return
	}
	g.ResponseWriter.WriteHeader(g.statusCode)
	g.headerSent = true
}

func (g *gzipResponseWriter) Write(data []byte) (int, error) {
	if g.gzipping {
		if _, err := g.gz.Write(data); err != nil {
			return 0, err
		}
		return len(data), nil
	}

	g.buf = append(g.buf, data...)
	if len(g.buf) < minCompressSize {
		return len(data), nil
	}

	// Threshold reached — switch to gzip for this and all later writes.
	h := g.ResponseWriter.Header()
	h.Set("Content-Encoding", "gzip")
	h.Add("Vary", "Accept-Encoding")
	// The handler's Content-Length (if any) describes the uncompressed body.
	h.Del("Content-Length")
	g.sendHeader()

	g.gz = gzip.NewWriter(g.ResponseWriter)
	g.gzipping = true

	buffered := g.buf
	g.buf = nil
	if _, err := g.gz.Write(buffered); err != nil {
		return 0, err
	}
	return len(data), nil
}

// Close finalizes the response: it closes the gzip stream (flushing the
// trailer, without which the body is unreadable) or, when the body stayed
// under the compression threshold, writes it out uncompressed.
func (g *gzipResponseWriter) Close() {
	if g.gzipping {
		_ = g.gz.Close()
		return
	}

	if len(g.buf) > 0 {
		g.ResponseWriter.Header().Set("Content-Length", strconv.Itoa(len(g.buf)))
	}
	g.sendHeader()
	if len(g.buf) > 0 {
		_, _ = g.ResponseWriter.Write(g.buf)
		g.buf = nil
	}
}

// Flush implements http.Flusher. Buffered content is pushed downstream so
// streaming handlers aren't stalled by the size-based buffering above.
func (g *gzipResponseWriter) Flush() {
	if g.gzipping {
		_ = g.gz.Flush()
		if f, ok := g.ResponseWriter.(http.Flusher); ok {
			f.Flush()
		}
		return
	}

	g.sendHeader()
	if len(g.buf) > 0 {
		_, _ = g.ResponseWriter.Write(g.buf)
		g.buf = nil
	}
	if f, ok := g.ResponseWriter.(http.Flusher); ok {
		f.Flush()
	}
}

// Compress returns middleware that gzip-compresses responses for clients that
// accept gzip encoding. It skips compression for SSE (text/event-stream),
// WebSocket upgrades, and responses smaller than 1KB.
func Compress() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !acceptsGzip(r) || isSSEOrWebSocket(r) {
				next.ServeHTTP(w, r)
				return
			}

			gw := newGzipResponseWriter(w)
			defer gw.Close()

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
