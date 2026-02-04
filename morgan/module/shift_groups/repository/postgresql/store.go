package postgresql

import (
	"context"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

func (r *Repository) Store(ctx context.Context, s *domain.ShiftGroup) error {
	query := "INSERT INTO hr.shift_groups (id, name, status, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(ctx, query, s.Id, s.Name, s.Status, s.CreatedAt, s.CreatedBy, s.UpdatedAt, s.UpdatedBy)
	return err
}
