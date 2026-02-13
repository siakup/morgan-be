package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/siakup/morgan-be/morgan/module/divisions/domain"
)

func (u *UseCase) Create(ctx context.Context, division *domain.Division) error {
	ctx, span := u.tracer.Start(ctx, "Create")
	defer span.End()

	division.Id = uuid.NewString()
	now := time.Now()
	division.CreatedAt = now
	division.UpdatedAt = now

	return u.repository.Store(ctx, division)
}
