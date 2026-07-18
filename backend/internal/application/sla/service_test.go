package sla

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/trickreport/backend/internal/domain/sla"
)

// --- Mocks ---

type mockRepo struct {
	policies   []sla.Policy
	policy     *sla.Policy
	err        error
	upsertErr  error
	gotUpsert  *sla.Policy
	upsertCalls int
}

func (m *mockRepo) List(ctx context.Context, tenantID uuid.UUID) ([]sla.Policy, error) {
	return m.policies, m.err
}

func (m *mockRepo) Upsert(ctx context.Context, p *sla.Policy) error {
	m.upsertCalls++
	m.gotUpsert = p
	return m.upsertErr
}

// --- List tests ---

func TestService_List_Delegates(t *testing.T) {
	want := []sla.Policy{{ID: uuid.New(), Priority: "high"}}
	repo := &mockRepo{policies: want}
	svc := NewService(repo)

	got, err := svc.List(context.Background(), uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("got %d", len(got))
	}
}

func TestService_List_Error(t *testing.T) {
	repo := &mockRepo{err: errors.New("db")}
	svc := NewService(repo)
	if _, err := svc.List(context.Background(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}

// --- Upsert tests ---

func TestService_Upsert_Success(t *testing.T) {
	tenant := uuid.New()
	repo := &mockRepo{}
	svc := NewService(repo)

	p, err := svc.Upsert(context.Background(), UpsertInput{
		TenantID: tenant, Priority: "high",
		ResponseTimeMinutes: 60, ResolutionTimeMinutes: 480, EscalationMinutes: 720,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Priority != "high" || p.TenantID != tenant {
		t.Errorf("policy = %+v", p)
	}
	if p.ResponseTimeMinutes != 60 || p.ResolutionTimeMinutes != 480 || p.EscalationMinutes != 720 {
		t.Errorf("times = %+v", p)
	}
	if repo.upsertCalls != 1 {
		t.Errorf("upsert calls = %d", repo.upsertCalls)
	}
}

func TestService_Upsert_AllValidPriorities(t *testing.T) {
	for _, pr := range []string{"low", "medium", "high", "critical"} {
		t.Run(pr, func(t *testing.T) {
			repo := &mockRepo{}
			svc := NewService(repo)
			if _, err := svc.Upsert(context.Background(), UpsertInput{Priority: pr}); err != nil {
				t.Fatalf("priority %q: unexpected error %v", pr, err)
			}
		})
	}
}

func TestService_Upsert_InvalidPriority(t *testing.T) {
	svc := NewService(&mockRepo{})
	_, err := svc.Upsert(context.Background(), UpsertInput{Priority: "urgent"})
	if !errors.Is(err, sla.ErrInvalidPriority) {
		t.Errorf("expected ErrInvalidPriority, got %v", err)
	}
}

func TestService_Upsert_RepoError(t *testing.T) {
	repo := &mockRepo{upsertErr: errors.New("db")}
	svc := NewService(repo)
	_, err := svc.Upsert(context.Background(), UpsertInput{Priority: "high"})
	if err == nil {
		t.Fatal("expected error")
	}
}
