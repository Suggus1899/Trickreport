package postgres

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func TestPgToUUID_Valid(t *testing.T) {
	original := uuid.New()
	pgUUID := pgtype.UUID{Valid: true}
	copy(pgUUID.Bytes[:], original[:])

	got := pgToUUID(pgUUID)
	if got != original {
		t.Errorf("pgToUUID = %v, want %v", got, original)
	}
}

func TestPgToUUID_Invalid(t *testing.T) {
	pgUUID := pgtype.UUID{Valid: false}
	got := pgToUUID(pgUUID)
	if got != uuid.Nil {
		t.Errorf("pgToUUID = %v, want uuid.Nil", got)
	}
}

func TestPgToUUIDPtr_Valid(t *testing.T) {
	original := uuid.New()
	pgUUID := pgtype.UUID{Valid: true}
	copy(pgUUID.Bytes[:], original[:])

	got := pgToUUIDPtr(pgUUID)
	if got == nil {
		t.Fatal("expected non-nil pointer")
	}
	if *got != original {
		t.Errorf("pgToUUIDPtr = %v, want %v", *got, original)
	}
}

func TestPgToUUIDPtr_Invalid(t *testing.T) {
	pgUUID := pgtype.UUID{Valid: false}
	got := pgToUUIDPtr(pgUUID)
	if got != nil {
		t.Errorf("expected nil, got %v", *got)
	}
}
