package postgres

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type MFAChallengeRepo struct{ db *pgxpool.Pool }

func NewMFAChallengeRepo(db *pgxpool.Pool) *MFAChallengeRepo { return &MFAChallengeRepo{db: db} }

func (r *MFAChallengeRepo) CreateChallenge(ctx context.Context, jti string, userID uuid.UUID, expires time.Time) error {
	_, err := r.db.Exec(ctx, `INSERT INTO mfa_challenges(jti,user_id,expires_at) VALUES($1,$2,$3)`, jti, userID, expires)
	return err
}

func (r *MFAChallengeRepo) ConsumeChallenge(ctx context.Context, jti string, userID uuid.UUID, valid bool) (bool, error) {
	result, err := r.db.Exec(ctx, `UPDATE mfa_challenges SET attempts=attempts+1, consumed=$3 WHERE jti=$1 AND user_id=$2 AND expires_at>NOW() AND NOT consumed AND attempts<5`, jti, userID, valid)
	return valid && result.RowsAffected() == 1, err
}
