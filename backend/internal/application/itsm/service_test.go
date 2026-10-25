package itsm

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	domain "github.com/trickreport/backend/internal/domain/itsm"
	domainSLA "github.com/trickreport/backend/internal/domain/sla"
)

// fakeRepo implements the calls RequestService makes. Any other Repository
// method panics through the nil embedded interface, which keeps the test honest
// about what the catalog request path touches.
type fakeRepo struct {
	Repository
	service   domain.CatalogService
	tickets   []fakeTicket
	approvals []domain.Approval
}

type fakeTicket struct {
	priority  string
	deadlines domainSLA.Deadlines
}

func (f *fakeRepo) GetService(context.Context, uuid.UUID, uuid.UUID) (*domain.CatalogService, error) {
	svc := f.service
	return &svc, nil
}

func (f *fakeRepo) UserExists(context.Context, uuid.UUID, uuid.UUID) (bool, error) {
	return true, nil
}

func (f *fakeRepo) CreateTicketFromService(_ context.Context, _, _ uuid.UUID, _, _, _, priority string, dl domainSLA.Deadlines) (uuid.UUID, error) {
	f.tickets = append(f.tickets, fakeTicket{priority: priority, deadlines: dl})
	return uuid.New(), nil
}

func (f *fakeRepo) CreateApproval(_ context.Context, a *domain.Approval) error {
	f.approvals = append(f.approvals, *a)
	return nil
}

type fakePolicies struct {
	policy   *domainSLA.Policy
	err      error
	priority string
}

func (f *fakePolicies) GetByPriority(_ context.Context, _ uuid.UUID, priority string) (*domainSLA.Policy, error) {
	f.priority = priority
	return f.policy, f.err
}

func requestCatalogService(t *testing.T, policies *fakePolicies, priority string) (*fakeRepo, error) {
	t.Helper()
	repo := &fakeRepo{service: domain.CatalogService{ID: uuid.New(), TenantID: uuid.New(), CategoryName: "Access", Name: "VPN", Active: true}}
	svc := NewService(repo, policies)
	_, err := svc.RequestService(context.Background(), repo.service.TenantID, uuid.New(), repo.service.ID,
		CatalogRequest{Title: "Need VPN access", Priority: priority})
	return repo, err
}

func TestRequestServiceAppliesPriorityPolicyDeadlines(t *testing.T) {
	policies := &fakePolicies{policy: &domainSLA.Policy{ResponseTimeMinutes: 30, ResolutionTimeMinutes: 240, EscalationMinutes: 120}}

	before := time.Now()
	repo, err := requestCatalogService(t, policies, "high")
	after := time.Now()
	if err != nil {
		t.Fatalf("RequestService: %v", err)
	}

	if policies.priority != "high" {
		t.Fatalf("policy looked up for %q, want high", policies.priority)
	}
	if len(repo.tickets) != 1 {
		t.Fatalf("tickets created = %d, want 1", len(repo.tickets))
	}
	dl := repo.tickets[0].deadlines
	assertWithin(t, "response", dl.Response, before.Add(30*time.Minute), after.Add(30*time.Minute))
	assertWithin(t, "resolution", dl.Resolution, before.Add(240*time.Minute), after.Add(240*time.Minute))
	assertWithin(t, "escalation", dl.Escalation, before.Add(120*time.Minute), after.Add(120*time.Minute))
}

func TestRequestServiceWithoutPolicyHasNoClocks(t *testing.T) {
	policies := &fakePolicies{err: domainSLA.ErrNotFound}

	repo, err := requestCatalogService(t, policies, "")
	if err != nil {
		t.Fatalf("RequestService: %v", err)
	}

	if policies.priority != "medium" {
		t.Fatalf("policy looked up for %q, want the default medium", policies.priority)
	}
	if len(repo.tickets) != 1 {
		t.Fatalf("tickets created = %d, want 1", len(repo.tickets))
	}
	dl := repo.tickets[0].deadlines
	if dl.Response != nil || dl.Resolution != nil || dl.Escalation != nil {
		t.Fatalf("deadlines = %+v, want none when no policy exists", dl)
	}
}

func TestRequestServiceRejectsInvalidPolicyBeforeCreating(t *testing.T) {
	policies := &fakePolicies{policy: &domainSLA.Policy{ResponseTimeMinutes: 0, ResolutionTimeMinutes: 240}}

	repo, err := requestCatalogService(t, policies, "high")
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("RequestService error = %v, want a validation error", err)
	}
	if len(repo.tickets) != 0 {
		t.Fatalf("tickets created = %d, want 0 for an invalid policy", len(repo.tickets))
	}
}

func TestRequestServicePropagatesPolicyLookupFailure(t *testing.T) {
	boom := errors.New("policy store down")
	policies := &fakePolicies{err: boom}

	repo, err := requestCatalogService(t, policies, "high")
	if !errors.Is(err, boom) {
		t.Fatalf("RequestService error = %v, want %v", err, boom)
	}
	if len(repo.tickets) != 0 {
		t.Fatalf("tickets created = %d, want 0 when the policy lookup fails", len(repo.tickets))
	}
}

func assertWithin(t *testing.T, name string, got *time.Time, lo, hi time.Time) {
	t.Helper()
	if got == nil {
		t.Fatalf("%s deadline is nil, want one between %s and %s", name, lo, hi)
	}
	if got.Before(lo) || got.After(hi) {
		t.Fatalf("%s deadline = %s, want between %s and %s", name, got, lo, hi)
	}
}
