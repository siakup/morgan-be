package usecase

import (
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/room_type/domain"
)

var _ domain.UseCase = (*UseCase)(nil)

// UseCase implements the domain.UseCase interface for room type.
type UseCase struct {
	repo domain.RoomTypeRepository
}

// NewUseCase creates a new instance of the UseCase.
func NewUseCase(repo domain.RoomTypeRepository) *UseCase {
	return &UseCase{repo: repo}
}
