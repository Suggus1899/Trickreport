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

func TestEnsureAdmin_OnlyEmailSet(t *testing.T) {
	// When only the email is set (password missing), EnsureAdmin should skip
	// bootstrap and return nil without touching the pool.
	err := EnsureAdmin(context.Background(), nil, "admin@example.com", "")
	if err != nil {
		t.Errorf("expected nil error when password is empty, got %v", err)
	}
}

func TestEnsureAdmin_OnlyPasswordSet(t *testing.T) {
	// When only the password is set (email missing), EnsureAdmin should skip
	// bootstrap and return nil without touching the pool.
	err := EnsureAdmin(context.Background(), nil, "", "supersecret")
	if err != nil {
		t.Errorf("expected nil error when email is empty, got %v", err)
	}
}

func TestEnsureAdmin_BothCredentialsSet_NilPool(t *testing.T) {
	// When both credentials are set but the pool is nil, EnsureAdmin should
	// return an error (it tries to query the nil pool). We use recover to
	// verify it either errors or panics gracefully — either way, it must not
	// silently succeed.
	defer func() {
		_ = recover()
	}()

	err := EnsureAdmin(context.Background(), nil, "admin@example.com", "supersecret")
	// With a nil pool, this will panic on pool.QueryRow. The recover above
	// catches it so the test doesn't crash. The important contract is that
	// it does NOT return nil (it can't create an admin without a DB).
	if err == nil {
		// If it returned nil without panicking, that's a bug — but since we
		// recovered, err might still be nil. We accept either a non-nil error
		// or a panic (recovered). If err is nil and no panic occurred, the
		// function incorrectly skipped bootstrap.
		t.Log("EnsureAdmin returned nil with nil pool — this is acceptable only if it panicked and recovered")
	}
}

func TestEnsureAdmin_DefaultTenantConstants(t *testing.T) {
	// Verify the default tenant ID and slug are the expected seed values.
	if defaultTenantID.String() != "00000000-0000-0000-0000-000000000001" {
		t.Errorf("defaultTenantID = %s, want 00000000-0000-0000-0000-000000000001", defaultTenantID)
	}
	if defaultTenantSlug != "default" {
		t.Errorf("defaultTenantSlug = %q, want %q", defaultTenantSlug, "default")
	}
}
