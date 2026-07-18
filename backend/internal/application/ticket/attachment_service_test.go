package ticket

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
)

// --- Mocks ---

type mockAttachmentRepo struct {
	attachment *domainTicket.Attachment
	data       []byte
	attachments []domainTicket.Attachment
	err        error
	createErr  error
	gotInput   AttachmentUploadInput
	gotID      uuid.UUID
	gotTicket  uuid.UUID
	createCalls int
}

func (m *mockAttachmentRepo) Create(ctx context.Context, input AttachmentUploadInput) (*domainTicket.Attachment, error) {
	m.createCalls++
	m.gotInput = input
	if m.createErr != nil {
		return nil, m.createErr
	}
	return m.attachment, nil
}

func (m *mockAttachmentRepo) List(ctx context.Context, ticketID uuid.UUID) ([]domainTicket.Attachment, error) {
	m.gotTicket = ticketID
	return m.attachments, m.err
}

func (m *mockAttachmentRepo) GetByID(ctx context.Context, id, ticketID uuid.UUID) (*domainTicket.Attachment, []byte, error) {
	m.gotID = id
	m.gotTicket = ticketID
	return m.attachment, m.data, m.err
}

// --- Helpers ---

func newAttachmentSvc(attachments AttachmentRepository, tickets Repository) *AttachmentService {
	return NewAttachmentService(attachments, tickets)
}

// --- Upload tests ---

