package automation

import (
	"context"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/domain/automation"
)

// Repository is the port for automation rule persistence.
type Repository interface {
	List(ctx context.Context, tenantID uuid.UUID) ([]automation.Rule, error)
	Create(ctx context.Context, r *automation.Rule) error
	Update(ctx context.Context, r *automation.Rule) error
	Delete(ctx context.Context, id, tenantID uuid.UUID) error
}

// Service is the application service for automation operations.
type Service struct {
	repo Repository
}

// NewService creates a new automation application service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// List returns all automation rules for the tenant.
func (s *Service) List(ctx context.Context, tenantID uuid.UUID) ([]automation.Rule, error) {
	rules, err := s.repo.List(ctx, tenantID)
	if err != nil {
		return nil, err
	}
	for i := range rules {
		canonical, err := automation.NormalizeConditions(rules[i].Conditions)
		if err != nil {
			return nil, err
		}
		rules[i].Conditions = canonical
	}
	return rules, nil
}

// CreateInput holds the data for creating a new automation rule.
type CreateInput struct {
	TenantID    uuid.UUID
	Name        string
	Description string
	TriggerType string
	Conditions  map[string]any
	Actions     []any
	IsActive    bool
}

// Create creates a new automation rule.
func (s *Service) Create(ctx context.Context, input CreateInput) (*automation.Rule, error) {
	conditions, err := automation.ValidateRule(input.Name, input.TriggerType, input.Conditions, input.Actions)
	if err != nil {
		return nil, err
	}

	r := &automation.Rule{
		TenantID:    input.TenantID,
		Name:        input.Name,
		Description: input.Description,
		TriggerType: input.TriggerType,
		Conditions:  conditions,
		Actions:     input.Actions,
		IsActive:    input.IsActive,
	}

	if err := s.repo.Create(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

// UpdateInput holds the data for updating an automation rule.
type UpdateInput struct {
	ID          uuid.UUID
	TenantID    uuid.UUID
	Name        string
	Description string
	TriggerType string
	Conditions  map[string]any
	Actions     []any
	IsActive    bool
}

// Update modifies an existing automation rule.
func (s *Service) Update(ctx context.Context, input UpdateInput) (*automation.Rule, error) {
	conditions, err := automation.ValidateRule(input.Name, input.TriggerType, input.Conditions, input.Actions)
	if err != nil {
		return nil, err
	}

	r := &automation.Rule{
		ID:          input.ID,
		TenantID:    input.TenantID,
		Name:        input.Name,
		Description: input.Description,
		TriggerType: input.TriggerType,
		Conditions:  conditions,
		Actions:     input.Actions,
		IsActive:    input.IsActive,
	}

	if err := s.repo.Update(ctx, r); err != nil {
		return nil, err
	}
	return r, nil
}

// Delete removes an automation rule.
func (s *Service) Delete(ctx context.Context, id, tenantID uuid.UUID) error {
	return s.repo.Delete(ctx, id, tenantID)
}
