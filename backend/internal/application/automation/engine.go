package automation

import (
	"context"
	"errors"
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

// EmailSender is the port for the "send_email" automation action.
type EmailSender interface {
	Send(to, subject, body string) error
}

// Engine evaluates and executes automation rules when triggers fire.
type Engine struct {
	repo        Repository
	tickets     TicketExecutor
	emailSender EmailSender
	logger      zerolog.Logger
}

// NewEngine creates a new automation engine.
func NewEngine(repo Repository, tickets TicketExecutor, logger zerolog.Logger) *Engine {
	return &Engine{repo: repo, tickets: tickets, logger: logger}
}

// SetEmailSender wires the email sender used by the "send_email" action.
// Optional — when unset, matching actions are skipped with a warning log,
// same as an unknown action type.
func (e *Engine) SetEmailSender(sender EmailSender) {
	e.emailSender = sender
}

// triggerFromEvent maps an engine event type to the stored rule trigger_type.
// Stored trigger types: "ticket_created", "status_changed", "priority_changed", "sla_breach", "escalation".
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
	case "escalation":
		return "escalation"
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

	var failures []error
	for _, rule := range rules {
		if !rule.IsActive || rule.TriggerType != trigger {
			continue
		}

		if !e.matchConditions(rule.Conditions, snapshot, event) {
			continue
		}

		if err := e.executeActions(ctx, tenantID, rule, event, snapshot); err != nil {
			failures = append(failures, err)
			e.logger.Error().Err(err).
				Str("rule_id", rule.ID.String()).
				Str("rule_name", rule.Name).
				Msg("automation engine: failed to execute actions")
		}
	}

	return errors.Join(failures...)
}

// matchConditions evaluates simple equality conditions against the ticket snapshot.
// Conditions is a map of field name -> expected value (e.g. {"priority": "high"}).
func (e *Engine) matchConditions(conditions map[string]any, snapshot *TicketSnapshot, event Event) bool {
	canonical, err := domainautomation.NormalizeConditions(conditions)
	if err != nil {
		return false
	}
	for _, raw := range canonical["conditions"].([]any) {
		clause := raw.(map[string]any)
		field, expected := clause["field"].(string), clause["value"].(string)
		actual, ok := conditionActual(field, snapshot, event)
		if !ok || !clauseHolds(clause["op"], actual, expected) {
			return false
		}
	}
	return true
}

// conditionActual returns the value a condition field refers to. It reports false
// for unknown fields and for ticket fields when no snapshot is available.
func conditionActual(field string, snapshot *TicketSnapshot, event Event) (string, bool) {
	switch strings.ToLower(field) {
	case "old_value", "oldvalue":
		// Without a snapshot, only the exact snake_case spelling matches.
		return event.OldValue, snapshot != nil || field == "old_value"
	case "new_value", "newvalue":
		return event.NewValue, snapshot != nil || field == "new_value"
	case "status", "priority", "category", "title":
		if snapshot == nil {
			return "", false
		}
		return ticketAttribute(snapshot, strings.ToLower(field)), true
	default:
		// Unknown condition field — fail safe (no match).
		return "", false
	}
}

// ticketAttribute returns the snapshot attribute named by a lowercase condition field.
func ticketAttribute(snapshot *TicketSnapshot, field string) string {
	switch field {
	case "status":
		return snapshot.Status
	case "priority":
		return snapshot.Priority
	case "category":
		return snapshot.Category
	case "title":
		return snapshot.Title
	default:
		return ""
	}
}

// clauseHolds reports whether actual satisfies one condition operator.
// Unknown operators do not fail the clause.
func clauseHolds(op any, actual, expected string) bool {
	switch op {
	case "equals":
		return strings.EqualFold(actual, expected)
	case "not_equals":
		return !strings.EqualFold(actual, expected)
	case "contains":
		return strings.Contains(strings.ToLower(actual), strings.ToLower(expected))
	default:
		return true
	}
}

// executeActions runs all actions defined in a rule.
// Each action is a map[string]any with an "type" key and action-specific fields.
func (e *Engine) executeActions(ctx context.Context, tenantID uuid.UUID, rule domainautomation.Rule, event Event, snapshot *TicketSnapshot) error {
	if executor, ok := e.tickets.(interface {
		RunInTx(context.Context, func(context.Context) error) error
	}); ok {
		return executor.RunInTx(ctx, func(ctx context.Context) error { return e.runActions(ctx, tenantID, rule, event, snapshot) })
	}
	return e.runActions(ctx, tenantID, rule, event, snapshot)
}

func (e *Engine) runActions(ctx context.Context, tenantID uuid.UUID, rule domainautomation.Rule, event Event, snapshot *TicketSnapshot) error {
	var failures []error
	for _, raw := range rule.Actions {
		action, ok := raw.(map[string]any)
		if !ok {
			failures = append(failures, domainautomation.ErrValidation)
			continue
		}

		actionType, _ := action["type"].(string)
		if actionType == "" {
			failures = append(failures, domainautomation.ErrValidation)
			continue
		}

		if err := e.executeAction(ctx, tenantID, actionType, action, event, snapshot); err != nil {
			failures = append(failures, err)
			e.logger.Error().Err(err).
				Str("action", actionType).
				Str("rule_id", rule.ID.String()).
				Msg("automation engine: action failed")
			// Continue executing remaining actions.
		}
	}
	return errors.Join(failures...)
}

func (e *Engine) executeAction(ctx context.Context, tenantID uuid.UUID, actionType string, action map[string]any, event Event, snapshot *TicketSnapshot) error {
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
		return e.addCommentAction(ctx, tenantID, event, action)

	case "add_tag":
		return fmt.Errorf("%w: add_tag is not supported", domainautomation.ErrValidation)

	case "send_email":
		return e.sendEmailAction(ctx, event, snapshot, action)

	default:
		return fmt.Errorf("%w: unknown action %q", domainautomation.ErrValidation, actionType)
	}
}

// addCommentAction posts the comment described by an add_comment action.
func (e *Engine) addCommentAction(ctx context.Context, tenantID uuid.UUID, event Event, action map[string]any) error {
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
}

// sendEmailAction sends the notification described by a send_email action.
func (e *Engine) sendEmailAction(ctx context.Context, event Event, snapshot *TicketSnapshot, action map[string]any) error {
	to, _ := action["to"].(string)
	if to == "" {
		return fmt.Errorf("send_email: missing 'to' field")
	}
	if e.emailSender == nil {
		return fmt.Errorf("send_email: no email sender configured")
	}
	subject, _ := action["subject"].(string)
	if subject == "" {
		subject = "Trickreport automation notification"
	}
	body := fmt.Sprintf("Ticket ID: %s", event.TicketID)
	if snapshot != nil {
		body = fmt.Sprintf("Ticket: %s\nStatus: %s\nPriority: %s", snapshot.Title, snapshot.Status, snapshot.Priority)
	}
	if sender, ok := e.emailSender.(interface {
		SendContext(context.Context, string, string, string) error
	}); ok {
		return sender.SendContext(ctx, to, subject, body)
	}
	return e.emailSender.Send(to, subject, body)
}
