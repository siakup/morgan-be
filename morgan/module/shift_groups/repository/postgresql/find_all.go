package postgresql

import (
	"context"
	"fmt"
	"strings"

	"github.com/siakup/morgan-be/morgan/module/shift_groups/domain"
)

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
