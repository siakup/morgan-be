package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/severity_levels/domain"
)

func (u *UseCase) Create(ctx context.Context, sl *domain.SeverityLevel) error {
	sl.Id = uuid.NewString()
	now := time.Now()
	sl.CreatedAt = now
	sl.UpdatedAt = now
	return u.repo.Store(ctx, sl)
}
