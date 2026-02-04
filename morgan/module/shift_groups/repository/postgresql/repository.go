package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

var _ domain.ShiftGroupRepository = (*Repository)(nil)

// ShiftGroupEntity represents the schema in the database.
type ShiftGroupEntity struct {
	Id        string         `db:"id"`
	Name      string         `db:"name"`
	Status    bool           `db:"status"`
	CreatedAt time.Time      `db:"created_at"`
	CreatedBy sql.NullString `db:"created_by"`
	UpdatedAt time.Time      `db:"updated_at"`
	UpdatedBy sql.NullString `db:"updated_by"`
	DeletedAt sql.NullTime   `db:"deleted_at"`
	DeletedBy sql.NullString `db:"deleted_by"`
}

// Repository implements the domain.ShiftGroupRepository interface for PostgreSQL.
type Repository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new instance of the PostgreSQL Repository.
func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindAll(ctx context.Context, filter domain.ShiftGroupFilter) ([]*domain.ShiftGroup, int64, error) {
	var where []string
	var args []interface{}

	// InstitutionId filter removed as per ERD alignment

	if filter.Search != "" {
		where = append(where, "name ILIKE $"+fmt.Sprint(len(args)+1))
		args = append(args, "%"+filter.Search+"%")
	}

	where = append(where, "deleted_at IS NULL")

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM hr.shift_groups " + whereClause
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT id, name, status, created_at, created_by, updated_at, updated_by FROM hr.shift_groups " + whereClause + " ORDER BY created_at DESC LIMIT $" + fmt.Sprint(len(args)+1) + " OFFSET $" + fmt.Sprint(len(args)+2)
	args = append(args, filter.GetLimit(), filter.GetOffset())

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var shiftGroups []*domain.ShiftGroup
	for rows.Next() {
		var e ShiftGroupEntity
		if err := rows.Scan(&e.Id, &e.Name, &e.Status, &e.CreatedAt, &e.CreatedBy, &e.UpdatedAt, &e.UpdatedBy); err != nil {
			return nil, 0, err
		}
		shiftGroups = append(shiftGroups, &domain.ShiftGroup{
			Id:        e.Id,
			Name:      e.Name,
			Status:    e.Status,
			CreatedAt: e.CreatedAt,
			CreatedBy: nullStringToPointer(e.CreatedBy),
			UpdatedAt: e.UpdatedAt,
			UpdatedBy: nullStringToPointer(e.UpdatedBy),
		})
	}

	return shiftGroups, total, nil
}

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

func (r *Repository) Store(ctx context.Context, s *domain.ShiftGroup) error {
	query := "INSERT INTO hr.shift_groups (id, name, status, created_at, created_by, updated_at, updated_by) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := r.db.Exec(ctx, query, s.Id, s.Name, s.Status, s.CreatedAt, s.CreatedBy, s.UpdatedAt, s.UpdatedBy)
	return err
}

func (r *Repository) Update(ctx context.Context, s *domain.ShiftGroup) error {
	query := "UPDATE hr.shift_groups SET name = $1, status = $2, updated_at = $3, updated_by = $4 WHERE id = $5"
	_, err := r.db.Exec(ctx, query, s.Name, s.Status, s.UpdatedAt, s.UpdatedBy, s.Id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id string, deletedBy string) error {
	// Soft delete
	query := "UPDATE hr.shift_groups SET deleted_at = NOW(), deleted_by = $2 WHERE id = $1"
	_, err := r.db.Exec(ctx, query, id, deletedBy)
	return err
}

func nullStringToPointer(ns sql.NullString) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func nullTimeToPointer(nt sql.NullTime) *time.Time {
	if nt.Valid {
		return &nt.Time
	}
	return nil
}
