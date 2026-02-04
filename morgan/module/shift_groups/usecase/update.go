package usecase

import (
	"context"
	"time"

	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

func (u *UseCase) Update(ctx context.Context, shiftGroup *domain.ShiftGroup) error {
	if len(shiftGroup.Name) > 2 {
		return &errors.AppError{Code: 400, Type: "BAD_REQUEST", Message: "Name too long, max 2 characters"}
	}
	shiftGroup.UpdatedAt = time.Now()
	return u.repo.Update(ctx, shiftGroup)
}
