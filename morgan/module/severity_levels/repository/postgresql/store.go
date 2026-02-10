package postgresql

import (
	"context"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/severity_levels/domain"
)

func (r *Repository) Store(ctx context.Context, s *domain.SeverityLevel) error {
	query := "INSERT INTO hr.severity_levels (id, name, status, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(ctx, query, s.Id, s.Name, s.Status, s.CreatedAt, s.CreatedBy, s.UpdatedAt, s.UpdatedBy)
	return err
}
