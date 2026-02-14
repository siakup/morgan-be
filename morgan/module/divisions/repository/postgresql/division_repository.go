package postgresql

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/siakup/morgan-be/morgan/module/divisions/domain"
)

type DivisionRepository struct {
	db *pgx.Conn
}

func NewDivisionRepository(db *pgx.Conn) *DivisionRepository {
	return &DivisionRepository{db}
}

func (r *DivisionRepository) FindAll(ctx context.Context, filter domain.DivisionFilter) ([]*domain.Division, int64, error) {
	query := `SELECT id, name, status, created_at, created_by, updated_at, updated_by FROM master.divisions WHERE deleted_at IS NULL`
	countQuery := `SELECT COUNT(*) FROM master.divisions WHERE deleted_at IS NULL`

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

	var divisions []*domain.Division
	for rows.Next() {
		var d domain.Division
		if err := rows.Scan(&d.Id, &d.Name, &d.Status, &d.CreatedAt, &d.CreatedBy, &d.UpdatedAt, &d.UpdatedBy); err != nil {
			return nil, 0, err
		}
		divisions = append(divisions, &d)
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

	return divisions, count, nil
}

func (r *DivisionRepository) FindByID(ctx context.Context, id string) (*domain.Division, error) {
	query := `SELECT id, name, status, created_at, created_by, updated_at, updated_by FROM master.divisions WHERE id = $1 AND deleted_at IS NULL`
	var d domain.Division
	err := r.db.QueryRow(ctx, query, id).Scan(&d.Id, &d.Name, &d.Status, &d.CreatedAt, &d.CreatedBy, &d.UpdatedAt, &d.UpdatedBy)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *DivisionRepository) Store(ctx context.Context, d *domain.Division) error {
	query := `INSERT INTO master.divisions (id, name, status, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.db.Exec(ctx, query, d.Id, d.Name, d.Status, d.CreatedAt, d.CreatedBy, d.UpdatedAt, d.UpdatedBy)
	return err
}

func (r *DivisionRepository) Update(ctx context.Context, d *domain.Division) error {
	query := `UPDATE master.divisions SET name = $1, status = $2, updated_at = $3, updated_by = $4 WHERE id = $5`
	_, err := r.db.Exec(ctx, query, d.Name, d.Status, d.UpdatedAt, d.UpdatedBy, d.Id)
	return err
}

func (r *DivisionRepository) Delete(ctx context.Context, id string, deletedBy string) error {
	query := `UPDATE master.divisions SET deleted_at = $1, deleted_by = $2 WHERE id = $3`
	_, err := r.db.Exec(ctx, query, time.Now(), deletedBy, id)
	return err
}
