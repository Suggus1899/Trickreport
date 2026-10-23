package automation

import (
	"fmt"
	"net/mail"
	"sort"
	"strings"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/domain/ticket"
)

// NormalizeConditions reads legacy equality maps and returns the canonical AND clauses.
func NormalizeConditions(input map[string]any) (map[string]any, error) {
	clauses := make([]any, 0, len(input))
	if raw, ok := input["conditions"]; ok {
		var valid bool
		clauses, valid = raw.([]any)
		if !valid || len(input) != 1 {
			return nil, ErrValidation
		}
	} else {
		keys := make([]string, 0, len(input))
		for key := range input {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			clauses = append(clauses, map[string]any{"field": key, "op": "equals", "value": input[key]})
		}
	}
	result := make([]any, 0, len(clauses))
	for _, raw := range clauses {
		clause, ok := raw.(map[string]any)
		if !ok || len(clause) != 3 {
			return nil, ErrValidation
		}
		field, fieldOK := clause["field"].(string)
		op, opOK := clause["op"].(string)
		value, valueOK := clause["value"].(string)
		if !fieldOK || !opOK || !valueOK {
			return nil, ErrValidation
		}
		field = strings.ToLower(field)
		if field == "oldvalue" {
			field = "old_value"
		}
		if field == "newvalue" {
			field = "new_value"
		}
		switch field {
		case "status", "priority", "category", "title", "old_value", "new_value":
		default:
			return nil, ErrValidation
		}
		switch op {
		case "equals", "not_equals", "contains":
		default:
			return nil, ErrValidation
		}
		result = append(result, map[string]any{"field": field, "op": op, "value": value})
	}
	return map[string]any{"conditions": result}, nil
}

// ValidateRule rejects contracts the engine cannot execute.
func ValidateRule(name, trigger string, conditions map[string]any, actions []any) (map[string]any, error) {
	if strings.TrimSpace(name) == "" {
		return nil, ErrValidation
	}
	switch trigger {
	case "ticket_created", "status_changed", "priority_changed", "sla_breach", "escalation":
	default:
		return nil, ErrValidation
	}
	canonical, err := NormalizeConditions(conditions)
	if err != nil {
		return nil, err
	}
	for _, raw := range actions {
		action, ok := raw.(map[string]any)
		if !ok {
			return nil, ErrValidation
		}
		kind, _ := action["type"].(string)
		allowed := map[string]bool{"type": true}
		var field string
		switch kind {
		case "set_priority":
			field = "priority"
		case "set_status":
			field = "status"
		case "assign_to":
			field = "user_id"
		case "add_comment":
			field = "content"
			allowed["is_internal"] = true
		case "send_email":
			field = "to"
			allowed["subject"] = true
		default:
			return nil, fmt.Errorf("%w: unsupported action %q", ErrValidation, kind)
		}
		allowed[field] = true
		if kind == "set_priority" || kind == "set_status" || kind == "add_comment" {
			allowed["value"] = true
		}
		for key := range action {
			if !allowed[key] {
				return nil, ErrValidation
			}
		}
		value, _ := action[field].(string)
		if value == "" && allowed["value"] {
			value, _ = action["value"].(string)
		}
		if strings.TrimSpace(value) == "" {
			return nil, ErrValidation
		}
		switch kind {
		case "set_priority":
			if _, err := ticket.ParsePriority(value); err != nil {
				return nil, ErrValidation
			}
		case "set_status":
			if _, err := ticket.ParseStatus(value); err != nil {
				return nil, ErrValidation
			}
		case "assign_to":
			if id, err := uuid.Parse(value); err != nil || id == uuid.Nil {
				return nil, ErrValidation
			}
		case "add_comment":
			if raw, ok := action["is_internal"]; ok {
				if _, ok := raw.(bool); !ok {
					return nil, ErrValidation
				}
			}
		case "send_email":
			if _, err := mail.ParseAddress(value); err != nil {
				return nil, ErrValidation
			}
			if subject, ok := action["subject"]; ok {
				if _, ok := subject.(string); !ok {
					return nil, ErrValidation
				}
			}
		}
	}
	return canonical, nil
}
