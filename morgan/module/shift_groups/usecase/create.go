package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

func (u *UseCase) Create(ctx context.Context, shiftGroup *domain.ShiftGroup) error {
	if len(shiftGroup.Name) > 2 {
		return &errors.AppError{Code: 400, Type: "BAD_REQUEST", Message: "Name too long, max 2 characters"}
	}
	shiftGroup.Id = uuid.NewString()
	now := time.Now()
	shiftGroup.CreatedAt = now
	shiftGroup.UpdatedAt = now
	// Add validation here if needed
	return u.repo.Store(ctx, shiftGroup)
}
