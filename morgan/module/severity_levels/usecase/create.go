package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
)

func (u *UseCase) Create(ctx context.Context, sl *domain.SeverityLevel) error {
	sl.Id = uuid.NewString()
	now := time.Now()
	sl.CreatedAt = now
	sl.UpdatedAt = now
	return u.repository.Store(ctx, sl)
}
