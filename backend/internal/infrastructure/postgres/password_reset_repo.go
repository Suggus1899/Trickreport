package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	appAuth "github.com/trickreport/backend/internal/application/auth"
)

// PasswordResetRepo implements appAuth.PasswordResetRepository using PostgreSQL.
type PasswordResetRepo struct {
	db *pgxpool.Pool
}

// NewPasswordResetRepo creates a new PasswordResetRepo.
func NewPasswordResetRepo(db *pgxpool.Pool) *PasswordResetRepo {
	return &PasswordResetRepo{db: db}
}

// CreateToken stores a new password reset token.
func (r *PasswordResetRepo) CreateToken(ctx context.Context, userID uuid.UUID, token string, expiresAt time.Time) error {
	const q = `INSERT INTO password_reset_tokens (user_id, token, expires_at) VALUES ($1, $2, $3)`
	_, err := r.db.Exec(ctx, q, userID, token, expiresAt)
	if err != nil {
		return fmt.Errorf("password_reset_repo.CreateToken: %w", err)
	}
	return nil
}

// FindValidToken looks up a reset token that is valid (not used, not expired).
// Returns the associated user ID and expiration time.
func (r *PasswordResetRepo) FindValidToken(ctx context.Context, token string) (uuid.UUID, time.Time, error) {
	const q = `SELECT user_id, expires_at FROM password_reset_tokens WHERE token = $1 AND used = FALSE LIMIT 1`

	var userID pgtype.UUID
	var expiresAt time.Time
	err := r.db.QueryRow(ctx, q, token).Scan(&userID, &expiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uuid.Nil, time.Time{}, appAuth.ErrPasswordResetTokenInvalid
		}
		return uuid.Nil, time.Time{}, fmt.Errorf("password_reset_repo.FindValidToken: %w", err)
	}
	return pgToUUID(userID), expiresAt, nil
}

// MarkUsed marks a reset token as used.
func (r *PasswordResetRepo) MarkUsed(ctx context.Context, token string) error {
	const q = `UPDATE password_reset_tokens SET used = TRUE WHERE token = $1`
	_, err := r.db.Exec(ctx, q, token)
	if err != nil {
		return fmt.Errorf("password_reset_repo.MarkUsed: %w", err)
	}
	return nil
}

// Compile-time interface assertion.
var _ appAuth.PasswordResetRepository = (*PasswordResetRepo)(nil)
