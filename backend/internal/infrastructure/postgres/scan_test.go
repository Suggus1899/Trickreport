package postgres

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	domainarticle "github.com/trickreport/backend/internal/domain/article"
	domainautomation "github.com/trickreport/backend/internal/domain/automation"
)

// mockRow implements pgx.Row for testing scan functions.
type mockRow struct {
	err  error
	set  func(dest []any)
}

func (m *mockRow) Scan(dest ...any) error {
	if m.err != nil {
		return m.err
	}
	if m.set != nil {
		m.set(dest)
	}
	return nil
}

func TestScanArticle_Success(t *testing.T) {
	row := &mockRow{
		set: func(dest []any) {
			// Order: &id, &tid, &a.Title, &a.Content, &a.Category, &a.Tags, &a.Published, &createdBy, &a.CreatedAt, &a.UpdatedAt, &a.AuthorName
			*dest[0].(*pgtype.UUID) = pgtype.UUID{Valid: true, Bytes: [16]byte{1, 2, 3, 4}}
			*dest[1].(*pgtype.UUID) = pgtype.UUID{Valid: true, Bytes: [16]byte{5, 6, 7, 8}}
			*dest[2].(*string) = "Test Article"
			*dest[3].(*string) = "Content here"
			*dest[4].(*string) = "general"
			*dest[5].(*[]string) = []string{"tag1", "tag2"}
			*dest[6].(*bool) = true
			*dest[7].(*pgtype.UUID) = pgtype.UUID{Valid: true, Bytes: [16]byte{9, 10, 11, 12}}
			*dest[8].(*time.Time) = time.Now()
			*dest[9].(*time.Time) = time.Now()
			*dest[10].(*string) = "Admin User"
		},
	}

	a, err := scanArticle(row)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Title != "Test Article" {
		t.Errorf("Title = %q", a.Title)
	}
	if a.Category != "general" {
		t.Errorf("Category = %q", a.Category)
	}
	if !a.Published {
		t.Error("expected Published=true")
	}
	if a.AuthorName != "Admin User" {
		t.Errorf("AuthorName = %q", a.AuthorName)
	}
}

func TestScanArticle_NotFound(t *testing.T) {
	row := &mockRow{err: pgx.ErrNoRows}
	_, err := scanArticle(row)
	if !errors.Is(err, domainarticle.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestScanArticle_ScanError(t *testing.T) {
	row := &mockRow{err: errors.New("scan failed")}
	_, err := scanArticle(row)
	if err == nil {
		t.Error("expected error")
	}
	if errors.Is(err, domainarticle.ErrNotFound) {
		t.Error("should not be ErrNotFound")
	}
}

func TestScanArticleRow_Success(t *testing.T) {
	row := &mockRow{
		set: func(dest []any) {
			*dest[0].(*pgtype.UUID) = pgtype.UUID{Valid: true, Bytes: [16]byte{1}}
			*dest[1].(*pgtype.UUID) = pgtype.UUID{Valid: true, Bytes: [16]byte{2}}
			*dest[2].(*string) = "Row Article"
			*dest[3].(*string) = "Content"
			*dest[4].(*string) = "cat"
			*dest[5].(*[]string) = nil
			*dest[6].(*bool) = false
			*dest[7].(*pgtype.UUID) = pgtype.UUID{Valid: false}
			*dest[8].(*time.Time) = time.Time{}
			*dest[9].(*time.Time) = time.Time{}
			*dest[10].(*string) = ""
		},
	}

	a, err := scanArticleRow(row)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Title != "Row Article" {
		t.Errorf("Title = %q", a.Title)
	}
}

func TestScanRule_Success(t *testing.T) {
	row := &mockRow{
		set: func(dest []any) {
			// Order: &id, &tid, &rl.Name, &rl.Description, &rl.TriggerType, &rl.Conditions, &rl.Actions, &rl.IsActive, &rl.CreatedAt, &rl.UpdatedAt
			*dest[0].(*pgtype.UUID) = pgtype.UUID{Valid: true, Bytes: [16]byte{1, 2, 3}}
			*dest[1].(*pgtype.UUID) = pgtype.UUID{Valid: true, Bytes: [16]byte{4, 5, 6}}
			*dest[2].(*string) = "Auto Rule"
			*dest[3].(*string) = "Description"
			*dest[4].(*string) = "on_create"
			*dest[5].(*map[string]any) = map[string]any{"key": "value"}
			*dest[6].(*[]any) = []any{map[string]any{"action": "set_priority"}}
			*dest[7].(*bool) = true
			*dest[8].(*time.Time) = time.Now()
			*dest[9].(*time.Time) = time.Now()
		},
	}

	r, err := scanRule(row)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Name != "Auto Rule" {
		t.Errorf("Name = %q", r.Name)
	}
	if r.TriggerType != "on_create" {
		t.Errorf("TriggerType = %q", r.TriggerType)
	}
	if !r.IsActive {
		t.Error("expected IsActive=true")
	}
}

func TestScanRule_NotFound(t *testing.T) {
	row := &mockRow{err: pgx.ErrNoRows}
	_, err := scanRule(row)
	if !errors.Is(err, domainautomation.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestScanRule_ScanError(t *testing.T) {
	row := &mockRow{err: errors.New("scan failed")}
	_, err := scanRule(row)
	if err == nil {
		t.Error("expected error")
	}
	if errors.Is(err, domainautomation.ErrNotFound) {
		t.Error("should not be ErrNotFound")
	}
}
