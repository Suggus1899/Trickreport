package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	infraEmail "github.com/trickreport/backend/internal/infrastructure/email"
	"time"
)

type EmailOutboxRepo struct{ db *pgxpool.Pool }

func NewEmailOutboxRepo(db *pgxpool.Pool) *EmailOutboxRepo { return &EmailOutboxRepo{db: db} }

func (r *EmailOutboxRepo) Enqueue(ctx context.Context, to, subject, body string) error {
	conn := DBTX(r.db)
	if tx, ok := TxFromContext(ctx); ok {
		conn = tx
	}
	_, err := conn.Exec(ctx, `INSERT INTO email_outbox(recipient,subject,body) VALUES($1,$2,$3)`, to, subject, body)
	return err
}

func (r *EmailOutboxRepo) Claim(ctx context.Context) (*infraEmail.QueuedEmail, error) {
	var msg infraEmail.QueuedEmail
	err := r.db.QueryRow(ctx, `UPDATE email_outbox SET status='sending',attempts=attempts+1,leased_until=NOW()+INTERVAL '1 minute' WHERE id=(SELECT id FROM email_outbox WHERE (status='pending' AND available_at<=NOW()) OR (status='sending' AND leased_until<NOW()) ORDER BY created_at FOR UPDATE SKIP LOCKED LIMIT 1) RETURNING id,recipient,subject,body,attempts`).Scan(&msg.ID, &msg.To, &msg.Subject, &msg.Body, &msg.Attempts)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &msg, err
}

func (r *EmailOutboxRepo) Finish(ctx context.Context, msg *infraEmail.QueuedEmail, sent bool) error {
	status := "pending"
	reason := "delivery failed"
	var delivered *time.Time
	if sent {
		status = "sent"
		reason = ""
		now := time.Now()
		delivered = &now
	} else if msg.Attempts >= 5 {
		status = "failed"
	}
	_, err := r.db.Exec(ctx, `UPDATE email_outbox SET status=$1,last_error=$2,available_at=$3,leased_until=NULL,delivered_at=$4,body=CASE WHEN $1='sent' THEN '' ELSE body END WHERE id=$5 AND status='sending' AND attempts=$6`, status, reason, time.Now().Add(time.Second<<min(msg.Attempts, 6)), delivered, msg.ID, msg.Attempts)
	return err
}

func (r *EmailOutboxRepo) Counts(ctx context.Context) (pending, failed int64, err error) {
	err = r.db.QueryRow(ctx, `SELECT COUNT(*) FILTER(WHERE status IN ('pending','sending')),COUNT(*) FILTER(WHERE status='failed') FROM email_outbox`).Scan(&pending, &failed)
	return
}
