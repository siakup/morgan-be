package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/object"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_sessions/domain"
)

var queryFindById = `
	SELECT
		id, name, start, end, status
	FROM schedule.shift_sessions
	WHERE id = @id
	LIMIT 1
`

// FindByID retrieves a single shift session by its ID.
func (r *Repository) FindByID(ctx context.Context, id string) (*domain.ShiftSession, error) {
	rows, err := r.db.Query(ctx, queryFindById, pgx.NamedArgs{
		"id": id,
	})
	if err != nil {
		return nil, err
	}

	record, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[ShiftSessionEntity])
	if err != nil {
		return nil, err
	}

	shiftSession, err := object.Parse[*ShiftSessionEntity, *domain.ShiftSession](object.TagDB, object.TagObject, record)
	if err != nil {
		return nil, err
	}

	return shiftSession, nil
}
