package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

var queryDelete = `UPDATE master.domains
SET
    deleted_at = NOW(),
    deleted_by = @deleted_by
WHERE id = @id
AND deleted_at IS NULL
`

// Delete soft removes a domain from the database.
func (r *Repository) Delete(ctx context.Context, id string, deletedBy string) error {

	_, err := r.db.Exec(ctx, queryDelete, pgx.NamedArgs{
		"id":         id,
		"deleted_by": deletedBy,
	})

	return err
}
