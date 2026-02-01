package postgresql

import (
	"context"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

func (r *Repository) Update(ctx context.Context, s *domain.ShiftGroup) error {
	query := "UPDATE hr.shift_groups SET name = $1, status = $2, updated_at = $3, updated_by = $4 WHERE id = $5"
	_, err := r.db.Exec(ctx, query, s.Name, s.Status, s.UpdatedAt, s.UpdatedBy, s.Id)
	return err
}
