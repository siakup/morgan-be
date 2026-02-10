package postgresql

import (
	"context"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

func (r *Repository) Store(ctx context.Context, rt *domain.RoomType) error {
	query := `INSERT INTO master.room_types (id, name, description, is_active, created_at, created_by, updated_at, updated_by)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := r.db.Exec(ctx, query,
		rt.Id, rt.Name, rt.Description, rt.IsActive,
		rt.CreatedAt, rt.CreatedBy, rt.UpdatedAt, rt.UpdatedBy,
	)
	return err
}
