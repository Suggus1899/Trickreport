package sla

import (
	"context"
	"errors"
	domain "github.com/trickreport/backend/internal/domain/sla"
	"testing"
)

func TestService_ValidatesDurations(t *testing.T) {
	for _, durations := range [][3]int{{0, 60, 0}, {1, 0, 0}, {1, 60, -1}, {-1, 60, 1}} {
		repo := &mockRepo{}
		_, err := NewService(repo).Upsert(context.Background(), UpsertInput{Priority: "high", ResponseTimeMinutes: durations[0], ResolutionTimeMinutes: durations[1], EscalationMinutes: durations[2]})
		if !errors.Is(err, domain.ErrValidation) || repo.upsertCalls != 0 {
			t.Errorf("durations=%v err=%v writes=%d", durations, err, repo.upsertCalls)
		}
	}
	if _, err := NewService(&mockRepo{}).Upsert(context.Background(), UpsertInput{Priority: "high", ResponseTimeMinutes: 1, ResolutionTimeMinutes: 60, EscalationMinutes: 0}); err != nil {
		t.Fatalf("zero escalation must disable: %v", err)
	}
}
