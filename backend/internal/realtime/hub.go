package realtime

import (
	"encoding/json"
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

	mu sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]map[*Client]bool),
		Broadcast:  make(chan BroadcastPayload),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	log.Info().Msg("WebSocket Hub is running")
	for {
		select {
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

	h.Broadcast <- BroadcastPayload{
		TenantID: tenantID,
		Message:  b,
	}
}
