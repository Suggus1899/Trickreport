package postgres

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	appAuth "github.com/trickreport/backend/internal/application/auth"
)

func (r *PasswordResetRepo) ConsumeAndUpdatePassword(ctx context.Context, token, passwordHash string) error {
	return withTx(ctx, r.db, func(ctx context.Context) error {
		tx, _ := TxFromContext(ctx)
		var userID uuid.UUID
		err := tx.QueryRow(ctx, `SELECT user_id FROM password_reset_tokens WHERE token=$1 AND NOT used AND expires_at>NOW() FOR UPDATE`, token).Scan(&userID)
		if errors.Is(err, pgx.ErrNoRows) {
			return appAuth.ErrPasswordResetTokenInvalid
		}
		if err != nil {
			return err
		}
		result, err := tx.Exec(ctx, `UPDATE users SET password=$1,updated_at=NOW() WHERE id=$2 AND active`, passwordHash, userID)
		if err != nil {
			return err
		}
		if result.RowsAffected() != 1 {
			return appAuth.ErrPasswordResetTokenInvalid
		}
		if _, err := tx.Exec(ctx, `UPDATE password_reset_tokens SET used=TRUE WHERE user_id=$1 AND NOT used`, userID); err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `UPDATE user_sessions SET revoked=TRUE WHERE user_id=$1 AND NOT revoked`, userID)
		return err
	})
}
