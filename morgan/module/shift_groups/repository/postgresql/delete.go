package postgresql

import (
	"context"
)

func (r *Repository) Delete(ctx context.Context, id string, deletedBy string) error {
	// Soft delete
	query := "UPDATE hr.shift_groups SET deleted_at = NOW(), deleted_by = $2 WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id, deletedBy)
	return err
}
