package realtime

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 4096
)

// allowedOrigins holds the whitelist of origins permitted to open WebSocket
// connections. It is configured at startup via ConfigureAllowedOrigins.
var allowedOrigins = []string{}

// ConfigureAllowedOrigins sets the list of origins allowed by the WebSocket
// upgrader. Origins are compared case-insensitively, without trailing slashes.
// Call this once at server startup before any WebSocket connection is served.
func ConfigureAllowedOrigins(origins []string) {
	normalized := make([]string, 0, len(origins))
	for _, o := range origins {
		o = strings.ToLower(strings.TrimSpace(o))
		o = strings.TrimRight(o, "/")
		if o != "" {
			normalized = append(normalized, o)
		}
	}
	allowedOrigins = normalized
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := strings.ToLower(strings.TrimSpace(r.Header.Get("Origin")))
		if origin == "" {
			// Non-browser clients (curl, etc.) have no Origin header.
			// Allow them only when no whitelist is configured (dev mode).
			return len(allowedOrigins) == 0
		}
		origin = strings.TrimRight(origin, "/")
		for _, allowed := range allowedOrigins {
			if allowed == origin {
				return true
			}
		}
		log.Warn().Str("origin", origin).Msg("websocket origin rejected")
		return false
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub      *Hub
	TenantID uuid.UUID
	UserID   uuid.UUID

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
	done chan struct{}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		if c.done != nil {
			close(c.done)
		}
		select {
		case c.Hub.Unregister <- c:
		case <-c.Hub.done:
		}
		_ = c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}
	c.conn.SetPongHandler(func(string) error { return c.conn.SetReadDeadline(time.Now().Add(pongWait)) })
	for {
		_, _, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error().Err(err).Msg("websocket close error")
			}
			break
		}
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if !ok {
				// The hub closed the channel.
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(message); err != nil {
				_ = w.Close()
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				if _, err := w.Write([]byte{'\n'}); err != nil {
					_ = w.Close()
					return
				}
				if _, err := w.Write(<-c.send); err != nil {
					_ = w.Close()
					return
				}
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ServeWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request, tenantID uuid.UUID, userID uuid.UUID, guards ...func(context.Context) error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Error().Err(err).Msg("failed to upgrade websocket connection")
		return
	}

	client := &Client{
		Hub:      hub,
		TenantID: tenantID,
		UserID:   userID,
		conn:     conn,
		send:     make(chan []byte, 256),
		done:     make(chan struct{}),
	}
	client.Hub.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	if len(guards) > 0 {
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-client.done:
					return
				case <-hub.done:
					return
				case <-ticker.C:
					ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
					err := guards[0](ctx)
					cancel()
					if err != nil {
						_ = client.conn.Close()
						return
					}
				}
			}
		}()
	}
}
