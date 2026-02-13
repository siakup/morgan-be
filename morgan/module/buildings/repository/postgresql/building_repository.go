package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/siakup/morgan-be/morgan/module/buildings/domain"
)

type BuildingRepository struct {
	db *pgx.Conn
}

func NewBuildingRepository(db *pgx.Conn) *BuildingRepository {
	return &BuildingRepository{db}
}

func (r *BuildingRepository) FindAll(ctx context.Context, filter domain.BuildingFilter) ([]*domain.Building, int64, error) {
	query := `SELECT id, name, status, created_at, created_by, updated_at, updated_by FROM master.buildings WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM master.buildings WHERE deleted_at IS NULL`

	var args []interface{}
	argIdx := 1

	if filter.Search != "" {
		query += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		countQuery += fmt.Sprintf(" AND name ILIKE $%d", argIdx)
		args = append(args, "%"+filter.Search+"%")
		argIdx++
	}

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIdx, argIdx+1)
	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var buildings []*domain.Building
	for rows.Next() {
		var b domain.Building
		if err := rows.Scan(&b.Id, &b.Name, &b.Status, &b.CreatedAt, &b.CreatedBy, &b.UpdatedAt, &b.UpdatedBy); err != nil {
			return nil, 0, err
		}
		buildings = append(buildings, &b)
	}

	var count int64
	if filter.Search != "" {
		err = r.db.QueryRow(ctx, countQuery, "%"+filter.Search+"%").Scan(&count)
	} else {
		err = r.db.QueryRow(ctx, countQuery).Scan(&count)
	}
	if err != nil {
		return nil, 0, err
	}

	return buildings, count, nil
}

func (r *BuildingRepository) FindByID(ctx context.Context, id string) (*domain.Building, error) {
	query := `SELECT id, name, status, created_at, created_by, updated_at, updated_by FROM master.buildings WHERE id = $1 AND deleted_at IS NULL`
	var b domain.Building
	err := r.db.QueryRow(ctx, query, id).Scan(&b.Id, &b.Name, &b.Status, &b.CreatedAt, &b.CreatedBy, &b.UpdatedAt, &b.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *BuildingRepository) Store(ctx context.Context, b *domain.Building) error {
	query := `INSERT INTO master.buildings (id, name, status, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query, b.Id, b.Name, b.Status, b.CreatedAt, b.CreatedBy, b.UpdatedAt, b.UpdatedBy)
	return err
}

func (r *BuildingRepository) Update(ctx context.Context, b *domain.Building) error {
	query := `UPDATE master.buildings SET name = $1, status = $2, updated_at = $3, updated_by = $4 WHERE id = $5`
	_, err := r.db.Exec(ctx, query, b.Name, b.Status, b.UpdatedAt, b.UpdatedBy, b.Id)
	return err
}

func (r *BuildingRepository) Delete(ctx context.Context, id string, deletedBy string) error {
	query := `UPDATE master.buildings SET deleted_at = $1, deleted_by = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, time.Now(), deletedBy, id)
	return err
}
