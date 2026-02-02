package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var queryDelete = `DELETE FROM schedule.shift_sessions WHERE id = @id`

// Delete removes a shift session from the database.
func (r *Repository) Delete(ctx context.Context, id string) error {

	_, err := r.db.Exec(ctx, queryDelete, pgx.NamedArgs{
		"id": id,
	})

	return err
}
