package usecase

import (
	"context"
	"time"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/severity_levels/domain"
)

func (u *UseCase) Update(ctx context.Context, sl *domain.SeverityLevel) error {
	sl.UpdatedAt = time.Now()
	// Should we check ID empty? Usually redundant if passed in struct.
	return u.repo.Update(ctx, sl)
}
