package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/domain"
)

var queryUpdate = `
	UPDATE hr.shift_sessions
	SET
		name = @name,
		start = @start,
		"end" = @end,
		status = @status,
		updated_by = @updated_by,
		updated_at = now()
	WHERE id = @id
`

// Update modifies an existing shift session record.
func (r *Repository) Update(ctx context.Context, shiftSession *domain.ShiftSession) error {
	_, err := r.db.Exec(ctx, queryUpdate, pgx.NamedArgs{
		"id":         shiftSession.Id,
		"name":       shiftSession.Name,
		"start":      shiftSession.Start,
		"end":        shiftSession.End,
		"status":     shiftSession.Status,
		"updated_by": shiftSession.UpdatedBy,
	})
	return err
}
