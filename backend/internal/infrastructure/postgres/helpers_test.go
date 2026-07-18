package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"

	domainarticle "github.com/trickreport/backend/internal/domain/article"
	domainautomation "github.com/trickreport/backend/internal/domain/automation"
	domainsla "github.com/trickreport/backend/internal/domain/sla"
	domainticket "github.com/trickreport/backend/internal/domain/ticket"
	domainuser "github.com/trickreport/backend/internal/domain/user"
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

// ---------------------------------------------------------------------------
// Mock DBTX — for testing repository methods without a real database.
// ---------------------------------------------------------------------------

// mockDBTX implements the DBTX interface for unit testing repos.
type mockDBTX struct {
	execResult  pgconn.CommandTag
	execErr     error
	queryErr    error
	queryRowRow pgx.Row
	rows        pgx.Rows
}

func (m *mockDBTX) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	return m.execResult, m.execErr
}

func (m *mockDBTX) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	return m.rows, m.queryErr
}

func (m *mockDBTX) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return m.queryRowRow
}

// mockRows implements pgx.Rows for testing List queries.
type mockRows struct {
	rows  [][]any
	idx   int
	closed bool
}

func newMockRows(data [][]any) *mockRows {
	return &mockRows{rows: data}
}

func (m *mockRows) Next() bool {
	if m.closed {
		return false
	}
	m.idx++
	return m.idx <= len(m.rows)
}

func (m *mockRows) Scan(dest ...any) error {
	if m.idx < 1 || m.idx > len(m.rows) {
		return errors.New("no rows to scan")
	}
	row := m.rows[m.idx-1]
	for i, val := range row {
		if i >= len(dest) {
			break
		}
		switch d := dest[i].(type) {
		case *pgtype.UUID:
			if v, ok := val.(pgtype.UUID); ok {
				*d = v
			}
		case *string:
			if v, ok := val.(string); ok {
				*d = v
			}
		case *[]string:
			if v, ok := val.([]string); ok {
				*d = v
			}
		case *bool:
			if v, ok := val.(bool); ok {
				*d = v
			}
		case *int:
			if v, ok := val.(int); ok {
				*d = v
			}
		case *time.Time:
			if v, ok := val.(time.Time); ok {
				*d = v
			}
		case *map[string]any:
			if v, ok := val.(map[string]any); ok {
				*d = v
			}
		case *[]any:
			if v, ok := val.([]any); ok {
				*d = v
			}
		}
	}
	return nil
}

func (m *mockRows) Close() { m.closed = true }
func (m *mockRows) Err() error { return nil }
func (m *mockRows) CommandTag() pgconn.CommandTag { return pgconn.NewCommandTag("") }
func (m *mockRows) FieldDescriptions() []pgconn.FieldDescription { return nil }
func (m *mockRows) Values() ([]any, error) { return nil, nil }
func (m *mockRows) RawValues() [][]byte { return nil }
func (m *mockRows) Conn() *pgx.Conn { return nil }

// ---------------------------------------------------------------------------
// Test entity builders — for creating domain objects in tests.
// ---------------------------------------------------------------------------