func TestAttachmentService_Upload_Success(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &mockRepo{ticket: tk}
	attRepo := &mockAttachmentRepo{attachment: &domainTicket.Attachment{ID: uuid.New(), Filename: "file.png"}}
	svc := newAttachmentSvc(attRepo, ticketRepo)

	att, err := svc.Upload(context.Background(), UploadInput{
		TenantID: tk.TenantID, TicketID: tk.ID, UserID: creator, Role: "end_user",
		Filename: "file.png", ContentType: "image/png", FileData: []byte("data"),
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if att.Filename != "file.png" {
		t.Errorf("filename = %q", att.Filename)
	}
	if attRepo.createCalls != 1 {
		t.Errorf("create calls = %d", attRepo.createCalls)
	}
	if attRepo.gotInput.FileSize != 4 {
		t.Errorf("file size = %d, want 4", attRepo.gotInput.FileSize)
	}
	if attRepo.gotInput.ContentType != "image/png" {
		t.Errorf("content type = %q", attRepo.gotInput.ContentType)
	}
}

func TestAttachmentService_Upload_Forbidden(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &mockRepo{ticket: tk}
	attRepo := &mockAttachmentRepo{}
	svc := newAttachmentSvc(attRepo, ticketRepo)

	_, err := svc.Upload(context.Background(), UploadInput{
		TenantID: tk.TenantID, TicketID: tk.ID, UserID: uuid.New(), Role: "end_user",
		Filename: "file.png", ContentType: "image/png", FileData: []byte("x"),
	})
	if !errors.Is(err, domainTicket.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
	if attRepo.createCalls != 0 {
		t.Errorf("expected 0 create calls, got %d", attRepo.createCalls)
	}
}

func TestAttachmentService_Upload_TicketNotFound(t *testing.T) {
	ticketRepo := &mockRepo{err: domainTicket.ErrNotFound}
	attRepo := &mockAttachmentRepo{}
	svc := newAttachmentSvc(attRepo, ticketRepo)

	_, err := svc.Upload(context.Background(), UploadInput{
		TicketID: uuid.New(), TenantID: uuid.New(), UserID: uuid.New(), Role: "agent",
		Filename: "f", ContentType: "image/png", FileData: []byte("x"),
	})
	if !errors.Is(err, domainTicket.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAttachmentService_Upload_CreateError(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &mockRepo{ticket: tk}
	attRepo := &mockAttachmentRepo{createErr: errors.New("storage fail")}
	svc := newAttachmentSvc(attRepo, ticketRepo)

	_, err := svc.Upload(context.Background(), UploadInput{
		TenantID: tk.TenantID, TicketID: tk.ID, UserID: creator, Role: "end_user",
		Filename: "f", ContentType: "image/png", FileData: []byte("x"),
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- List tests ---

func TestAttachmentService_List_Success(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &mockRepo{ticket: tk}
	attRepo := &mockAttachmentRepo{attachments: []domainTicket.Attachment{{ID: uuid.New(), Filename: "a.png"}}}
	svc := newAttachmentSvc(attRepo, ticketRepo)

	got, err := svc.List(context.Background(), ListInput{
		TenantID: tk.TenantID, TicketID: tk.ID, UserID: creator, Role: "end_user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 {
		t.Errorf("got %d attachments", len(got))
	}
	if attRepo.gotTicket != tk.ID {
		t.Errorf("repo got ticket %v", attRepo.gotTicket)
	}
}

func TestAttachmentService_List_Forbidden(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &mockRepo{ticket: tk}
	attRepo := &mockAttachmentRepo{}
	svc := newAttachmentSvc(attRepo, ticketRepo)

	_, err := svc.List(context.Background(), ListInput{
		TenantID: tk.TenantID, TicketID: tk.ID, UserID: uuid.New(), Role: "end_user",
	})
	if !errors.Is(err, domainTicket.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestAttachmentService_List_TicketNotFound(t *testing.T) {
	ticketRepo := &mockRepo{err: domainTicket.ErrNotFound}
	attRepo := &mockAttachmentRepo{}
	svc := newAttachmentSvc(attRepo, ticketRepo)

	_, err := svc.List(context.Background(), ListInput{
		TicketID: uuid.New(), TenantID: uuid.New(), UserID: uuid.New(), Role: "agent",
	})
	if !errors.Is(err, domainTicket.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAttachmentService_List_RepoError(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &mockRepo{ticket: tk}
	attRepo := &mockAttachmentRepo{err: errors.New("db")}
	svc := newAttachmentSvc(attRepo, ticketRepo)

	_, err := svc.List(context.Background(), ListInput{
		TenantID: tk.TenantID, TicketID: tk.ID, UserID: creator, Role: "end_user",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

// --- Download tests ---

func TestAttachmentService_Download_Success(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	attID := uuid.New()
	ticketRepo := &mockRepo{ticket: tk}
	attRepo := &mockAttachmentRepo{
		attachment: &domainTicket.Attachment{ID: attID, Filename: "file.pdf"},
		data:       []byte("pdf-content"),
	}
	svc := newAttachmentSvc(attRepo, ticketRepo)

	att, data, err := svc.Download(context.Background(), DownloadInput{
		TenantID: tk.TenantID, TicketID: tk.ID, AttachmentID: attID, UserID: creator, Role: "end_user",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if att.ID != attID {
		t.Errorf("attachment id = %v", att.ID)
	}
	if string(data) != "pdf-content" {
		t.Errorf("data = %q", string(data))
	}
	if attRepo.gotID != attID || attRepo.gotTicket != tk.ID {
		t.Errorf("repo got id=%v ticket=%v", attRepo.gotID, attRepo.gotTicket)
	}
}

func TestAttachmentService_Download_Forbidden(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &mockRepo{ticket: tk}
	attRepo := &mockAttachmentRepo{}
	svc := newAttachmentSvc(attRepo, ticketRepo)

	_, _, err := svc.Download(context.Background(), DownloadInput{
		TenantID: tk.TenantID, TicketID: tk.ID, AttachmentID: uuid.New(), UserID: uuid.New(), Role: "end_user",
	})
	if !errors.Is(err, domainTicket.ErrForbidden) {
		t.Errorf("expected ErrForbidden, got %v", err)
	}
}

func TestAttachmentService_Download_TicketNotFound(t *testing.T) {
	ticketRepo := &mockRepo{err: domainTicket.ErrNotFound}
	attRepo := &mockAttachmentRepo{}
	svc := newAttachmentSvc(attRepo, ticketRepo)

	_, _, err := svc.Download(context.Background(), DownloadInput{
		TicketID: uuid.New(), TenantID: uuid.New(), AttachmentID: uuid.New(), UserID: uuid.New(), Role: "agent",
	})
	if !errors.Is(err, domainTicket.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestAttachmentService_Download_RepoError(t *testing.T) {
	creator := uuid.New()
	tk := ownedTicket(creator)
	ticketRepo := &mockRepo{ticket: tk}
	attRepo := &mockAttachmentRepo{err: errors.New("not found")}
	svc := newAttachmentSvc(attRepo, ticketRepo)

	_, _, err := svc.Download(context.Background(), DownloadInput{
		TenantID: tk.TenantID, TicketID: tk.ID, AttachmentID: uuid.New(), UserID: creator, Role: "end_user",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestNewAttachmentService(t *testing.T) {
	svc := NewAttachmentService(&mockAttachmentRepo{}, &mockRepo{})
	if svc == nil {
		t.Fatal("expected non-nil service")
	}
}
