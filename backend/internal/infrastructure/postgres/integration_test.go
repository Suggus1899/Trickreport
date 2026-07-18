package postgres

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	appArticle "github.com/trickreport/backend/internal/application/article"
	appTicket "github.com/trickreport/backend/internal/application/ticket"
	appUser "github.com/trickreport/backend/internal/application/user"
	"github.com/trickreport/backend/internal/db"
	domainarticle "github.com/trickreport/backend/internal/domain/article"
	domainautomation "github.com/trickreport/backend/internal/domain/automation"
	domainsla "github.com/trickreport/backend/internal/domain/sla"
	domainticket "github.com/trickreport/backend/internal/domain/ticket"
	domainuser "github.com/trickreport/backend/internal/domain/user"
)

// testPool is shared across all integration tests in this file. It is
// initialized once in TestMain when TEST_DATABASE_URL is set.
var testPool *pgxpool.Pool

// TestMain sets up the shared connection pool and runs migrations when
// TEST_DATABASE_URL is available. When the env var is not set, all
// integration tests are skipped.
func TestMain(m *testing.M) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		// No DB available — run the non-integration tests (the rest of the
		// package) and skip the integration suite.
		os.Exit(m.Run())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pool, err := db.New(ctx, dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "integration_test: cannot connect to TEST_DATABASE_URL: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := db.RunMigrations(dsn); err != nil {
		fmt.Fprintf(os.Stderr, "integration_test: migrations failed: %v\n", err)
		os.Exit(1)
	}

	testPool = pool
	os.Exit(m.Run())
}

// skipIfNoDB skips the calling test when no test database is configured.
func skipIfNoDB(t *testing.T) {
	t.Helper()
	if testPool == nil {
		t.Skip("TEST_DATABASE_URL not set — skipping integration test")
	}
}

// seedTenantUser inserts a fresh tenant and admin user, returning their IDs.
// Each test gets its own tenant so tests are isolated.
func seedTenantUser(t *testing.T, ctx context.Context) (uuid.UUID, uuid.UUID) {
	t.Helper()
	tenantID := uuid.New()
	_, err := testPool.Exec(ctx,
		`INSERT INTO tenants (id, name, slug, settings) VALUES ($1, $2, $3, '{}')`,
		tenantID, "Test Tenant "+tenantID.String()[:8], "t-"+tenantID.String()[:8],
	)
	if err != nil {
		t.Fatalf("seed tenant: %v", err)
	}

	userID := uuid.New()
	_, err = testPool.Exec(ctx,
		`INSERT INTO users (id, tenant_id, name, email, role) VALUES ($1, $2, $3, $4, 'admin')`,
		userID, tenantID, "Test Admin", "admin-"+tenantID.String()[:8]+"@test.local",
	)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	return tenantID, userID
}

// cleanupTenant removes the tenant and all cascaded data.
func cleanupTenant(ctx context.Context, tenantID uuid.UUID) {
	_, _ = testPool.Exec(ctx, `DELETE FROM tenants WHERE id = $1`, tenantID)
}

// ---------------------------------------------------------------------------
// Ticket CRUD
// ---------------------------------------------------------------------------

func TestIntegration_TicketRepo_CRUD(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenantID, userID := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenantID)

	repo := NewTicketRepo(testPool)

	// Create
	tk := &domainticket.Ticket{
		TenantID:    tenantID,
		Title:       "Integration test ticket",
		Description: "Something is broken",
		Priority:    domainticket.PriorityHigh,
		Category:    "bug",
		CreatedBy:   userID,
	}
	if err := repo.Create(ctx, tk); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if tk.ID == uuid.Nil {
		t.Fatal("expected non-nil ticket ID after Create")
	}
	if tk.Status != domainticket.StatusOpen {
		t.Errorf("Status = %v, want open", tk.Status)
	}

	// Get
	got, err := repo.GetByID(ctx, tk.ID, tenantID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Title != "Integration test ticket" {
		t.Errorf("Title = %q", got.Title)
	}

	// List
	list, err := repo.List(ctx, tenantID, appTicket.Filter{Limit: 10}, "admin", userID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("List returned %d tickets, want 1", len(list))
	}

	// UpdateStatus
	updated, err := repo.UpdateStatus(ctx, tk.ID, tenantID, tk.Version, domainticket.StatusInProgress, userID)
	if err != nil {
		t.Fatalf("UpdateStatus: %v", err)
	}
	if updated.Status != domainticket.StatusInProgress {
		t.Errorf("Status = %v, want in_progress", updated.Status)
	}

	// Delete (no repo Delete method; use direct SQL to clean up is fine, but
	// verify the ticket still exists then remove via cascade cleanup).
	got2, err := repo.GetByID(ctx, tk.ID, tenantID)
	if err != nil {
		t.Fatalf("GetByID after update: %v", err)
	}
	if got2.Status != domainticket.StatusInProgress {
		t.Errorf("Status after re-fetch = %v", got2.Status)
	}
}

