package postgresql

import (
	"context"

	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
)

func (r *Repository) Update(ctx context.Context, s *domain.SeverityLevel) error {
	query := "UPDATE hr.severity_levels SET name = $2, status = $3, updated_at = $4, updated_by = $5 WHERE id = $1"
	_, err := r.db.Exec(ctx, query, s.Id, s.Name, s.Status, s.UpdatedAt, s.UpdatedBy)
	return err
}
