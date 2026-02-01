package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.ShiftGroup, error) {
	query := "SELECT id, name, status, created_at, created_by, updated_at, updated_by, deleted_at, deleted_by FROM hr.shift_groups WHERE id = $1 AND deleted_at IS NULL"
	var e ShiftGroupEntity
	err := r.db.QueryRow(ctx, query, id).Scan(&e.Id, &e.Name, &e.Status, &e.CreatedAt, &e.CreatedBy, &e.UpdatedAt, &e.UpdatedBy, &e.DeletedAt, &e.DeletedBy)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // Or custom error
		}
		return nil, err
	}
	return &domain.ShiftGroup{
		Id:        e.Id,
		Name:      e.Name,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
		CreatedBy: nullStringToPointer(e.CreatedBy),
		UpdatedAt: e.UpdatedAt,
		UpdatedBy: nullStringToPointer(e.UpdatedBy),
		DeletedAt: nullTimeToPointer(e.DeletedAt),
		DeletedBy: nullStringToPointer(e.DeletedBy),
	}, nil
}
