package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_sessions/domain"
)

var queryStore = `
	INSERT INTO schedule.shift_sessions (
		name, start, "end", created_by, updated_by
	) VALUES (
		@name, @start, @end, @created_by, @updated_by
	)
	RETURNING id
`

// Store persists a new shift session to the database.
func (r *Repository) Store(ctx context.Context, shiftSession *domain.ShiftSession) error {
	rows, err := r.db.Query(ctx, queryStore, pgx.NamedArgs{
		"name":       shiftSession.Name,
		"start":      shiftSession.Start,
		"end":        shiftSession.End,
		"created_by": shiftSession.CreatedBy,
		"updated_by": shiftSession.UpdatedBy,
	})
	if err != nil {
		return err
	}

	// Scan returning ID
	var id string
	if _, err := pgx.ForEachRow(rows, []any{&id}, func() error { return nil }); err != nil {
		return err
	}
	shiftSession.Id = id
	return nil
}
