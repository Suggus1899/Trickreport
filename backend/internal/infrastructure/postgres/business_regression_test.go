package postgres

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	appAnalytics "github.com/trickreport/backend/internal/application/analytics"
	appAuth "github.com/trickreport/backend/internal/application/auth"
	appAutomation "github.com/trickreport/backend/internal/application/automation"
	appTicket "github.com/trickreport/backend/internal/application/ticket"
	domainAutomation "github.com/trickreport/backend/internal/domain/automation"
	domainSLA "github.com/trickreport/backend/internal/domain/sla"
	domain "github.com/trickreport/backend/internal/domain/ticket"
	"github.com/trickreport/backend/internal/interfaces/http/handler"
	"github.com/trickreport/backend/internal/interfaces/http/middleware"
)

func businessTicket(t *testing.T, ctx context.Context, tenant, creator uuid.UUID) *domain.Ticket {
	t.Helper()
	tk := &domain.Ticket{TenantID: tenant, CreatedBy: creator, Title: "business regression", Description: "regression test ticket", Priority: domain.PriorityHigh, Category: "general"}
	if err := NewTicketRepo(testPool).Create(ctx, tk); err != nil {
		t.Fatal(err)
	}
	return tk
}

func TestIntegration_TicketSLATrackingPreservesContract(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenant, actor := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenant)
	slaRepo := NewSLARepo(testPool)
	policy := &domainSLA.Policy{TenantID: tenant, Priority: "high", ResponseTimeMinutes: 10, ResolutionTimeMinutes: 60, EscalationMinutes: 0}
	if err := slaRepo.Upsert(ctx, policy); err != nil {
		t.Fatal(err)
	}
	repo := NewTicketRepo(testPool)
	comments := NewCommentRepo(testPool)
	svc := appTicket.NewService(repo, comments, NewHistoryRepo(testPool), nil, nil)
	svc.SetTxManager(NewTxManager(testPool))
	svc.SetSLAPolicyFetcher(slaRepo)
	tk, err := svc.Create(ctx, appTicket.CreateInput{TenantID: tenant, CreatedBy: actor, Title: "SLA ticket", Description: "Full SLA contract description", Priority: "high"})
	if err != nil {
		t.Fatal(err)
	}
	if tk.SLADeadline == nil || tk.ResponseDeadline == nil || tk.EscalationDeadline != nil || tk.FirstRespondedAt != nil || tk.ResolvedAt != nil {
		t.Fatalf("new SLA fields=%+v", tk)
	}
	resolutionDeadline, responseDeadline := *tk.SLADeadline, *tk.ResponseDeadline
	policy.ResponseTimeMinutes = 500
	policy.ResolutionTimeMinutes = 1000
	policy.EscalationMinutes = 1
	if err := slaRepo.Upsert(ctx, policy); err != nil {
		t.Fatal(err)
	}
	if err := repo.UpdatePriority(ctx, tk.ID, tenant, actor, "critical"); err != nil {
		t.Fatal(err)
	}
	if err := comments.Create(ctx, &domain.Comment{TicketID: tk.ID, UserID: actor, Content: "internal", IsInternal: true}); err != nil {
		t.Fatal(err)
	}
	got, err := repo.GetByID(ctx, tk.ID, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if got.FirstRespondedAt != nil {
		t.Fatal("internal note counted as customer response")
	}
	if err := comments.Create(ctx, &domain.Comment{TicketID: tk.ID, UserID: actor, Content: "customer response"}); err != nil {
		t.Fatal(err)
	}
	got, err = repo.GetByID(ctx, tk.ID, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if got.FirstRespondedAt == nil {
		t.Fatal("staff public response not tracked")
	}
	first := *got.FirstRespondedAt
	if err := comments.Create(ctx, &domain.Comment{TicketID: tk.ID, UserID: actor, Content: "second response"}); err != nil {
		t.Fatal(err)
	}
	got, err = repo.UpdateStatus(ctx, tk.ID, tenant, got.Version, domain.StatusResolved, actor)
	if err != nil {
		t.Fatal(err)
	}
	if got.ResolvedAt == nil {
		t.Fatal("resolution timestamp missing")
	}
	got, err = repo.UpdateStatus(ctx, tk.ID, tenant, got.Version, domain.StatusInProgress, actor)
	if err != nil {
		t.Fatal(err)
	}
	if got.ResolvedAt != nil || got.FirstRespondedAt == nil || !got.FirstRespondedAt.Equal(first) || !got.SLADeadline.Equal(resolutionDeadline) || !got.ResponseDeadline.Equal(responseDeadline) || got.EscalationDeadline != nil {
		t.Fatalf("policy/priority/reopen changed contract: %+v", got)
	}
}

func TestIntegration_AutomationRuleRollsBackAllActions(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenant, actor := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenant)
	tk := businessTicket(t, ctx, tenant, actor)
	rules := NewAutomationRepo(testPool)
	rule := &domainAutomation.Rule{TenantID: tenant, Name: "Atomic rule", TriggerType: "ticket_created", Conditions: map[string]any{}, Actions: []any{map[string]any{"type": "set_priority", "priority": "critical"}, map[string]any{"type": "assign_to", "user_id": uuid.New().String()}}, IsActive: true}
	if err := rules.Create(ctx, rule); err != nil {
		t.Fatal(err)
	}
	engine := appAutomation.NewEngine(rules, NewAutomationExecutor(testPool).WithSystemUserID(actor), zerolog.Nop())
	if err := engine.Evaluate(ctx, tenant, appAutomation.Event{Type: "ticket_created", TicketID: tk.ID, TenantID: tenant}); err == nil {
		t.Fatal("invalid assignee must fail entire rule")
	}
	got, err := NewTicketRepo(testPool).GetByID(ctx, tk.ID, tenant)
	if err != nil {
		t.Fatal(err)
	}
	var count int
	if err := testPool.QueryRow(ctx, `SELECT COUNT(*) FROM ticket_history WHERE ticket_id=$1`, tk.ID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if got.Priority != domain.PriorityHigh || got.Version != tk.Version || count != 0 {
		t.Fatalf("partial rule persisted: priority=%s version=%d history=%d", got.Priority, got.Version, count)
	}
}

func TestIntegration_AnalyticsUsesEventTimesAndViewerScope(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenant, actor := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenant)
	other := uuid.New()
	if _, err := testPool.Exec(ctx, `INSERT INTO users(id,tenant_id,name,email,role)VALUES($1,$2,'Other',$3,'end_user')`, other, tenant, other.String()+"@test.local"); err != nil {
		t.Fatal(err)
	}
	day := time.Now().UTC().Truncate(24 * time.Hour)
	if _, err := testPool.Exec(ctx, `INSERT INTO tickets(tenant_id,created_by,title,description,priority,status,created_at,resolved_at,updated_at,sla_deadline,response_deadline,first_responded_at)VALUES
($1,$2,'known','test','high','closed',$4,$5,$6,$6,$5,$5),
($1,$3,'historical','test','high','closed',$4,NULL,$6,NULL,NULL,NULL),
($1,$2,'open','test','high','open',$6,NULL,$6,NULL,NULL,NULL)`, tenant, actor, other, day.Add(-48*time.Hour), day.Add(-24*time.Hour), day); err != nil {
		t.Fatal(err)
	}
	repo := NewAnalyticsRepo(testPool)
	summary, err := repo.GetSummary(ctx, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if summary.TotalTickets != 3 || summary.OpenTickets != 1 || summary.ResolvedTickets != 1 {
		t.Errorf("summary=%+v", summary)
	}
	resolution, err := repo.GetResolutionTime(ctx, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if len(resolution) != 1 || resolution[0].AvgHours != 24 {
		t.Errorf("resolution=%+v want24knownhours", resolution)
	}
	volume, err := repo.GetVolume(ctx, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if len(volume) != 30 {
		t.Errorf("volume days=%d want30includingzero", len(volume))
	}
	h := handler.NewAnalyticsHandler(appAnalytics.NewService(repo))
	endpoint := middleware.Tenant(NewTenantResolver(testPool))(http.HandlerFunc(h.GetSummary))
	req := httptest.NewRequest("GET", "/api/analytics/summary", nil)
	req.Header.Set("X-Tenant-ID", tenant.String())
	req = req.WithContext(middleware.ContextWithClaims(req.Context(), &appAuth.Claims{TenantID: tenant, UserID: actor, Role: "end_user"}))
	rec := httptest.NewRecorder()
	endpoint.ServeHTTP(rec, req)
	var scoped appAnalytics.Summary
	if err := json.Unmarshal(rec.Body.Bytes(), &scoped); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 200 || scoped.TotalTickets != 2 {
		t.Errorf("viewer scope status=%d summary=%+v", rec.Code, scoped)
	}
}

func TestCreateTicket_IdempotencySurvivesResponseLoss(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenant, actor := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenant)
	svc := appTicket.NewService(NewTicketRepo(testPool), NewCommentRepo(testPool), NewHistoryRepo(testPool), nil, nil)
	svc.SetTxManager(NewTxManager(testPool))
	h := handler.NewTicketHandler(svc)
	endpoint := middleware.Tenant(NewTenantResolver(testPool))(http.HandlerFunc(h.Create))
	key := uuid.New().String()
	create := func(user uuid.UUID, body string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/api/tickets", strings.NewReader(body))
		req.Header.Set("Idempotency-Key", key)
		req.Header.Set("X-Tenant-ID", tenant.String())
		req = req.WithContext(middleware.ContextWithClaims(req.Context(), &appAuth.Claims{TenantID: tenant, UserID: user, Role: "admin"}))
		rec := httptest.NewRecorder()
		endpoint.ServeHTTP(rec, req)
		return rec
	}
	body := `{"title":"Idempotent issue","description":"A complete description","priority":"high","category":"general"}`
	first := create(actor, body)
	if first.Code != 201 {
		t.Fatalf("first: %d %s", first.Code, first.Body)
	}
	var result map[string]any
	if err := json.Unmarshal(first.Body.Bytes(), &result); err != nil {
		t.Fatal(err)
	}
	var wg sync.WaitGroup
	responses := make(chan *httptest.ResponseRecorder, 8)
	for range 8 {
		wg.Go(func() { responses <- create(actor, body) })
	}
	wg.Wait()
	close(responses)
	for rec := range responses {
		var replay map[string]any
		_ = json.Unmarshal(rec.Body.Bytes(), &replay)
		if rec.Code != 201 || replay["id"] != result["id"] {
			t.Errorf("replay: %d %s", rec.Code, rec.Body)
		}
	}
	var count int
	if err := testPool.QueryRow(ctx, `SELECT COUNT(*) FROM tickets WHERE tenant_id=$1`, tenant).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Errorf("response loss/retries created %d tickets want1", count)
	}
	if rec := create(actor, strings.Replace(body, "Idempotent issue", "Different payload", 1)); rec.Code != 409 {
		t.Errorf("changed payload status=%d want409", rec.Code)
	}
	other := uuid.New()
	if _, err := testPool.Exec(ctx, `INSERT INTO users(id,tenant_id,name,email,role)VALUES($1,$2,'Other',$3,'admin')`, other, tenant, other.String()+"@test.local"); err != nil {
		t.Fatal(err)
	}
	rec := create(other, body)
	var separate map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &separate)
	if rec.Code != 201 || separate["id"] == result["id"] {
		t.Errorf("other identity retrieved original: %d %s", rec.Code, rec.Body)
	}
	key = uuid.New().String()
	initialResponses := make(chan *httptest.ResponseRecorder, 8)
	start := make(chan struct{})
	for range 8 {
		wg.Go(func() { <-start; initialResponses <- create(actor, body) })
	}
	close(start)
	wg.Wait()
	close(initialResponses)
	var sharedID any
	for rec := range initialResponses {
		var created map[string]any
		_ = json.Unmarshal(rec.Body.Bytes(), &created)
		if sharedID == nil {
			sharedID = created["id"]
		}
		if rec.Code != 201 || created["id"] != sharedID {
			t.Errorf("concurrent initial create: %d %s", rec.Code, rec.Body)
		}
	}
	if err := testPool.QueryRow(ctx, `SELECT COUNT(*) FROM tickets WHERE tenant_id=$1`, tenant).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("concurrent initial creates left%d tickets want3", count)
	}
}

