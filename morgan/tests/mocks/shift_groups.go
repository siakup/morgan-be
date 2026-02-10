package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

// ShiftGroupsUseCaseMock is a mock implementation of domain.UseCase
type ShiftGroupsUseCaseMock struct {
	mock.Mock
}

// FindAll mocks the FindAll method
func (m *ShiftGroupsUseCaseMock) FindAll(ctx context.Context, filter domain.ShiftGroupFilter) ([]*domain.ShiftGroup, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.ShiftGroup), args.Get(1).(int64), args.Error(2)
}

// FindByID mocks the FindByID method
func (m *ShiftGroupsUseCaseMock) FindByID(ctx context.Context, id string) (*domain.ShiftGroup, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ShiftGroup), args.Error(1)
}

// Create mocks the Create method
func (m *ShiftGroupsUseCaseMock) Create(ctx context.Context, shiftGroup *domain.ShiftGroup) error {
	args := m.Called(ctx, shiftGroup)
	return args.Error(0)
}

// Update mocks the Update method
func (m *ShiftGroupsUseCaseMock) Update(ctx context.Context, shiftGroup *domain.ShiftGroup) error {
	args := m.Called(ctx, shiftGroup)
	return args.Error(0)
}

// Delete mocks the Delete method
func (m *ShiftGroupsUseCaseMock) Delete(ctx context.Context, id string, deletedBy string) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}
