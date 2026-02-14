package usecase

import (
	"context"
	"time"

	"github.com/siakup/morgan-be/morgan/module/divisions/domain"
)

func (u *UseCase) Update(ctx context.Context, division *domain.Division) error {
	ctx, span := u.tracer.Start(ctx, "Update")
	defer span.End()

	division.UpdatedAt = time.Now()
	return u.repository.Update(ctx, division)
}
