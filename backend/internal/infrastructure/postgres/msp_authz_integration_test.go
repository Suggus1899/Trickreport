package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	appMSP "github.com/trickreport/backend/internal/application/msp"
	domain "github.com/trickreport/backend/internal/domain/msp"
)

// TestIntegration_MSP_ContactIsScopedToOwnClient checks the per-client rule on
// the two MSP operations that did not apply it: listing a project's tasks and
// requesting a remote session. A contact linked to client A must get NotFound
// for client B, the same answer the other MSP endpoints give. Staff keep full
// access, and the contact keeps access to its own client.
func TestIntegration_MSP_ContactIsScopedToOwnClient(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenantID, adminID := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenantID)

	contactID := uuid.New()
	if _, err := testPool.Exec(ctx,
		`INSERT INTO users (id, tenant_id, name, email, role) VALUES ($1, $2, $3, $4, 'end_user')`,
		contactID, tenantID, "Contact", "contact-"+tenantID.String()[:8]+"@test.local",
	); err != nil {
		t.Fatalf("seed contact user: %v", err)
	}

	repo := NewMSPRepo(testPool)
	svc := appMSP.NewService(repo)

	clientA := &domain.Client{TenantID: tenantID, Name: "Client A", Domain: "a-" + tenantID.String()[:8] + ".test.local", Status: "active"}
	clientB := &domain.Client{TenantID: tenantID, Name: "Client B", Domain: "b-" + tenantID.String()[:8] + ".test.local", Status: "active"}
	for _, c := range []*domain.Client{clientA, clientB} {
		if err := repo.CreateClient(ctx, c); err != nil {
			t.Fatalf("CreateClient: %v", err)
		}
	}
	if err := repo.CreateContact(ctx, &domain.ClientContact{
		TenantID: tenantID, ClientID: clientA.ID, UserID: &contactID,
		Name: "Contact", Email: "contact@test.local", Role: "contact",
	}); err != nil {
		t.Fatalf("CreateContact: %v", err)
	}

	tech := &domain.Technician{TenantID: tenantID, UserID: adminID, Skills: []string{}, Active: true}
	if err := repo.CreateTechnician(ctx, tech); err != nil {
		t.Fatalf("CreateTechnician: %v", err)
	}

	projectA := &domain.Project{TenantID: tenantID, ClientID: clientA.ID, Name: "Project A", Status: "planned"}
	projectB := &domain.Project{TenantID: tenantID, ClientID: clientB.ID, Name: "Project B", Status: "planned"}
	for _, p := range []*domain.Project{projectA, projectB} {
		if err := repo.CreateProject(ctx, p); err != nil {
			t.Fatalf("CreateProject: %v", err)
		}
		task := &domain.ProjectTask{TenantID: tenantID, ProjectID: p.ID, Title: "Task " + p.Name, Status: "todo"}
		if err := repo.CreateTask(ctx, task); err != nil {
			t.Fatalf("CreateTask: %v", err)
		}
	}

	t.Run("contact lists tasks of own client", func(t *testing.T) {
		tasks, err := svc.ListTasks(ctx, tenantID, contactID, "end_user", projectA.ID)
		if err != nil {
			t.Fatalf("ListTasks: %v", err)
		}
		if len(tasks) != 1 {
			t.Fatalf("got %d tasks, want 1", len(tasks))
		}
	})

	t.Run("contact cannot list tasks of another client", func(t *testing.T) {
		_, err := svc.ListTasks(ctx, tenantID, contactID, "end_user", projectB.ID)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("ListTasks err = %v, want ErrNotFound", err)
		}
	})

	t.Run("staff lists tasks of any client", func(t *testing.T) {
		tasks, err := svc.ListTasks(ctx, tenantID, adminID, "admin", projectB.ID)
		if err != nil {
			t.Fatalf("ListTasks: %v", err)
		}
		if len(tasks) != 1 {
			t.Fatalf("got %d tasks, want 1", len(tasks))
		}
	})

	t.Run("contact requests a remote session for own client", func(t *testing.T) {
		sess := &domain.RemoteSession{ClientID: clientA.ID, TechnicianID: tech.ID}
		if _, err := svc.RequestRemoteSession(ctx, tenantID, contactID, "end_user", sess); err != nil {
			t.Fatalf("RequestRemoteSession: %v", err)
		}
	})

	t.Run("contact cannot request a remote session for another client", func(t *testing.T) {
		sess := &domain.RemoteSession{ClientID: clientB.ID, TechnicianID: tech.ID}
		_, err := svc.RequestRemoteSession(ctx, tenantID, contactID, "end_user", sess)
		if !errors.Is(err, domain.ErrNotFound) {
			t.Fatalf("RequestRemoteSession err = %v, want ErrNotFound", err)
		}
	})
}
