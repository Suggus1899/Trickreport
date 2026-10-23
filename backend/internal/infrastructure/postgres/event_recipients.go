package postgres

import (
	"context"
	"github.com/google/uuid"
)

func (r *UserRepo) EventRecipients(ctx context.Context, tenantID, ticketID uuid.UUID, includeCreator bool) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx, `SELECT u.id FROM users u JOIN tickets t ON t.tenant_id=u.tenant_id WHERE t.id=$1 AND t.tenant_id=$2 AND u.active AND (u.role IN ('admin','agent') OR ($3 AND u.id=t.created_by))`, ticketID, tenantID, includeCreator)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := []uuid.UUID{}
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		users = append(users, id)
	}
	return users, rows.Err()
}