// newTestTicket returns a fully populated ticket for use in tests.
func newTestTicket() *domainticket.Ticket {
	tenantID := uuid.New()
	creatorID := uuid.New()
	return &domainticket.Ticket{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Title:       "Test ticket",
		Description: "Test description",
		Status:      domainticket.StatusOpen,
		Priority:    domainticket.PriorityMedium,
		Category:    "general",
		CreatedBy:   creatorID,
		Version:     0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// newTestUser returns a fully populated user for use in tests.
func newTestUser() *domainuser.User {
	return &domainuser.User{
		ID:        uuid.New(),
		TenantID:  uuid.New(),
		Name:      "Test User",
		Email:     "test@example.com",
		Role:      domainuser.RoleAgent,
		Active:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// newTestArticle returns a fully populated article for use in tests.
func newTestArticle() *domainarticle.Article {
	return &domainarticle.Article{
		ID:        uuid.New(),
		TenantID:  uuid.New(),
		Title:     "Test article",
		Content:   "Test content",
		Category:  "general",
		Tags:      []string{"test"},
		Published: true,
		CreatedBy: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// newTestSLAPolicy returns a fully populated SLA policy for use in tests.
func newTestSLAPolicy() *domainsla.Policy {
	return &domainsla.Policy{
		ID:                    uuid.New(),
		TenantID:              uuid.New(),
		Priority:              "high",
		ResponseTimeMinutes:   60,
		ResolutionTimeMinutes: 480,
		EscalationMinutes:     120,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}
}

// newTestRule returns a fully populated automation rule for use in tests.
func newTestRule() *domainautomation.Rule {
	return &domainautomation.Rule{
		ID:          uuid.New(),
		TenantID:    uuid.New(),
		Name:        "Test rule",
		Description: "Test description",
		TriggerType: "ticket_created",
		Conditions:  map[string]any{"priority": "high"},
		Actions:     []any{map[string]any{"type": "set_status", "status": "in_progress"}},
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ---------------------------------------------------------------------------
// Mock DBTX tests
// ---------------------------------------------------------------------------

func TestMockDBTX_Exec(t *testing.T) {
	db := &mockDBTX{execErr: errors.New("exec failed")}
	_, err := db.Exec(context.Background(), "INSERT INTO t VALUES (1)")
	if err == nil || err.Error() != "exec failed" {
		t.Errorf("expected exec failed error, got %v", err)
	}
}

func TestMockDBTX_QueryRow(t *testing.T) {
	row := &mockRow{set: func(dest []any) {
		*dest[0].(*string) = "hello"
	}}
	db := &mockDBTX{queryRowRow: row}
	var s string
	err := db.QueryRow(context.Background(), "SELECT 1").Scan(&s)
	if err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if s != "hello" {
		t.Errorf("got %q, want hello", s)
	}
}

func TestMockDBTX_Query(t *testing.T) {
	rows := newMockRows([][]any{
		{"row1"},
		{"row2"},
	})
	db := &mockDBTX{rows: rows}

	result, err := db.Query(context.Background(), "SELECT name FROM t")
	if err != nil {
		t.Fatalf("Query: %v", err)
	}
	defer result.Close()

	count := 0
	for result.Next() {
		var name string
		if err := result.Scan(&name); err != nil {
			t.Fatalf("Scan: %v", err)
		}
		count++
	}
	if count != 2 {
		t.Errorf("scanned %d rows, want 2", count)
	}
}

func TestMockRows_Close(t *testing.T) {
	rows := newMockRows([][]any{{"a"}})
	rows.Close()
	if rows.Next() {
		t.Error("Next should return false after Close")
	}
}

func TestMockRows_Empty(t *testing.T) {
	rows := newMockRows(nil)
	if rows.Next() {
		t.Error("Next should return false for empty rows")
	}
}

func TestNewTestTicket(t *testing.T) {
	tk := newTestTicket()
	if tk.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
	if tk.Title != "Test ticket" {
		t.Errorf("Title = %q", tk.Title)
	}
	if tk.Status != domainticket.StatusOpen {
		t.Errorf("Status = %v", tk.Status)
	}
}

func TestNewTestUser(t *testing.T) {
	u := newTestUser()
	if u.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
	if u.Role != domainuser.RoleAgent {
		t.Errorf("Role = %v", u.Role)
	}
	if !u.Active {
		t.Error("expected Active=true")
	}
}

func TestNewTestArticle(t *testing.T) {
	a := newTestArticle()
	if a.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
	if !a.Published {
		t.Error("expected Published=true")
	}
	if len(a.Tags) != 1 {
		t.Errorf("Tags len = %d, want 1", len(a.Tags))
	}
}

func TestNewTestSLAPolicy(t *testing.T) {
	p := newTestSLAPolicy()
	if p.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
	if p.ResolutionTimeMinutes != 480 {
		t.Errorf("ResolutionTimeMinutes = %d, want 480", p.ResolutionTimeMinutes)
	}
}

func TestNewTestRule(t *testing.T) {
	r := newTestRule()
	if r.ID == uuid.Nil {
		t.Error("expected non-nil ID")
	}
	if r.TriggerType != "ticket_created" {
		t.Errorf("TriggerType = %q", r.TriggerType)
	}
	if !r.IsActive {
		t.Error("expected IsActive=true")
	}
}
