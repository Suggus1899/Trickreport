package db

import (
	"context"
	"testing"
)

func TestNew_InvalidDSN(t *testing.T) {
	_, err := New(context.Background(), "not-a-valid-dsn")
	if err == nil {
		t.Error("expected error for invalid DSN")
	}
}

func TestNew_UnreachableDB(t *testing.T) {
	_, err := New(context.Background(), "postgres://nobody:nopass@nonexistent.invalid:5432/nodb?sslmode=disable&connect_timeout=1")
	if err == nil {
		t.Error("expected error for unreachable database")
	}
}

func TestRunMigrations_InvalidURL(t *testing.T) {
	err := RunMigrations("not-a-valid-url")
	if err == nil {
		t.Error("expected error for invalid migration URL")
	}
}

func TestRollbackLastMigration_InvalidURL(t *testing.T) {
	err := RollbackLastMigration("not-a-valid-url")
	if err == nil {
		t.Error("expected error for invalid migration URL")
	}
}
