package postgres

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	appAuth "github.com/trickreport/backend/internal/application/auth"
)

// SessionRepo implements appAuth.SessionRepository using PostgreSQL.
type SessionRepo struct {
	db *pgxpool.Pool
}

// NewSessionRepo creates a new SessionRepo.
func NewSessionRepo(db *pgxpool.Pool) *SessionRepo {
	return &SessionRepo{db: db}
}

// CreateSession inserts a new session record.
func (r *SessionRepo) CreateSession(ctx context.Context, input appAuth.SessionInput) error {
	const q = `INSERT INTO user_sessions (user_id, jti, refresh_jti, ip_address, user_agent, expires_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := r.db.Exec(ctx, q, input.UserID, input.JTI, input.RefreshJTI, input.IPAddress, input.UserAgent, input.ExpiresAt)
	if err != nil {
		return fmt.Errorf("session_repo.CreateSession: %w", err)
	}
	return nil
}

// RevokeSession marks a session as revoked by its access token jti.
func (r *SessionRepo) RevokeSession(ctx context.Context, jti string) error {
	const q = `UPDATE user_sessions SET revoked = TRUE WHERE jti = $1`
	_, err := r.db.Exec(ctx, q, jti)
	if err != nil {
		return fmt.Errorf("session_repo.RevokeSession: %w", err)
	}
	return nil
}

// RevokeByRefreshJTI marks a session as revoked by its refresh token jti.
func (r *SessionRepo) RevokeByRefreshJTI(ctx context.Context, refreshJTI string) error {
	const q = `UPDATE user_sessions SET revoked = TRUE WHERE refresh_jti = $1`
	_, err := r.db.Exec(ctx, q, refreshJTI)
	if err != nil {
		return fmt.Errorf("session_repo.RevokeByRefreshJTI: %w", err)
	}
	return nil
}

// Compile-time interface assertion.
var _ appAuth.SessionRepository = (*SessionRepo)(nil)

// PurgeExpiredSessions deletes sessions that have expired. This is a
// maintenance method that can be called periodically.
func (r *SessionRepo) PurgeExpiredSessions(ctx context.Context) error {
	const q = `DELETE FROM user_sessions WHERE expires_at < NOW()`
	_, err := r.db.Exec(ctx, q)
	if err != nil {
		return fmt.Errorf("session_repo.PurgeExpiredSessions: %w", err)
	}
	return nil
}

// GetActiveSessions returns all active (non-expired, non-revoked) sessions
// for a user.
func (r *SessionRepo) GetActiveSessions(ctx context.Context, userID uuid.UUID) ([]appAuth.SessionInput, error) {
	const q = `SELECT jti, refresh_jti, ip_address, user_agent, expires_at FROM user_sessions WHERE user_id = $1 AND revoked = FALSE AND expires_at > NOW() ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, q, userID)
	if err != nil {
		return nil, fmt.Errorf("session_repo.GetActiveSessions: %w", err)
	}
	defer rows.Close()

	var sessions []appAuth.SessionInput
	for rows.Next() {
		var s appAuth.SessionInput
		s.UserID = userID
		if err := rows.Scan(&s.JTI, &s.RefreshJTI, &s.IPAddress, &s.UserAgent, &s.ExpiresAt); err != nil {
			return nil, fmt.Errorf("session_repo.GetActiveSessions: scan: %w", err)
		}
		sessions = append(sessions, s)
	}
	return sessions, rows.Err()
}
