package middleware

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestResponseMiddlewarePreservesWebSocketUpgrade(t *testing.T) {
	for _, tc := range []struct {
		name       string
		middleware func(http.Handler) http.Handler
	}{{"logger", RequestLogger}, {"metrics", Metrics}} {
		t.Run(tc.name, func(t *testing.T) {
			u := websocket.Upgrader{}
			srv := httptest.NewServer(tc.middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				c, err := u.Upgrade(w, r, nil)
				if err == nil {
					_ = c.Close()
				}
			})))
			defer srv.Close()
			c, resp, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(srv.URL, "http"), nil)
			if err != nil {
				status := 0
				if resp != nil {
					status = resp.StatusCode
					_ = resp.Body.Close()
				}
				t.Fatalf("upgrade failed (status %d): %v", status, err)
			}
			_ = c.Close()
		})
	}
}
