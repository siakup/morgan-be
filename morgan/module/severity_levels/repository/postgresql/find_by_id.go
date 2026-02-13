package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
)

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.SeverityLevel, error) {
	query := "SELECT id, name, status, created_at, created_by, updated_at, updated_by FROM hr.severity_levels WHERE id = $1 AND deleted_at IS NULL"
	var e SeverityLevelEntity
	err := r.db.QueryRow(ctx, query, id).Scan(&e.Id, &e.Name, &e.Status, &e.CreatedAt, &e.CreatedBy, &e.UpdatedAt, &e.UpdatedBy)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &domain.SeverityLevel{
		Id:        e.Id,
		Name:      e.Name,
		Status:    e.Status,
		CreatedAt: e.CreatedAt,
		CreatedBy: nullStringToPointer(e.CreatedBy),
		UpdatedAt: e.UpdatedAt,
		UpdatedBy: nullStringToPointer(e.UpdatedBy),
	}, nil
}
