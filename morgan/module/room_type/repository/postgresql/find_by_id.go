package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.RoomType, error) {
	query := `SELECT id, name, description, is_active, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by
	FROM master.room_types WHERE id = $1 AND deleted_at IS NULL`
	var e RoomTypeEntity
	err := r.db.QueryRow(ctx, query, id).Scan(
		&e.Id, &e.Name, &e.Description, &e.IsActive,
		&e.CreatedAt, &e.CreatedBy, &e.UpdatedAt, &e.UpdatedBy,
		&e.DeletedAt, &e.DeletedBy,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &domain.RoomType{
		Id:          e.Id,
		Name:        e.Name,
		Description: e.Description,
		IsActive:    e.IsActive,
		CreatedAt:   e.CreatedAt,
		CreatedBy:   nullStringToPointer(e.CreatedBy),
		UpdatedAt:   e.UpdatedAt,
		UpdatedBy:   nullStringToPointer(e.UpdatedBy),
		DeletedAt:   nullTimeToPointer(e.DeletedAt),
		DeletedBy:   nullStringToPointer(e.DeletedBy),
	}, nil
}
