package postgresql

import (
	"context"
	"fmt"
	"strings"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

func (r *Repository) FindAll(ctx context.Context, filter domain.RoomTypeFilter) ([]*domain.RoomType, int64, error) {
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

	countQuery := "SELECT COUNT(*) FROM master.room_types " + whereClause
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, name, description, is_active, created_at, created_by, updated_at, updated_by
	FROM master.room_types ` + whereClause + `
	ORDER BY created_at DESC LIMIT $` + fmt.Sprint(len(args)+1) + ` OFFSET $` + fmt.Sprint(len(args)+2)
	args = append(args, filter.GetLimit(), filter.GetOffset())

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var list []*domain.RoomType
	for rows.Next() {
		var e RoomTypeEntity
		if err := rows.Scan(&e.Id, &e.Name, &e.Description, &e.IsActive, &e.CreatedAt, &e.CreatedBy, &e.UpdatedAt, &e.UpdatedBy); err != nil {
			return nil, 0, err
		}
		list = append(list, &domain.RoomType{
			Id:          e.Id,
			Name:        e.Name,
			Description: e.Description,
			IsActive:    e.IsActive,
			CreatedAt:   e.CreatedAt,
			CreatedBy:   nullStringToPointer(e.CreatedBy),
			UpdatedAt:   e.UpdatedAt,
			UpdatedBy:   nullStringToPointer(e.UpdatedBy),
		})
	}

	return list, total, nil
}
