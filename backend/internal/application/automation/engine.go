package automation

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	domainautomation "github.com/trickreport/backend/internal/domain/automation"
)

// Event represents a trigger event that the engine can evaluate.
type Event struct {
	Type     string // "ticket_created", "ticket_status_changed", "ticket_priority_changed", "sla_breach"
	TenantID uuid.UUID
	TicketID uuid.UUID
	UserID   uuid.UUID
	OldValue string
	NewValue string
}

// TicketSnapshot is a lightweight view of a ticket used for condition evaluation.
type TicketSnapshot struct {
	ID       uuid.UUID
	Status   string
	Priority string
	Category string
	Title    string
}

// TicketExecutor is the port for executing actions on tickets.
type TicketExecutor interface {
	GetTicket(ctx context.Context, tenantID, ticketID uuid.UUID) (*TicketSnapshot, error)
	SetPriority(ctx context.Context, tenantID, ticketID uuid.UUID, priority string) error
	SetStatus(ctx context.Context, tenantID, ticketID uuid.UUID, status string) error
	AssignTo(ctx context.Context, tenantID, ticketID, userID uuid.UUID) error
	AddComment(ctx context.Context, tenantID, ticketID uuid.UUID, content string, isInternal bool) error
	AddTag(ctx context.Context, tenantID, ticketID uuid.UUID, tag string) error
}

// Engine evaluates and executes automation rules when triggers fire.
type Engine struct {
	repo    Repository
	tickets TicketExecutor
	logger  zerolog.Logger
}

// NewEngine creates a new automation engine.
func NewEngine(repo Repository, tickets TicketExecutor, logger zerolog.Logger) *Engine {
	return &Engine{repo: repo, tickets: tickets, logger: logger}
}

// triggerFromEvent maps an engine event type to the stored rule trigger_type.
// Stored trigger types: "ticket_created", "status_changed", "priority_changed", "sla_breach".
func triggerFromEvent(eventType string) string {
	switch eventType {
	case "ticket_created":
		return "ticket_created"
	case "ticket_status_changed":
		return "status_changed"
	case "ticket_priority_changed":
		return "priority_changed"
	case "sla_breach":
		return "sla_breach"
	default:
		return eventType
	}
}

// Evaluate runs all active rules for a given trigger event.
func (e *Engine) Evaluate(ctx context.Context, tenantID uuid.UUID, event Event) error {
	if e == nil || e.repo == nil || e.tickets == nil {
		return nil
	}

	trigger := triggerFromEvent(event.Type)

	rules, err := e.repo.List(ctx, tenantID)
	if err != nil {
		return fmt.Errorf("automation engine: list rules: %w", err)
	}

	// Fetch the ticket snapshot once for condition evaluation.
	var snapshot *TicketSnapshot
	if event.TicketID != uuid.Nil {
		s, err := e.tickets.GetTicket(ctx, tenantID, event.TicketID)
		if err != nil {
			e.logger.Warn().Err(err).Str("event", event.Type).Msg("automation engine: failed to load ticket snapshot")
		} else {
			snapshot = s
		}
	}

	for _, rule := range rules {
		if !rule.IsActive || rule.TriggerType != trigger {
			continue
		}

		if !e.matchConditions(rule.Conditions, snapshot, event) {
			continue
		}

		if err := e.executeActions(ctx, tenantID, rule, event); err != nil {
			e.logger.Error().Err(err).
				Str("rule_id", rule.ID.String()).
				Str("rule_name", rule.Name).
				Msg("automation engine: failed to execute actions")
		}
	}

	return nil
}

// matchConditions evaluates simple equality conditions against the ticket snapshot.
// Conditions is a map of field name -> expected value (e.g. {"priority": "high"}).
func (e *Engine) matchConditions(conditions map[string]any, snapshot *TicketSnapshot, event Event) bool {
	if len(conditions) == 0 {
		return true
	}

	for field, expected := range conditions {
		expectedStr := fmt.Sprintf("%v", expected)
		var actual string

		switch strings.ToLower(field) {
		case "status":
			if snapshot != nil {
				actual = snapshot.Status
			}
		case "priority":
			if snapshot != nil {
				actual = snapshot.Priority
			}
		case "category":
			if snapshot != nil {
				actual = snapshot.Category
			}
		case "title":
			if snapshot != nil {
				actual = snapshot.Title
			}
		case "old_value", "oldvalue":
			actual = event.OldValue
		case "new_value", "newvalue":
			actual = event.NewValue
		default:
			// Unknown condition field — fail safe (no match).
			return false
		}

		if !strings.EqualFold(actual, expectedStr) {
			return false
		}
	}
	return true
}

// executeActions runs all actions defined in a rule.
// Each action is a map[string]any with an "type" key and action-specific fields.
func (e *Engine) executeActions(ctx context.Context, tenantID uuid.UUID, rule domainautomation.Rule, event Event) error {
	for _, raw := range rule.Actions {
		action, ok := raw.(map[string]any)
		if !ok {
			e.logger.Warn().Str("rule_id", rule.ID.String()).Msg("automation engine: skipping non-map action")
			continue
		}

		actionType, _ := action["type"].(string)
		if actionType == "" {
			e.logger.Warn().Str("rule_id", rule.ID.String()).Msg("automation engine: skipping action without type")
			continue
		}

		if err := e.executeAction(ctx, tenantID, actionType, action, event); err != nil {
			e.logger.Error().Err(err).
				Str("action", actionType).
				Str("rule_id", rule.ID.String()).
				Msg("automation engine: action failed")
			// Continue executing remaining actions.
		}
	}
	return nil
}

func (e *Engine) executeAction(ctx context.Context, tenantID uuid.UUID, actionType string, action map[string]any, event Event) error {
	switch actionType {
	case "set_priority":
		priority, _ := action["priority"].(string)
		if priority == "" {
			priority, _ = action["value"].(string)
		}
		if priority == "" {
			return fmt.Errorf("set_priority: missing 'priority' field")
		}
		return e.tickets.SetPriority(ctx, tenantID, event.TicketID, priority)

	case "set_status":
		status, _ := action["status"].(string)
		if status == "" {
			status, _ = action["value"].(string)
		}
		if status == "" {
			return fmt.Errorf("set_status: missing 'status' field")
		}
		return e.tickets.SetStatus(ctx, tenantID, event.TicketID, status)

	case "assign_to":
		userIDStr, _ := action["user_id"].(string)
		if userIDStr == "" {
			return fmt.Errorf("assign_to: missing 'user_id' field")
		}
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return fmt.Errorf("assign_to: invalid user_id: %w", err)
		}
		return e.tickets.AssignTo(ctx, tenantID, event.TicketID, userID)

	case "add_comment":
		content, _ := action["content"].(string)
		if content == "" {
			content, _ = action["value"].(string)
		}
		if content == "" {
			return fmt.Errorf("add_comment: missing 'content' field")
		}
		isInternal := true
		if v, ok := action["is_internal"]; ok {
			isInternal, _ = v.(bool)
		}
		return e.tickets.AddComment(ctx, tenantID, event.TicketID, content, isInternal)

	case "add_tag":
		tag, _ := action["tag"].(string)
		if tag == "" {
			tag, _ = action["value"].(string)
		}
		if tag == "" {
			return fmt.Errorf("add_tag: missing 'tag' field")
		}
		return e.tickets.AddTag(ctx, tenantID, event.TicketID, tag)

	default:
		e.logger.Warn().Str("action", actionType).Msg("automation engine: unknown action type")
		return nil
	}
}
