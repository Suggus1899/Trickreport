package sla

import (
	"errors"
	"testing"
	"time"
)

func TestNewDeadlines(t *testing.T) {
	created := time.Date(2026, 10, 25, 12, 0, 0, 0, time.UTC)
	at := func(minutes int) *time.Time {
		v := created.Add(time.Duration(minutes) * time.Minute)
		return &v
	}

	tests := []struct {
		name                                         string
		response, resolution, escalation             int
		wantResponse, wantResolution, wantEscalation *time.Time
	}{
		{
			name: "all clocks", response: 30, resolution: 240, escalation: 120,
			wantResponse: at(30), wantResolution: at(240), wantEscalation: at(120),
		},
		{
			name: "escalation clock only when set", response: 30, resolution: 240, escalation: 0,
			wantResponse: at(30), wantResolution: at(240),
		},
		{
			name: "zero minutes mean no SLA", response: 0, resolution: 0, escalation: 0,
		},
		{
			name: "zero response leaves that clock unset", response: 0, resolution: 240, escalation: 0,
			wantResolution: at(240),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDeadlines(created, tt.response, tt.resolution, tt.escalation)
			assertSameInstant(t, "response", got.Response, tt.wantResponse)
			assertSameInstant(t, "resolution", got.Resolution, tt.wantResolution)
			assertSameInstant(t, "escalation", got.Escalation, tt.wantEscalation)
		})
	}
}

func TestPolicyValidate(t *testing.T) {
	tests := []struct {
		name    string
		policy  Policy
		wantErr bool
	}{
		{name: "complete policy", policy: Policy{ResponseTimeMinutes: 30, ResolutionTimeMinutes: 240, EscalationMinutes: 120}},
		{name: "escalation may be unset", policy: Policy{ResponseTimeMinutes: 30, ResolutionTimeMinutes: 240}},
		{name: "response must be positive", policy: Policy{ResolutionTimeMinutes: 240}, wantErr: true},
		{name: "resolution must be positive", policy: Policy{ResponseTimeMinutes: 30}, wantErr: true},
		{name: "escalation cannot be negative", policy: Policy{ResponseTimeMinutes: 30, ResolutionTimeMinutes: 240, EscalationMinutes: -1}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.policy.Validate()
			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrValidation) {
				t.Fatalf("Validate() = %v, want ErrValidation", err)
			}
		})
	}
}

func assertSameInstant(t *testing.T, name string, got, want *time.Time) {
	t.Helper()
	if (got == nil) != (want == nil) {
		t.Fatalf("%s: got %v, want %v", name, got, want)
	}
	if got != nil && !got.Equal(*want) {
		t.Fatalf("%s: got %s, want %s", name, got, want)
	}
}
