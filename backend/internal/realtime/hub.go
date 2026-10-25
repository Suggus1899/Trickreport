package realtime

import (
	"context"
	"encoding/json"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients mapped by tenant_id.
	// We want to broadcast only to clients belonging to the same tenant.
	// Structure: map[tenant_id]map[*Client]bool
	clients map[uuid.UUID]map[*Client]bool

	// Inbound messages from the clients.
	Broadcast chan BroadcastPayload

	// Register requests from the clients.
	Register chan *Client

	// Unregister requests from clients.
	Unregister chan *Client

	mu   sync.RWMutex
	done chan struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]map[*Client]bool),
		Broadcast:  make(chan BroadcastPayload),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		done:       make(chan struct{}),
	}
}

func (h *Hub) Run(contexts ...context.Context) {
	ctx := context.Background()
	if len(contexts) > 0 {
		ctx = contexts[0]
	}
	defer close(h.done)
	log.Info().Msg("WebSocket Hub is running")
	for {
		select {
		case <-ctx.Done():
			h.mu.Lock()
			for _, clients := range h.clients {
				for client := range clients {
					close(client.send)
					if client.conn != nil {
						_ = client.conn.Close()
					}
				}
			}
			h.clients = make(map[uuid.UUID]map[*Client]bool)
			h.mu.Unlock()
			return
		case client := <-h.Register:
			h.mu.Lock()
			if _, ok := h.clients[client.TenantID]; !ok {
				h.clients[client.TenantID] = make(map[*Client]bool)
			}
			h.clients[client.TenantID][client] = true
			h.mu.Unlock()
			log.Debug().Str("tenant_id", client.TenantID.String()).Str("user_id", client.UserID.String()).Msg("Client connected")

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.TenantID]; ok {
				if _, ok := h.clients[client.TenantID][client]; ok {
					delete(h.clients[client.TenantID], client)
					close(client.send)
					if len(h.clients[client.TenantID]) == 0 {
						delete(h.clients, client.TenantID)
					}
					log.Debug().Str("tenant_id", client.TenantID.String()).Str("user_id", client.UserID.String()).Msg("Client disconnected")
				}
			}
			h.mu.Unlock()

		case payload := <-h.Broadcast:
			h.mu.Lock()
			clients := h.clients[payload.TenantID]
			for client := range clients {
				if payload.Scoped && !slices.Contains(payload.UserIDs, client.UserID) {
					continue
				}
				select {
				case client.send <- payload.Message:
				default:
					close(client.send)
					delete(h.clients[payload.TenantID], client)
				}
			}
			h.mu.Unlock()
		}
	}
}

// BroadcastEvent is a helper to encode and send a structured event.
func (h *Hub) BroadcastEvent(tenantID uuid.UUID, eventType EventType, data any) {
	msg := Message{
		Type:     eventType,
		TenantID: tenantID,
		Data:     data,
	}

	b, err := json.Marshal(msg)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal websocket message")
		return
	}

	select {
	case h.Broadcast <- BroadcastPayload{
		TenantID: tenantID,
		Message:  b,
	}:
	case <-h.done:
	}
}

func (h *Hub) SendToUsers(tenantID uuid.UUID, userIDs []uuid.UUID, eventType EventType, data any) {
	if len(userIDs) == 0 {
		return
	}
	b, err := json.Marshal(Message{Type: eventType, TenantID: tenantID, Data: data})
	if err != nil {
		return
	}
	select {
	case h.Broadcast <- BroadcastPayload{TenantID: tenantID, Message: b, UserIDs: userIDs, Scoped: true}:
	case <-h.done:
	}
}