func TestIntegration_AttachmentUploadListDownloadScope(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenant, actor := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenant)
	tk := businessTicket(t, ctx, tenant, actor)
	h := handler.NewAttachmentHandler(appTicket.NewAttachmentService(NewAttachmentRepo(testPool), NewTicketRepo(testPool)))
	router := chi.NewRouter()
	router.Use(middleware.Tenant(NewTenantResolver(testPool)))
	router.Post("/tickets/{id}/attachments", h.Upload)
	router.Get("/tickets/{id}/attachments", h.List)
	router.Get("/tickets/{id}/attachments/{aid}", h.Download)
	request := func(method, url string, data *bytes.Buffer, contentType string, viewerTenant, viewer uuid.UUID) *httptest.ResponseRecorder {
		if data == nil {
			data = &bytes.Buffer{}
		}
		req := httptest.NewRequest(method, url, data)
		req.Header.Set("X-Tenant-ID", viewerTenant.String())
		req.Header.Set("Content-Type", contentType)
		req = req.WithContext(middleware.ContextWithClaims(req.Context(), &appAuth.Claims{TenantID: viewerTenant, UserID: viewer, Role: "admin"}))
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)
		return rec
	}
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("file", "safe.txt")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = part.Write([]byte("hello attachment"))
	_ = writer.Close()
	// CreateFormFile uses octet-stream, so declare the actual content type for this fixture.
	payload := bytes.Replace(body.Bytes(), []byte("application/octet-stream"), []byte("text/plain"), 1)
	url := "/tickets/" + tk.ID.String() + "/attachments"
	rec := request("POST", url, bytes.NewBuffer(payload), writer.FormDataContentType(), tenant, actor)
	if rec.Code != 201 {
		t.Fatalf("upload%d %s", rec.Code, rec.Body)
	}
	var uploaded map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &uploaded); err != nil {
		t.Fatal(err)
	}
	attachmentID := uploaded["id"].(string)
	if rec := request("GET", url, nil, "", tenant, actor); rec.Code != 200 || !strings.Contains(rec.Body.String(), attachmentID) {
		t.Fatalf("list%d %s", rec.Code, rec.Body)
	}
	if rec := request("GET", url+"/"+attachmentID, nil, "", tenant, actor); rec.Code != 200 || rec.Body.String() != "hello attachment" || rec.Header().Get("X-Content-Type-Options") != "nosniff" {
		t.Fatalf("download%d %s", rec.Code, rec.Body)
	}
	otherTicket := businessTicket(t, ctx, tenant, actor)
	if rec := request("GET", "/tickets/"+otherTicket.ID.String()+"/attachments/"+attachmentID, nil, "", tenant, actor); rec.Code != 404 {
		t.Fatalf("other ticket status%d", rec.Code)
	}
	otherTenant, otherActor := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, otherTenant)
	if rec := request("GET", url+"/"+attachmentID, nil, "", otherTenant, otherActor); rec.Code != 404 {
		t.Fatalf("other tenant status%d", rec.Code)
	}
}

