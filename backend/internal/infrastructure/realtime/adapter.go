package realtime

import (
	"context"
	"github.com/google/uuid"
	ticketapp "github.com/trickreport/backend/internal/application/ticket"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
	rt "github.com/trickreport/backend/internal/realtime"
	"time"
)

// HubAdapter wraps an existing realtime.Hub to implement the
// ticket.EventBroadcaster interface.
type HubAdapter struct {
	hub        *rt.Hub
	recipients RecipientResolver
}

type RecipientResolver interface {
	EventRecipients(context.Context, uuid.UUID, uuid.UUID, bool) ([]uuid.UUID, error)
}

// NewHubAdapter creates a new HubAdapter wrapping the given Hub.
func NewHubAdapter(hub *rt.Hub, resolvers ...RecipientResolver) *HubAdapter {
	a := &HubAdapter{hub: hub}
	if len(resolvers) > 0 {
		a.recipients = resolvers[0]
	}
	return a
}

// BroadcastEvent delegates to the underlying Hub's BroadcastEvent method,
// casting the string event type to realtime.EventType.
func (a *HubAdapter) BroadcastEvent(tenantID uuid.UUID, eventType string, data any) {
	var ticketID uuid.UUID
	includeCreator := true
	var dto any
	switch v := data.(type) {
	case *domainTicket.Notification:
		if v.TenantID != tenantID {
			return
		}
		a.hub.SendToUsers(tenantID, []uuid.UUID{v.UserID}, rt.EventType(eventType), v)
		return
	case *domainTicket.Ticket:
		ticketID = v.ID
		dto = map[string]any{"id": v.ID, "title": v.Title, "status": v.Status, "priority": v.Priority, "created_by": v.CreatedBy, "assigned_to": v.AssignedTo}
	case *domainTicket.Comment:
		ticketID = v.TicketID
		includeCreator = !v.IsInternal
		dto = map[string]any{"id": v.ID, "ticket_id": v.TicketID, "user_id": v.UserID, "content": v.Content, "is_internal": v.IsInternal, "created_at": v.CreatedAt, "user_name": v.UserName}
	case map[string]any:
		if id, ok := v["id"].(uuid.UUID); ok {
			ticketID = id
		} else if id, ok := v["id"].(string); ok {
			ticketID, _ = uuid.Parse(id)
		}
		dto = v
	default:
		return
	}
	if a.recipients == nil || ticketID == uuid.Nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	users, err := a.recipients.EventRecipients(ctx, tenantID, ticketID, includeCreator)
	if err != nil {
		return
	}
	a.hub.SendToUsers(tenantID, users, rt.EventType(eventType), dto)
}

// Compile-time interface assertion.
var _ ticketapp.EventBroadcaster = (*HubAdapter)(nil)
