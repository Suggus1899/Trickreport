package realtime

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func TestConfigureAllowedOrigins_Normalizes(t *testing.T) {
	ConfigureAllowedOrigins([]string{"  HTTPS://Example.COM/  ", "http://localhost:3000/", ""})

	if len(allowedOrigins) != 2 {
		t.Fatalf("expected 2 origins, got %d: %v", len(allowedOrigins), allowedOrigins)
	}
	if allowedOrigins[0] != "https://example.com" {
		t.Errorf("origin[0] = %q", allowedOrigins[0])
	}
	if allowedOrigins[1] != "http://localhost:3000" {
		t.Errorf("origin[1] = %q", allowedOrigins[1])
	}
}

func TestUpgrader_CheckOrigin_EmptyOriginNoWhitelist(t *testing.T) {
	allowedOrigins = []string{}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Del("Origin")
	if !upgrader.CheckOrigin(req) {
		t.Error("expected empty origin to be allowed when no whitelist configured")
	}
}

func TestUpgrader_CheckOrigin_EmptyOriginWithWhitelist(t *testing.T) {
	allowedOrigins = []string{"https://example.com"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Del("Origin")
	if upgrader.CheckOrigin(req) {
		t.Error("expected empty origin to be rejected when whitelist is configured")
	}
}

func TestUpgrader_CheckOrigin_AllowedOrigin(t *testing.T) {
	allowedOrigins = []string{"https://example.com"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://example.com/")
	if !upgrader.CheckOrigin(req) {
		t.Error("expected allowed origin to pass")
	}
}

func TestUpgrader_CheckOrigin_RejectedOrigin(t *testing.T) {
	allowedOrigins = []string{"https://example.com"}
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Origin", "https://evil.com")
	if upgrader.CheckOrigin(req) {
		t.Error("expected non-whitelisted origin to be rejected")
	}
}

func TestServeWs_AndBroadcast(t *testing.T) {
	// Allow all origins for this test
	allowedOrigins = []string{}

	h := NewHub()
	go h.Run()

	tenantID := uuid.New()
	userID := uuid.New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(h, w, r, tenantID, userID)
	}))
	defer server.Close()

	// Convert http:// to ws://
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	// Wait for client registration
	time.Sleep(50 * time.Millisecond)

	// Broadcast a message — the client should receive it
	h.BroadcastEvent(tenantID, EventTicketCreated, map[string]string{"id": "123"})

	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, msg, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	if len(msg) == 0 {
		t.Error("expected non-empty message")
	}
	if !strings.Contains(string(msg), "TICKET_CREATED") {
		t.Errorf("expected message to contain TICKET_CREATED, got %s", string(msg))
	}
}

func TestServeWs_UpgradeError(t *testing.T) {
	allowedOrigins = []string{}

	h := NewHub()

	// Use a response that can't be hijacked (not an httptest.NewRecorder properly)
	// Actually, httptest.NewRecorder doesn't support hijacking, so the upgrade will fail.
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Connection", "upgrade")
	r.Header.Set("Upgrade", "websocket")
	r.Header.Set("Sec-WebSocket-Version", "13")
	r.Header.Set("Sec-WebSocket-Key", "REDACTED")

	// This should fail because httptest.NewRecorder doesn't implement http.Hijacker
	ServeWs(h, w, r, uuid.New(), uuid.New())

	// If we reach here without panicking, the error path was exercised.
}