func TestIntegration_AssignmentNullAndEligibleAgent(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenant, actor := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenant)
	tk := businessTicket(t, ctx, tenant, actor)
	repo := NewTicketRepo(testPool)
	if err := repo.Assign(ctx, tk.ID, tenant, actor, actor); err != nil {
		t.Fatal(err)
	}
	if err := repo.Assign(ctx, tk.ID, tenant, uuid.Nil, actor); err != nil {
		t.Errorf("unassign must persist NULL: %v", err)
	}
	var assigned *uuid.UUID
	if err := testPool.QueryRow(ctx, `SELECT assigned_to FROM tickets WHERE id=$1`, tk.ID).Scan(&assigned); err != nil {
		t.Fatal(err)
	}
	if assigned != nil {
		t.Errorf("unassigned=%v want NULL", assigned)
	}
	otherTenant, otherUser := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, otherTenant)
	for _, user := range []struct {
		id     uuid.UUID
		role   string
		active bool
	}{{uuid.New(), "end_user", true}, {uuid.New(), "agent", false}} {
		if _, err := testPool.Exec(ctx, `INSERT INTO users(id,tenant_id,name,email,role,active) VALUES($1,$2,'test',$3,$4,$5)`, user.id, tenant, user.id.String()+"@test.local", user.role, user.active); err != nil {
			t.Fatal(err)
		}
		if err := repo.Assign(ctx, tk.ID, tenant, user.id, actor); !errors.Is(err, domain.ErrValidation) {
			t.Errorf("assignee role=%s active=%t err=%v", user.role, user.active, err)
		}
	}
	if err := repo.Assign(ctx, tk.ID, tenant, otherUser, actor); !errors.Is(err, domain.ErrValidation) {
		t.Errorf("cross-tenant assignee err=%v", err)
	}
}

