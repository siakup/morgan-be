package mocks

import (
	"context"

	"github.com/siakup/morgan-be/morgan/module/buildings/domain"
	"github.com/stretchr/testify/mock"
)

type BuildingUseCaseMock struct {
	mock.Mock
}

func (m *BuildingUseCaseMock) GetAll(ctx context.Context, filter domain.BuildingFilter) ([]*domain.Building, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Building), args.Get(1).(int64), args.Error(2)
}

func (m *BuildingUseCaseMock) Get(ctx context.Context, id string) (*domain.Building, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Building), args.Error(1)
}

func (m *BuildingUseCaseMock) Create(ctx context.Context, building *domain.Building) error {
	args := m.Called(ctx, building)
	return args.Error(0)
}

func (m *BuildingUseCaseMock) Update(ctx context.Context, building *domain.Building) error {
	args := m.Called(ctx, building)
	return args.Error(0)
}

func (m *BuildingUseCaseMock) Delete(ctx context.Context, id string, deletedBy string) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}
