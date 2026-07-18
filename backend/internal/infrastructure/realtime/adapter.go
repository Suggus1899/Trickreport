package realtime

import (
	ticketapp "github.com/trickreport/backend/internal/application/ticket"
	rt "github.com/trickreport/backend/internal/realtime"
)

// HubAdapter wraps an existing realtime.Hub to implement the
// ticket.EventBroadcaster interface.
type HubAdapter struct {
	hub *rt.Hub
}

// NewHubAdapter creates a new HubAdapter wrapping the given Hub.
func NewHubAdapter(hub *rt.Hub) *HubAdapter {
	return &HubAdapter{hub: hub}
}

// BroadcastEvent delegates to the underlying Hub's BroadcastEvent method,
// casting the string event type to realtime.EventType.
func (a *HubAdapter) BroadcastEvent(tenantID string, eventType string, data any) {
	a.hub.BroadcastEvent(tenantID, rt.EventType(eventType), data)
}

// Compile-time interface assertion.
var _ ticketapp.EventBroadcaster = (*HubAdapter)(nil)
