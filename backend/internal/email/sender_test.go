package email

import "testing"

func TestConsoleSender_Send(t *testing.T) {
	s := NewConsoleSender()
	if s == nil {
		t.Fatal("expected non-nil sender")
	}
	if err := s.Send("user@example.com", "Test Subject", "Test body"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