// ---------------------------------------------------------------------------
// User CRUD
// ---------------------------------------------------------------------------

func TestIntegration_UserRepo_CRUD(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenantID, _ := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenantID)

	repo := NewUserRepo(testPool)

	// Create
	u := &domainuser.User{
		TenantID: tenantID,
		Name:     "Jane Agent",
		Email:    "jane-" + tenantID.String()[:8] + "@test.local",
		Role:     domainuser.RoleAgent,
	}
	if err := repo.Create(ctx, u, "hashedpassword"); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if u.ID == uuid.Nil {
		t.Fatal("expected non-nil user ID")
	}

	// Get
	got, err := repo.GetByID(ctx, u.ID, tenantID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Email != u.Email {
		t.Errorf("Email = %q", got.Email)
	}

	// List
	users, err := repo.List(ctx, tenantID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(users) < 2 { // seeded admin + jane
		t.Fatalf("List returned %d users, want >= 2", len(users))
	}

	// Update
	active := false
	updated, err := repo.Update(ctx, u.ID, tenantID, appUser.UpdateFields{Active: &active})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Active {
		t.Error("expected Active=false after update")
	}

	// Deactivate
	if err := repo.Deactivate(ctx, u.ID, tenantID); err != nil {
		t.Fatalf("Deactivate: %v", err)
	}
}

// ---------------------------------------------------------------------------
// Article CRUD
// ---------------------------------------------------------------------------

func TestIntegration_ArticleRepo_CRUD(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenantID, userID := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenantID)

	repo := NewArticleRepo(testPool)

	// Create
	a := &domainarticle.Article{
		TenantID:  tenantID,
		Title:     "How to reset password",
		Content:   "Go to settings and click reset.",
		Category:  "general",
		Tags:      []string{"password", "security"},
		Published: true,
		CreatedBy: userID,
	}
	if err := repo.Create(ctx, a); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if a.ID == uuid.Nil {
		t.Fatal("expected non-nil article ID")
	}

	// Get
	got, err := repo.GetByID(ctx, a.ID, tenantID, "admin")
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Title != "How to reset password" {
		t.Errorf("Title = %q", got.Title)
	}

	// List
	articles, err := repo.List(ctx, tenantID, appArticle.Filter{}, "admin")
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(articles) != 1 {
		t.Fatalf("List returned %d, want 1", len(articles))
	}

	// Update
	a.Title = "How to reset password (updated)"
	if err := repo.Update(ctx, a); err != nil {
		t.Fatalf("Update: %v", err)
	}
	got2, _ := repo.GetByID(ctx, a.ID, tenantID, "admin")
	if got2.Title != "How to reset password (updated)" {
		t.Errorf("Title after update = %q", got2.Title)
	}

	// Delete
	if err := repo.Delete(ctx, a.ID, tenantID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	if _, err := repo.GetByID(ctx, a.ID, tenantID, "admin"); err == nil {
		t.Error("expected error after delete")
	}
}

// ---------------------------------------------------------------------------
// SLA Policy CRUD
// ---------------------------------------------------------------------------

func TestIntegration_SLARepo_CRUD(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenantID, _ := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenantID)

	repo := NewSLARepo(testPool)

	// Upsert (create)
	p := &domainsla.Policy{
		TenantID:              tenantID,
		Priority:              "high",
		ResponseTimeMinutes:   30,
		ResolutionTimeMinutes: 480,
		EscalationMinutes:     60,
	}
	if err := repo.Upsert(ctx, p); err != nil {
		t.Fatalf("Upsert (create): %v", err)
	}
	if p.ID == uuid.Nil {
		t.Fatal("expected non-nil policy ID")
	}

	// List
	policies, err := repo.List(ctx, tenantID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(policies) != 1 {
		t.Fatalf("List returned %d, want 1", len(policies))
	}

	// Upsert (update)
	p.ResolutionTimeMinutes = 720
	if err := repo.Upsert(ctx, p); err != nil {
		t.Fatalf("Upsert (update): %v", err)
	}
	policies2, _ := repo.List(ctx, tenantID)
	if len(policies2) != 1 {
		t.Fatalf("List after upsert returned %d, want 1", len(policies2))
	}
	if policies2[0].ResolutionTimeMinutes != 720 {
		t.Errorf("ResolutionTimeMinutes = %d, want 720", policies2[0].ResolutionTimeMinutes)
	}
}

// ---------------------------------------------------------------------------
// Automation Rule CRUD
// ---------------------------------------------------------------------------

func TestIntegration_AutomationRepo_CRUD(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenantID, _ := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenantID)

	repo := NewAutomationRepo(testPool)

	// Create
	rule := &domainautomation.Rule{
		TenantID:    tenantID,
		Name:        "Auto-assign high priority",
		Description: "Assigns high priority tickets to agent",
		TriggerType: "ticket_created",
		Conditions:  map[string]any{"priority": "high"},
		Actions:     []any{map[string]any{"type": "set_status", "status": "in_progress"}},
		IsActive:    true,
	}
	if err := repo.Create(ctx, rule); err != nil {
		t.Fatalf("Create: %v", err)
	}
	if rule.ID == uuid.Nil {
		t.Fatal("expected non-nil rule ID")
	}

	// List
	rules, err := repo.List(ctx, tenantID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(rules) != 1 {
		t.Fatalf("List returned %d, want 1", len(rules))
	}

	// Update
	rule.Name = "Auto-assign (updated)"
	if err := repo.Update(ctx, rule); err != nil {
		t.Fatalf("Update: %v", err)
	}

	// Delete
	if err := repo.Delete(ctx, rule.ID, tenantID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	rules2, _ := repo.List(ctx, tenantID)
	if len(rules2) != 0 {
		t.Fatalf("List after delete returned %d, want 0", len(rules2))
	}
}

// ---------------------------------------------------------------------------
// Attachment CRUD
// ---------------------------------------------------------------------------

func TestIntegration_AttachmentRepo_CRUD(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenantID, userID := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenantID)

	// Create a ticket to attach to.
	ticketRepo := NewTicketRepo(testPool)
	tk := &domainticket.Ticket{
		TenantID:    tenantID,
		Title:       "Ticket with attachment",
		Description: "See attached file",
		Priority:    domainticket.PriorityLow,
		Category:    "general",
		CreatedBy:   userID,
	}
	if err := ticketRepo.Create(ctx, tk); err != nil {
		t.Fatalf("seed ticket: %v", err)
	}

	repo := NewAttachmentRepo(testPool)
	fileData := []byte("hello world attachment content")

	// Create
	att, err := repo.Create(ctx, appTicket.AttachmentUploadInput{
		TicketID:    tk.ID,
		UserID:      userID,
		Filename:    "report.txt",
		ContentType: "text/plain",
		FileSize:    int64(len(fileData)),
		FileData:    fileData,
	})
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if att.ID == uuid.Nil {
		t.Fatal("expected non-nil attachment ID")
	}

	// List
	atts, err := repo.List(ctx, tk.ID)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(atts) != 1 {
		t.Fatalf("List returned %d, want 1", len(atts))
	}
	if atts[0].Filename != "report.txt" {
		t.Errorf("Filename = %q", atts[0].Filename)
	}

	// GetByID (download)
	got, data, err := repo.GetByID(ctx, att.ID, tk.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Filename != "report.txt" {
		t.Errorf("Filename = %q", got.Filename)
	}
	if string(data) != string(fileData) {
		t.Errorf("file data mismatch")
	}
}
