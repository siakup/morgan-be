package postgresql

import (
	"context"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

func (r *Repository) Update(ctx context.Context, rt *domain.RoomType) error {
	query := `UPDATE master.room_types SET name = $1, description = $2, is_active = $3, updated_at = $4, updated_by = $5 WHERE id = $6`
	_, err := r.db.Exec(ctx, query, rt.Name, rt.Description, rt.IsActive, rt.UpdatedAt, rt.UpdatedBy, rt.Id)
	return err
}