func TestIntegration_AutomationPreservesDomainAndVersion(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenant, actor := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenant)
	tk := businessTicket(t, ctx, tenant, actor)
	executor := NewAutomationExecutor(testPool).WithSystemUserID(actor)
	if err := executor.SetStatus(ctx, tenant, tk.ID, "closed"); err != nil {
		t.Fatal(err)
	}
	if err := executor.SetStatus(ctx, tenant, tk.ID, "open"); !errors.Is(err, domain.ErrInvalidTransition) {
		t.Errorf("closed ticket must not reopen: %v", err)
	}
	if err := executor.SetPriority(ctx, tenant, tk.ID, "impossible"); err == nil {
		t.Error("invalid priority accepted")
	}
	got, err := NewTicketRepo(testPool).GetByID(ctx, tk.ID, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if got.Version != tk.Version+1 {
		t.Errorf("version=%d want%d", got.Version, tk.Version+1)
	}
	other, otherUser := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, other)
	if err := executor.AssignTo(ctx, tenant, tk.ID, otherUser); !errors.Is(err, domain.ErrValidation) {
		t.Errorf("cross tenant automation assignment err=%v", err)
	}
	if err := executor.AddComment(ctx, other, tk.ID, "cross tenant", true); err == nil {
		t.Error("cross tenant automation comment accepted")
	}
	if err := executor.AddTag(ctx, tenant, tk.ID, "unsupported"); err == nil {
		t.Error("unsupported tag reports success")
	}
}

func TestIntegration_StatusConflictAndHistory(t *testing.T) {
	skipIfNoDB(t)
	ctx := context.Background()
	tenant, actor := seedTenantUser(t, ctx)
	defer cleanupTenant(ctx, tenant)
	tk := businessTicket(t, ctx, tenant, actor)
	repo := NewTicketRepo(testPool)
	if _, err := repo.UpdateStatus(ctx, tk.ID, tenant, tk.Version, domain.StatusInProgress, actor); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.UpdateStatus(ctx, tk.ID, tenant, tk.Version, domain.StatusResolved, actor); !errors.Is(err, domain.ErrConcurrentModification) {
		t.Fatalf("stale version err=%v", err)
	}
	var old string
	var count int
	if err := testPool.QueryRow(ctx, `SELECT COUNT(*),MIN(old_value) FROM ticket_history WHERE ticket_id=$1 AND field='status'`, tk.ID).Scan(&count, &old); err != nil {
		t.Fatal(err)
	}
	if count != 1 || old != "open" {
		t.Fatalf("history count=%d old=%q", count, old)
	}
	_, err := repo.List(ctx, tenant, appTicket.Filter{Limit: 10}, "end_user", actor)
	if err != nil {
		t.Fatal(err)
	}
}
