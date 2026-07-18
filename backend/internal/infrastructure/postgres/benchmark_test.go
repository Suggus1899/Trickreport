package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"

	appTicket "github.com/trickreport/backend/internal/application/ticket"
	domainticket "github.com/trickreport/backend/internal/domain/ticket"
	domainuser "github.com/trickreport/backend/internal/domain/user"
)

// benchTenantUser creates a tenant + user for benchmark isolation.
func benchTenantUser(b *testing.B, ctx context.Context) (uuid.UUID, uuid.UUID) {
	b.Helper()
	tenantID := uuid.New()
	if _, err := testPool.Exec(ctx,
		`INSERT INTO tenants (id, name, slug, settings) VALUES ($1, $2, $3, '{}')`,
		tenantID, "Bench Tenant "+tenantID.String()[:8], "b-"+tenantID.String()[:8],
	); err != nil {
		b.Fatalf("seed tenant: %v", err)
	}
	userID := uuid.New()
	if _, err := testPool.Exec(ctx,
		`INSERT INTO users (id, tenant_id, name, email, role) VALUES ($1, $2, $3, $4, 'admin')`,
		userID, tenantID, "Bench User", "bench-"+tenantID.String()[:8]+"@test.local",
	); err != nil {
		b.Fatalf("seed user: %v", err)
	}
	return tenantID, userID
}

// skipBenchIfNoDB skips the benchmark when no test database is configured.
func skipBenchIfNoDB(b *testing.B) {
	b.Helper()
	if testPool == nil {
		b.Skip("TEST_DATABASE_URL not set — skipping benchmark")
	}
}

// BenchmarkTicketRepo_Create measures the cost of inserting a single ticket.
func BenchmarkTicketRepo_Create(b *testing.B) {
	skipBenchIfNoDB(b)
	ctx := context.Background()
	tenantID, userID := benchTenantUser(b, ctx)
	defer func() {
		_, _ = testPool.Exec(ctx, `DELETE FROM tenants WHERE id = $1`, tenantID)
	}()

	repo := NewTicketRepo(testPool)
	tk := &domainticket.Ticket{
		TenantID:    tenantID,
		Title:       "Bench ticket",
		Description: "Benchmark description",
		Priority:    domainticket.PriorityMedium,
		Category:    "general",
		CreatedBy:   userID,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tk.ID = uuid.Nil
		if err := repo.Create(ctx, tk); err != nil {
			b.Fatalf("Create: %v", err)
		}
	}
	b.StopTimer()
}

// BenchmarkTicketRepo_List measures the cost of listing tickets for a tenant.
func BenchmarkTicketRepo_List(b *testing.B) {
	skipBenchIfNoDB(b)
	ctx := context.Background()
	tenantID, userID := benchTenantUser(b, ctx)
	defer func() {
		_, _ = testPool.Exec(ctx, `DELETE FROM tenants WHERE id = $1`, tenantID)
	}()

	repo := NewTicketRepo(testPool)
	// Seed 50 tickets so the list query has real work to do.
	for i := 0; i < 50; i++ {
		tk := &domainticket.Ticket{
			TenantID:    tenantID,
			Title:       fmt.Sprintf("Seed ticket %d", i),
			Description: "Seed description",
			Priority:    domainticket.PriorityMedium,
			Category:    "general",
			CreatedBy:   userID,
		}
		if err := repo.Create(ctx, tk); err != nil {
			b.Fatalf("seed Create: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := repo.List(ctx, tenantID, appTicket.Filter{Limit: 50}, "admin", userID); err != nil {
			b.Fatalf("List: %v", err)
		}
	}
	b.StopTimer()
}

// BenchmarkTicketRepo_UpdateStatus measures the cost of a status transition.
func BenchmarkTicketRepo_UpdateStatus(b *testing.B) {
	skipBenchIfNoDB(b)
	ctx := context.Background()
	tenantID, userID := benchTenantUser(b, ctx)
	defer func() {
		_, _ = testPool.Exec(ctx, `DELETE FROM tenants WHERE id = $1`, tenantID)
	}()

	repo := NewTicketRepo(testPool)
	// Seed one ticket to update repeatedly.
	tk := &domainticket.Ticket{
		TenantID:    tenantID,
		Title:       "Bench update ticket",
		Description: "Will be updated",
		Priority:    domainticket.PriorityMedium,
		Category:    "general",
		CreatedBy:   userID,
	}
	if err := repo.Create(ctx, tk); err != nil {
		b.Fatalf("seed Create: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Alternate between open and in_progress to keep transitions valid.
		target := domainticket.StatusInProgress
		if i%2 == 0 {
			target = domainticket.StatusOpen
		}
		if _, err := repo.UpdateStatus(ctx, tk.ID, tenantID, tk.Version, target, userID); err != nil {
			// On concurrent-modification or invalid transition, re-seed version.
			got, gerr := repo.GetByID(ctx, tk.ID, tenantID)
			if gerr != nil {
				b.Fatalf("GetByID: %v", gerr)
			}
			tk.Version = got.Version
			i--
			continue
		}
		// Refresh version for next iteration.
		got, _ := repo.GetByID(ctx, tk.ID, tenantID)
		if got != nil {
			tk.Version = got.Version
		}
	}
	b.StopTimer()
}

// BenchmarkUserRepo_List measures the cost of listing users for a tenant.
func BenchmarkUserRepo_List(b *testing.B) {
	skipBenchIfNoDB(b)
	ctx := context.Background()
	tenantID, _ := benchTenantUser(b, ctx)
	defer func() {
		_, _ = testPool.Exec(ctx, `DELETE FROM tenants WHERE id = $1`, tenantID)
	}()

	repo := NewUserRepo(testPool)
	// Seed 50 users.
	for i := 0; i < 50; i++ {
		u := &domainuser.User{
			TenantID: tenantID,
			Name:     fmt.Sprintf("Bench user %d", i),
			Email:    fmt.Sprintf("bench-%d-%s@test.local", i, tenantID.String()[:8]),
			Role:     domainuser.RoleEndUser,
		}
		if err := repo.Create(ctx, u, ""); err != nil {
			b.Fatalf("seed Create: %v", err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := repo.List(ctx, tenantID); err != nil {
			b.Fatalf("List: %v", err)
		}
	}
	b.StopTimer()
}
