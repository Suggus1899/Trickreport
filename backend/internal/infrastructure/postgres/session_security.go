package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	appAuth "github.com/trickreport/backend/internal/application/auth"
)

func (r *SessionRepo) FindActiveByAccessJTI(ctx context.Context, jti string) (*appAuth.SessionInput, error) {
	var s appAuth.SessionInput
	err := r.db.QueryRow(ctx, `SELECT user_id,jti,COALESCE(refresh_jti,''),COALESCE(ip_address,''),COALESCE(user_agent,''),expires_at FROM user_sessions WHERE jti=$1 AND NOT revoked AND expires_at>NOW()`, jti).Scan(&s.UserID, &s.JTI, &s.RefreshJTI, &s.IPAddress, &s.UserAgent, &s.ExpiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, appAuth.ErrInvalidCredentials
	}
	return &s, err
}

func (r *SessionRepo) RotateSession(ctx context.Context, oldRefreshJTI string, s appAuth.SessionInput) (bool, error) {
	result, err := r.db.Exec(ctx, `UPDATE user_sessions SET jti=$1,refresh_jti=$2,ip_address=$3,user_agent=$4,expires_at=$5 WHERE refresh_jti=$6 AND user_id=$7 AND NOT revoked AND expires_at>NOW()`, s.JTI, s.RefreshJTI, s.IPAddress, s.UserAgent, s.ExpiresAt, oldRefreshJTI, s.UserID)
	return result.RowsAffected() == 1, err
}

func (r *SessionRepo) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `UPDATE user_sessions SET revoked=TRUE WHERE user_id=$1 AND NOT revoked`, userID)
	return err
}
