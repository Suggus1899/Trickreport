package ticket

import (
	"context"
	"io"

	"github.com/google/uuid"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
)

// AttachmentRepository is the port for attachment persistence.
type AttachmentRepository interface {
	Create(ctx context.Context, input AttachmentUploadInput) (*domainTicket.Attachment, error)
	List(ctx context.Context, ticketID uuid.UUID) ([]domainTicket.Attachment, error)
	GetByID(ctx context.Context, id, ticketID uuid.UUID) (*domainTicket.Attachment, []byte, error)
}

// AttachmentUploadInput holds the data for storing an attachment at the repo level.
type AttachmentUploadInput struct {
	TicketID    uuid.UUID
	UserID      uuid.UUID
	Filename    string
	ContentType string
	FileSize    int64
	FileData    []byte
}

// AttachmentService handles file attachment operations.
type AttachmentService struct {
	attachments AttachmentRepository
	tickets     Repository
}

// NewAttachmentService creates a new AttachmentService.
func NewAttachmentService(attachments AttachmentRepository, tickets Repository) *AttachmentService {
	return &AttachmentService{attachments: attachments, tickets: tickets}
}

// UploadInput holds the data for uploading an attachment.
type UploadInput struct {
	TenantID    uuid.UUID
	TicketID    uuid.UUID
	UserID      uuid.UUID
	Role        string
	Filename    string
	ContentType string
	FileData    []byte
}

// Upload stores a file attachment on a ticket after verifying access.
func (s *AttachmentService) Upload(ctx context.Context, input UploadInput) (*domainTicket.Attachment, error) {
	// Verify ticket access
	t, err := s.tickets.GetByID(ctx, input.TicketID, input.TenantID)
	if err != nil {
		return nil, err
	}
	if !t.CanBeViewedBy(input.UserID, input.Role) {
		return nil, domainTicket.ErrForbidden
	}

	return s.attachments.Create(ctx, AttachmentUploadInput{
		TicketID:    input.TicketID,
		UserID:      input.UserID,
		Filename:    input.Filename,
		ContentType: input.ContentType,
		FileSize:    int64(len(input.FileData)),
		FileData:    input.FileData,
	})
}

// ListInput holds the parameters for listing attachments.
type ListInput struct {
	TenantID uuid.UUID
	TicketID uuid.UUID
	UserID   uuid.UUID
	Role     string
}

// List returns metadata for all attachments on a ticket after verifying access.
func (s *AttachmentService) List(ctx context.Context, input ListInput) ([]domainTicket.Attachment, error) {
	// Verify ticket access
	t, err := s.tickets.GetByID(ctx, input.TicketID, input.TenantID)
	if err != nil {
		return nil, err
	}
	if !t.CanBeViewedBy(input.UserID, input.Role) {
		return nil, domainTicket.ErrForbidden
	}
	return s.attachments.List(ctx, input.TicketID)
}

// DownloadResult holds the attachment metadata and file content for download.
type DownloadResult struct {
	Attachment *domainTicket.Attachment
	Reader     io.Reader
}

// DownloadInput holds the parameters for downloading an attachment.
type DownloadInput struct {
	TenantID     uuid.UUID
	TicketID     uuid.UUID
	AttachmentID uuid.UUID
	UserID       uuid.UUID
	Role         string
}

// Download returns the attachment metadata and file content after verifying access.
func (s *AttachmentService) Download(ctx context.Context, input DownloadInput) (*domainTicket.Attachment, []byte, error) {
	// Verify ticket access
	t, err := s.tickets.GetByID(ctx, input.TicketID, input.TenantID)
	if err != nil {
		return nil, nil, err
	}
	if !t.CanBeViewedBy(input.UserID, input.Role) {
		return nil, nil, domainTicket.ErrForbidden
	}
	return s.attachments.GetByID(ctx, input.AttachmentID, input.TicketID)
}
