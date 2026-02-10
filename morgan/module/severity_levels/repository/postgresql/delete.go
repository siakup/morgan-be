package postgresql

import (
	"context"
)

func (r *Repository) Delete(ctx context.Context, id string, deletedBy string) error {
	query := "UPDATE hr.severity_levels SET deleted_at = NOW(), deleted_by = $2 WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id, deletedBy)
	return err
}
