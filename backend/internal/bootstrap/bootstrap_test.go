package bootstrap

import (
	"context"
	"testing"
)

func TestEnsureAdmin_EmptyCredentials(t *testing.T) {
	// With empty credentials, EnsureAdmin should return nil without touching the pool.
	err := EnsureAdmin(context.Background(), nil, "", "")
	if err != nil {
		t.Errorf("expected nil error for empty credentials, got %v", err)
	}
}
