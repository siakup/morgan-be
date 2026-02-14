package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/siakup/morgan-be/morgan/module/buildings/domain"
)

func (u *UseCase) Create(ctx context.Context, building *domain.Building) error {
	ctx, span := u.tracer.Start(ctx, "Create")
	defer span.End()

	building.Id = uuid.NewString()
	now := time.Now()
	building.CreatedAt = now
	building.UpdatedAt = now

	return u.repository.Store(ctx, building)
}
