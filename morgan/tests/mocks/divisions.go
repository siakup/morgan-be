package mocks

import (
	"context"

	"github.com/siakup/morgan-be/morgan/module/divisions/domain"
	"github.com/stretchr/testify/mock"
)

type DivisionUseCaseMock struct {
	mock.Mock
}

func (m *DivisionUseCaseMock) GetAll(ctx context.Context, filter domain.DivisionFilter) ([]*domain.Division, int64, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Get(1).(int64), args.Error(2)
	}
	return args.Get(0).([]*domain.Division), args.Get(1).(int64), args.Error(2)
}

func (m *DivisionUseCaseMock) Get(ctx context.Context, id string) (*domain.Division, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Division), args.Error(1)
}

func (m *DivisionUseCaseMock) Create(ctx context.Context, division *domain.Division) error {
	args := m.Called(ctx, division)
	return args.Error(0)
}

func (m *DivisionUseCaseMock) Update(ctx context.Context, division *domain.Division) error {
	args := m.Called(ctx, division)
	return args.Error(0)
}

func (m *DivisionUseCaseMock) Delete(ctx context.Context, id string, deletedBy string) error {
	args := m.Called(ctx, id, deletedBy)
	return args.Error(0)
}
