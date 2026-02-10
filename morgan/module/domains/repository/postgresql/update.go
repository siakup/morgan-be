package postgresql

import (
	"context"

	"github.com/jackc/pgx/v5"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/domain"
)

var queryUpdate = `
	UPDATE organization.domains
	SET
		name = @name,
		status = @status,
		updated_by = @updated_by,
		updated_at = now()
	WHERE id = @id
`

// Update modifies an existing domain record.
func (r *Repository) Update(ctx context.Context, domain *domain.Domain) error {
	_, err := r.db.Exec(ctx, queryUpdate, pgx.NamedArgs{
		"id":         domain.Id,
		"name":       domain.Name,
		"status":     domain.Status,
		"updated_by": domain.UpdatedBy,
	})
	return err
}
