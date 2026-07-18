package sla

import (
	"context"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/domain/sla"
)

// Repository is the port for SLA policy persistence.
type Repository interface {
	List(ctx context.Context, tenantID uuid.UUID) ([]sla.Policy, error)
	Upsert(ctx context.Context, p *sla.Policy) error
}

// Service is the application service for SLA operations.
type Service struct {
	repo Repository
}

// NewService creates a new SLA application service.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// List returns all SLA policies for the tenant.
func (s *Service) List(ctx context.Context, tenantID uuid.UUID) ([]sla.Policy, error) {
	return s.repo.List(ctx, tenantID)
}

// UpsertInput holds the data for creating or updating an SLA policy.
type UpsertInput struct {
	TenantID              uuid.UUID
	Priority              string
	ResponseTimeMinutes   int
	ResolutionTimeMinutes int
	EscalationMinutes     int
}

// Upsert creates or updates an SLA policy for a specific priority.
func (s *Service) Upsert(ctx context.Context, input UpsertInput) (*sla.Policy, error) {
	switch input.Priority {
	case "low", "medium", "high", "critical":
	default:
		return nil, sla.ErrInvalidPriority
	}

	p := &sla.Policy{
		TenantID:              input.TenantID,
		Priority:              input.Priority,
		ResponseTimeMinutes:   input.ResponseTimeMinutes,
		ResolutionTimeMinutes: input.ResolutionTimeMinutes,
		EscalationMinutes:     input.EscalationMinutes,
	}

	if err := s.repo.Upsert(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}
