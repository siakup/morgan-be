package usecase

import (
	"context"
	"time"

	"github.com/siakup/morgan-be/morgan/module/buildings/domain"
)

func (u *UseCase) Update(ctx context.Context, building *domain.Building) error {
	ctx, span := u.tracer.Start(ctx, "Update")
	defer span.End()

	building.UpdatedAt = time.Now()
	return u.repository.Update(ctx, building)
}
