package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	appTicket "github.com/trickreport/backend/internal/application/ticket"
	domainTicket "github.com/trickreport/backend/internal/domain/ticket"
)

// AttachmentRepo implements ticket attachment persistence using PostgreSQL.
type AttachmentRepo struct {
	db *pgxpool.Pool
}

// NewAttachmentRepo creates a new AttachmentRepo.
func NewAttachmentRepo(db *pgxpool.Pool) *AttachmentRepo {
	return &AttachmentRepo{db: db}
}

// Compile-time assertion that AttachmentRepo implements appTicket.AttachmentRepository.
var _ appTicket.AttachmentRepository = (*AttachmentRepo)(nil)

// Create inserts a new attachment and returns the populated domain entity.
func (r *AttachmentRepo) Create(ctx context.Context, input appTicket.AttachmentUploadInput) (*domainTicket.Attachment, error) {
	const q = `INSERT INTO ticket_attachments (ticket_id, user_id, filename, content_type, file_size, file_data)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	var id pgtype.UUID
	var createdAt pgtype.Timestamptz
	if err := r.db.QueryRow(ctx, q, input.TicketID, input.UserID, input.Filename, input.ContentType, input.FileSize, input.FileData).
		Scan(&id, &createdAt); err != nil {
		return nil, fmt.Errorf("attachment_repo.Create: %w", err)
	}

	return &domainTicket.Attachment{
		ID:          pgToUUID(id),
		TicketID:    input.TicketID,
		UserID:      input.UserID,
		Filename:    input.Filename,
		ContentType: input.ContentType,
		FileSize:    input.FileSize,
		CreatedAt:   createdAt.Time,
	}, nil
}

// List returns metadata for all attachments on a ticket (without file data).
func (r *AttachmentRepo) List(ctx context.Context, ticketID uuid.UUID) ([]domainTicket.Attachment, error) {
	const q = `SELECT a.id, a.ticket_id, a.user_id, a.filename, a.content_type, a.file_size, a.created_at, u.name
		FROM ticket_attachments a
		JOIN users u ON a.user_id = u.id
		WHERE a.ticket_id = $1
		ORDER BY a.created_at DESC`

	rows, err := r.db.Query(ctx, q, ticketID)
	if err != nil {
		return nil, fmt.Errorf("attachment_repo.List: query: %w", err)
	}
	defer rows.Close()

	var attachments []domainTicket.Attachment
	for rows.Next() {
		var a domainTicket.Attachment
		var id, ticketIDCol, userID pgtype.UUID
		if err := rows.Scan(&id, &ticketIDCol, &userID, &a.Filename, &a.ContentType, &a.FileSize, &a.CreatedAt, &a.UserName); err != nil {
			return nil, fmt.Errorf("attachment_repo.List: scan: %w", err)
		}
		a.ID = pgToUUID(id)
		a.TicketID = pgToUUID(ticketIDCol)
		a.UserID = pgToUUID(userID)
		attachments = append(attachments, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("attachment_repo.List: rows: %w", err)
	}
	return attachments, nil
}

// GetByID returns a single attachment including file data for download.
func (r *AttachmentRepo) GetByID(ctx context.Context, id, ticketID uuid.UUID) (*domainTicket.Attachment, []byte, error) {
	const q = `SELECT id, ticket_id, user_id, filename, content_type, file_size, created_at, file_data
		FROM ticket_attachments
		WHERE id = $1 AND ticket_id = $2`

	var a domainTicket.Attachment
	var idCol, ticketIDCol, userID pgtype.UUID
	var fileData []byte
	err := r.db.QueryRow(ctx, q, id, ticketID).Scan(&idCol, &ticketIDCol, &userID, &a.Filename, &a.ContentType, &a.FileSize, &a.CreatedAt, &fileData)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, domainTicket.ErrNotFound
		}
		return nil, nil, fmt.Errorf("attachment_repo.GetByID: %w", err)
	}
	a.ID = pgToUUID(idCol)
	a.TicketID = pgToUUID(ticketIDCol)
	a.UserID = pgToUUID(userID)
	return &a, fileData, nil
}
