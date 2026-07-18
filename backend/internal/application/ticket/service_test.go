package ticket

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	appAuto "github.com/trickreport/backend/internal/application/automation"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
)

// --- Mocks ---

type mockRepo struct {
	tickets     []domainTicket.Ticket
	ticket      *domainTicket.Ticket
	err         error
	updateErr   error
	assignErr   error
	creatorID   uuid.UUID
	creatorErr  error
	gotFilter   Filter
	gotRole     string
	gotUserID   uuid.UUID
	gotStatus   domainTicket.Status
	gotAssigned uuid.UUID
	createCalls int
	updateCalls int
	assignCalls int
}

func (m *mockRepo) List(ctx context.Context, tenantID uuid.UUID, filter Filter, role string, userID uuid.UUID) ([]domainTicket.Ticket, error) {
	m.gotFilter = filter
	m.gotRole = role
	m.gotUserID = userID
	return m.tickets, m.err
}

func (m *mockRepo) GetByID(ctx context.Context, id, tenantID uuid.UUID) (*domainTicket.Ticket, error) {
	return m.ticket, m.err
}

func (m *mockRepo) Create(ctx context.Context, t *domainTicket.Ticket) error {
	m.createCalls++
	if m.err != nil && m.ticket == nil {
		return m.err
	}
	return m.err
}

func (m *mockRepo) UpdateStatus(ctx context.Context, id, tenantID uuid.UUID, status domainTicket.Status, userID uuid.UUID) (*domainTicket.Ticket, error) {
	m.updateCalls++
	m.gotStatus = status
	m.gotUserID = userID
	if m.updateErr != nil {
		return nil, m.updateErr
	}
	return m.ticket, nil
}

func (m *mockRepo) Assign(ctx context.Context, id, tenantID, assignedTo, userID uuid.UUID) error {
	m.assignCalls++
	m.gotAssigned = assignedTo
	m.gotUserID = userID
	return m.assignErr
}

func (m *mockRepo) GetCreator(ctx context.Context, id, tenantID uuid.UUID) (uuid.UUID, error) {
	return m.creatorID, m.creatorErr
}

type mockCommentRepo struct {
	comments     []domainTicket.Comment
	comment      *domainTicket.Comment
	err          error
	createErr    error
	gotRole      string
	createCalls  int
	gotComment   *domainTicket.Comment
}

func (m *mockCommentRepo) List(ctx context.Context, ticketID uuid.UUID, role string) ([]domainTicket.Comment, error) {
	m.gotRole = role
	return m.comments, m.err
}

func (m *mockCommentRepo) Create(ctx context.Context, c *domainTicket.Comment) error {
	m.createCalls++
	m.gotComment = c
	return m.createErr
}

type mockHistoryRepo struct {
	entries []domainTicket.HistoryEntry
	err     error
}

func (m *mockHistoryRepo) List(ctx context.Context, ticketID uuid.UUID) ([]domainTicket.HistoryEntry, error) {
	return m.entries, m.err
}

type mockHub struct {
	broadcasts   int
	gotTenant    uuid.UUID
	gotEventType string
	gotData      any
}

func (m *mockHub) BroadcastEvent(tenantID uuid.UUID, eventType string, data any) {
	m.broadcasts++
	m.gotTenant = tenantID
	m.gotEventType = eventType
	m.gotData = data
}

type mockEmail struct {
	sent    int
	gotTo   string
	gotSubj string
	gotBody string
	err     error
}

func (m *mockEmail) Send(to, subject, body string) error {
	m.sent++
	m.gotTo = to
	m.gotSubj = subject
	m.gotBody = body
	return m.err
}

// --- Helpers ---

func newSvc(repo Repository, comments CommentRepository, history HistoryRepository, hub EventBroadcaster, email EmailNotifier) *UserService {
	return NewService(repo, comments, history, hub, email)
}

func ownedTicket(creator uuid.UUID) *domainTicket.Ticket {
	return &domainTicket.Ticket{
		ID:        uuid.New(),
		TenantID:  uuid.New(),
		Title:     "Broken login",
		Status:    domainTicket.StatusOpen,
		Priority:  domainTicket.PriorityMedium,
		Category:  "general",
		CreatedBy: creator,
	}
}

// --- List tests ---

func TestService_List_DelegatesToRepo(t *testing.T) {
	tenant := uuid.New()
	user := uuid.New()
	want := []domainTicket.Ticket{{ID: uuid.New(), Title: "t1"}}
	repo := &mockRepo{tickets: want}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	got, err := svc.List(context.Background(), tenant, Filter{Status: "open"}, "agent", user)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Title != "t1" {
		t.Errorf("got %v", got)
	}
	if repo.gotFilter.Status != "open" || repo.gotRole != "agent" || repo.gotUserID != user {
		t.Errorf("repo got filter=%+v role=%q user=%v", repo.gotFilter, repo.gotRole, repo.gotUserID)
	}
}

