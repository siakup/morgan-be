package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
)

func (r *Repository) FindAll(ctx context.Context, filter domain.SeverityLevelFilter) ([]*domain.SeverityLevel, int64, error) {
	var where []string
	var args []interface{}

	if filter.Search != "" {
		where = append(where, "name ILIKE $"+fmt.Sprint(len(args)+1))
		args = append(args, "%"+filter.Search+"%")
	}

	where = append(where, "deleted_at IS NULL")

	whereClause := ""
	if len(where) > 0 {
		whereClause = "WHERE " + strings.Join(where, " AND ")
	}

	countQuery := "SELECT COUNT(*) FROM hr.severity_levels " + whereClause
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := "SELECT id, name, status, created_at, created_by, updated_at, updated_by FROM hr.severity_levels " + whereClause + " ORDER BY created_at DESC LIMIT $" + fmt.Sprint(len(args)+1) + " OFFSET $" + fmt.Sprint(len(args)+2)
	args = append(args, filter.GetLimit(), filter.GetOffset())

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var severityLevels []*domain.SeverityLevel
	for rows.Next() {
		var e SeverityLevelEntity
		if err := rows.Scan(&e.Id, &e.Name, &e.Status, &e.CreatedAt, &e.CreatedBy, &e.UpdatedAt, &e.UpdatedBy); err != nil {
			return nil, 0, err
		}
		severityLevels = append(severityLevels, &domain.SeverityLevel{
			Id:        e.Id,
			Name:      e.Name,
			Status:    e.Status,
			CreatedAt: e.CreatedAt,
			CreatedBy: nullStringToPointer(e.CreatedBy),
			UpdatedAt: e.UpdatedAt,
			UpdatedBy: nullStringToPointer(e.UpdatedBy),
		})
	}

	return severityLevels, total, nil
}