func TestService_List_RepoError(t *testing.T) {
	repo := &mockRepo{err: errors.New("db down")}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	if _, err := svc.List(context.Background(), uuid.New(), Filter{}, "agent", uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}

// --- Get tests ---

func TestService_Get_Success(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &mockRepo{ticket: tk}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	got, err := svc.Get(context.Background(), tk.ID, tk.TenantID, "agent", uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != tk.ID {
		t.Errorf("got id %v", got.ID)
	}
}

func TestService_Get_ForbiddenForEndUser(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &mockRepo{ticket: tk}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	_, err := svc.Get(context.Background(), tk.ID, tk.TenantID, "end_user", uuid.New())
	if !errors.Is(err, domainTicket.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestService_Get_NotFound(t *testing.T) {
	repo := &mockRepo{err: domainTicket.ErrNotFound}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	_, err := svc.Get(context.Background(), uuid.New(), uuid.New(), "agent", uuid.New())
	if !errors.Is(err, domainTicket.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// --- Create tests ---

func TestService_Create_Success(t *testing.T) {
	tenant := uuid.New()
	creator := uuid.New()
	repo := &mockRepo{}
	hub := &mockHub{}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, hub, nil)

	tk, err := svc.Create(context.Background(), CreateInput{
		TenantID:    tenant,
		Title:       "New issue",
		Description: "Something broke",
		Priority:    "high",
		Category:    "bug",
		CreatedBy:   creator,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tk.Title != "New issue" {
		t.Errorf("title = %q", tk.Title)
	}
	if tk.Status != domainTicket.StatusOpen {
		t.Errorf("status = %q", tk.Status)
	}
	if tk.Priority != domainTicket.PriorityHigh {
		t.Errorf("priority = %q", tk.Priority)
	}
	if tk.Category != "bug" {
		t.Errorf("category = %q", tk.Category)
	}
	if tk.CreatedBy != creator {
		t.Errorf("created_by = %v", tk.CreatedBy)
	}
	if repo.createCalls != 1 {
		t.Errorf("create calls = %d", repo.createCalls)
	}
	if hub.broadcasts != 1 || hub.gotEventType != "TICKET_CREATED" || hub.gotTenant != tenant {
		t.Errorf("hub broadcast = %+v", hub)
	}
}

func TestService_Create_DefaultsCategoryAndPriority(t *testing.T) {
	repo := &mockRepo{}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	tk, err := svc.Create(context.Background(), CreateInput{
		TenantID:    uuid.New(),
		Title:       "t",
		Description: "d",
		CreatedBy:   uuid.New(),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tk.Category != "general" {
		t.Errorf("category = %q, want general", tk.Category)
	}
	if tk.Priority != domainTicket.PriorityLow {
		t.Errorf("priority = %q, want low", tk.Priority)
	}
}

func TestService_Create_EmptyTitle(t *testing.T) {
	svc := newSvc(&mockRepo{}, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)
	_, err := svc.Create(context.Background(), CreateInput{Title: "", Description: "d"})
	if !errors.Is(err, domainTicket.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Create_EmptyDescription(t *testing.T) {
	svc := newSvc(&mockRepo{}, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)
	_, err := svc.Create(context.Background(), CreateInput{Title: "t", Description: ""})
	if !errors.Is(err, domainTicket.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Create_InvalidPriority(t *testing.T) {
	svc := newSvc(&mockRepo{}, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)
	_, err := svc.Create(context.Background(), CreateInput{Title: "t", Description: "d", Priority: "urgent"})
	if !errors.Is(err, domainTicket.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_Create_RepoError(t *testing.T) {
	repo := &mockRepo{err: errors.New("insert failed")}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)
	_, err := svc.Create(context.Background(), CreateInput{Title: "t", Description: "d"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestService_Create_NoHubNoBroadcast(t *testing.T) {
	repo := &mockRepo{}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)
	if _, err := svc.Create(context.Background(), CreateInput{Title: "t", Description: "d"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// --- ChangeStatus tests ---

func TestService_ChangeStatus_Success(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &mockRepo{ticket: tk}
	hub := &mockHub{}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, hub, nil)

	if err := svc.ChangeStatus(context.Background(), tk.ID, tk.TenantID, uuid.New(), "agent", "in_progress"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.gotStatus != domainTicket.StatusInProgress {
		t.Errorf("repo got status %q", repo.gotStatus)
	}
	if repo.updateCalls != 1 {
		t.Errorf("update calls = %d", repo.updateCalls)
	}
	if hub.broadcasts != 1 || hub.gotEventType != "TICKET_UPDATED" {
		t.Errorf("hub = %+v", hub)
	}
}

func TestService_ChangeStatus_InvalidStatus(t *testing.T) {
	tk := ownedTicket(uuid.New())
	repo := &mockRepo{ticket: tk}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	err := svc.ChangeStatus(context.Background(), tk.ID, tk.TenantID, uuid.New(), "agent", "bogus")
	if !errors.Is(err, domainTicket.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_ChangeStatus_NotFound(t *testing.T) {
	repo := &mockRepo{err: domainTicket.ErrNotFound}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	err := svc.ChangeStatus(context.Background(), uuid.New(), uuid.New(), uuid.New(), "agent", "in_progress")
	if !errors.Is(err, domainTicket.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestService_ChangeStatus_ForbiddenEndUserReopen(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &mockRepo{ticket: tk}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	err := svc.ChangeStatus(context.Background(), tk.ID, tk.TenantID, creator, "end_user", "open")
	if !errors.Is(err, domainTicket.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
	if repo.updateCalls != 0 {
		t.Errorf("expected no update calls, got %d", repo.updateCalls)
	}
}

func TestService_ChangeStatus_InvalidTransition(t *testing.T) {
	tk := ownedTicket(uuid.New())
	tk.Status = domainTicket.StatusClosed
	repo := &mockRepo{ticket: tk}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	err := svc.ChangeStatus(context.Background(), tk.ID, tk.TenantID, uuid.New(), "agent", "open")
	if !errors.Is(err, domainTicket.ErrInvalidTransition) {
		t.Errorf("expected ErrInvalidTransition, got %v", err)
	}
}

func TestService_ChangeStatus_UpdateError(t *testing.T) {
	tk := ownedTicket(uuid.New())
	repo := &mockRepo{ticket: tk, updateErr: errors.New("update failed")}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	err := svc.ChangeStatus(context.Background(), tk.ID, tk.TenantID, uuid.New(), "agent", "in_progress")
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Assign tests ---

func TestService_Assign_DelegatesToRepo(t *testing.T) {
	repo := &mockRepo{}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	assignee := uuid.New()
	actor := uuid.New()
	if err := svc.Assign(context.Background(), uuid.New(), uuid.New(), assignee, actor); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if repo.gotAssigned != assignee || repo.gotUserID != actor {
		t.Errorf("repo got assigned=%v user=%v", repo.gotAssigned, repo.gotUserID)
	}
}

func TestService_Assign_RepoError(t *testing.T) {
	repo := &mockRepo{assignErr: errors.New("nope")}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	if err := svc.Assign(context.Background(), uuid.New(), uuid.New(), uuid.New(), uuid.New()); err == nil {
		t.Fatal("expected error")
	}
}

// --- ListComments tests ---

func TestService_ListComments_Success(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &mockRepo{ticket: tk}
	comments := &mockCommentRepo{comments: []domainTicket.Comment{{ID: uuid.New(), Content: "hi"}}}
	svc := newSvc(repo, comments, &mockHistoryRepo{}, nil, nil)

	got, err := svc.ListComments(context.Background(), tk.ID, tk.TenantID, "agent", uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("got %d comments", len(got))
	}
	if comments.gotRole != "agent" {
		t.Errorf("role = %q", comments.gotRole)
	}
}

func TestService_ListComments_Forbidden(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &mockRepo{ticket: tk}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	_, err := svc.ListComments(context.Background(), tk.ID, tk.TenantID, "end_user", uuid.New())
	if !errors.Is(err, domainTicket.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestService_ListComments_NotFound(t *testing.T) {
	repo := &mockRepo{err: domainTicket.ErrNotFound}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	_, err := svc.ListComments(context.Background(), uuid.New(), uuid.New(), "agent", uuid.New())
	if !errors.Is(err, domainTicket.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// --- AddComment tests ---

func TestService_AddComment_Success(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &mockRepo{ticket: tk}
	comments := &mockCommentRepo{}
	hub := &mockHub{}
	svc := newSvc(repo, comments, &mockHistoryRepo{}, hub, nil)

	actor := uuid.New()
	c, err := svc.AddComment(context.Background(), AddCommentInput{
		TicketID: tk.ID, TenantID: tk.TenantID, UserID: actor, Role: "agent",
		Content: "looks bad", IsInternal: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.Content != "looks bad" || !c.IsInternal {
		t.Errorf("comment = %+v", c)
	}
	if comments.createCalls != 1 {
		t.Errorf("create calls = %d", comments.createCalls)
	}
	if hub.broadcasts != 1 || hub.gotEventType != "NEW_COMMENT" {
		t.Errorf("hub = %+v", hub)
	}
}

func TestService_AddComment_EmptyContent(t *testing.T) {
	svc := newSvc(&mockRepo{}, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)
	_, err := svc.AddComment(context.Background(), AddCommentInput{Content: ""})
	if !errors.Is(err, domainTicket.ErrValidation) {
		t.Errorf("expected ErrValidation, got %v", err)
	}
}

func TestService_AddComment_EndUserCannotMakeInternal(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &mockRepo{ticket: tk}
	comments := &mockCommentRepo{}
	svc := newSvc(repo, comments, &mockHistoryRepo{}, nil, nil)

	c, err := svc.AddComment(context.Background(), AddCommentInput{
		TicketID: tk.ID, TenantID: tk.TenantID, UserID: creator, Role: "end_user",
		Content: "thanks", IsInternal: true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if c.IsInternal {
		t.Error("expected IsInternal=false for end_user")
	}
}

func TestService_AddComment_Forbidden(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &mockRepo{ticket: tk}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	_, err := svc.AddComment(context.Background(), AddCommentInput{
		TicketID: tk.ID, TenantID: tk.TenantID, UserID: uuid.New(), Role: "end_user",
		Content: "hi",
	})
	if !errors.Is(err, domainTicket.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestService_AddComment_NotFound(t *testing.T) {
	repo := &mockRepo{err: domainTicket.ErrNotFound}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	_, err := svc.AddComment(context.Background(), AddCommentInput{Content: "hi"})
	if !errors.Is(err, domainTicket.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestService_AddComment_CreateError(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &mockRepo{ticket: tk}
	comments := &mockCommentRepo{createErr: errors.New("insert failed")}
	svc := newSvc(repo, comments, &mockHistoryRepo{}, nil, nil)

	_, err := svc.AddComment(context.Background(), AddCommentInput{
		TicketID: tk.ID, TenantID: tk.TenantID, UserID: creator, Role: "end_user", Content: "hi",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- ListHistory tests ---

func TestService_ListHistory_Success(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &mockRepo{ticket: tk}
	history := &mockHistoryRepo{entries: []domainTicket.HistoryEntry{{ID: uuid.New(), Field: "status"}}}
	svc := newSvc(repo, &mockCommentRepo{}, history, nil, nil)

	got, err := svc.ListHistory(context.Background(), tk.ID, tk.TenantID, "agent", uuid.New())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("got %d entries", len(got))
	}
}

func TestService_ListHistory_Forbidden(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	repo := &mockRepo{ticket: tk}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	_, err := svc.ListHistory(context.Background(), tk.ID, tk.TenantID, "end_user", uuid.New())
	if !errors.Is(err, domainTicket.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestService_ListHistory_NotFound(t *testing.T) {
	repo := &mockRepo{err: domainTicket.ErrNotFound}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)

	_, err := svc.ListHistory(context.Background(), uuid.New(), uuid.New(), "agent", uuid.New())
	if !errors.Is(err, domainTicket.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

// --- SetEngine and evaluateAutomation tests ---

type mockEngine struct {
	evaluated bool
	tenantID  uuid.UUID
	err       error
}

func (m *mockEngine) Evaluate(ctx context.Context, tenantID uuid.UUID, event appAuto.Event) error {
	m.evaluated = true
	m.tenantID = tenantID
	return m.err
}

func TestService_SetEngine(t *testing.T) {
	svc := newSvc(&mockRepo{}, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)
	engine := &mockEngine{}
	svc.SetEngine(engine)
	// Verify engine is set by triggering evaluateAutomation through Create
	// which calls evaluateAutomation internally.
}

func TestService_EvaluateAutomation_WithEngine(t *testing.T) {
	repo := &mockRepo{}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)
	engine := &mockEngine{}
	svc.SetEngine(engine)

	tenantID := uuid.New()
	// Create a ticket which triggers evaluateAutomation
	_, _ = svc.Create(context.Background(), CreateInput{
		TenantID:    tenantID,
		Title:       "Test",
		Description: "Description",
		Priority:    "medium",
		Category:    "general",
		CreatedBy:   uuid.New(),
	})

	// Wait for the goroutine to complete
	time.Sleep(50 * time.Millisecond)
	if !engine.evaluated {
		t.Error("expected engine.Evaluate to be called")
	}
}

func TestService_EvaluateAutomation_NoEngine(t *testing.T) {
	repo := &mockRepo{}
	svc := newSvc(repo, &mockCommentRepo{}, &mockHistoryRepo{}, nil, nil)
	// No engine set — evaluateAutomation should be a no-op
	// This is tested implicitly by Create not panicking
	_, _ = svc.Create(context.Background(), CreateInput{
		TenantID:    uuid.New(),
		Title:       "Test",
		Description: "Description",
		Priority:    "medium",
		Category:    "general",
		CreatedBy:   uuid.New(),
	})
}
