package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/usecase"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/tests/mocks"
)

func TestUseCase_Domains(t *testing.T) {
	mockRepo := new(mocks.DomainsRepositoryMock)
	uc := usecase.NewUseCase(mockRepo)

	t.Run("FindAll", func(t *testing.T) {
		ctx := context.Background()
		filter := domain.DomainFilter{}
		domains := []*domain.Domain{{Id: "d1", Name: "Domain 1"}}
		count := int64(1)

		mockRepo.On("FindAll", mock.Anything, filter).Return(domains, count, nil).Once()

		res, c, err := uc.FindAll(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, domains, res)
		assert.Equal(t, count, c)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Get", func(t *testing.T) {
		ctx := context.Background()
		id := "d1"
		domain := &domain.Domain{Id: id, Name: "Domain 1"}

		mockRepo.On("FindByID", mock.Anything, id).Return(domain, nil).Once()

		res, err := uc.Get(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, domain, res)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Create", func(t *testing.T) {
		ctx := context.Background()
		newDomain := &domain.Domain{
			Name:   "New Domain",
			Status: true,
		}

		mockRepo.On("Store", mock.Anything, mock.MatchedBy(func(d *domain.Domain) bool {
			return d.Name == newDomain.Name && d.Status == newDomain.Status
		})).Return(nil).Once()

		err := uc.Create(ctx, newDomain)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Update", func(t *testing.T) {
		ctx := context.Background()
		domain := &domain.Domain{
			Id:     "d1",
			Name:   "Updated Domain",
			Status: true,
		}

		mockRepo.On("FindByID", mock.Anything, "d1").Return(domain, nil).Once()
		mockRepo.On("Update", mock.Anything, domain).Return(nil).Once()

		err := uc.Update(ctx, domain)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		ctx := context.Background()
		id := "d1"
		deletedBy := "user-1"

		mockRepo.On("Delete", mock.Anything, id, deletedBy).Return(nil).Once()

		err := uc.Delete(ctx, id, deletedBy)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Get_NotFound", func(t *testing.T) {
		ctx := context.Background()
		mockRepo.On("FindByID", mock.Anything, "d1").Return((*domain.Domain)(nil), pgx.ErrNoRows).Once()

		res, err := uc.Get(ctx, "d1")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		ctx := context.Background()
		updateDomain := &domain.Domain{Id: "d1"}

		mockRepo.On("FindByID", mock.Anything, "d1").Return((*domain.Domain)(nil), pgx.ErrNoRows).Once()

		err := uc.Update(ctx, updateDomain)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create_Error", func(t *testing.T) {
		ctx := context.Background()
		domain := &domain.Domain{Name: "Test"}
		mockRepo.On("Store", mock.Anything, mock.Anything).Return(errors.New("store failed")).Once()

		err := uc.Create(ctx, domain)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FindAll_Error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo.On("FindAll", mock.Anything, mock.Anything).Return(([]*domain.Domain)(nil), int64(0), errors.New("query failed")).Once()

		_, _, err := uc.FindAll(ctx, domain.DomainFilter{})
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete_Error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo.On("Delete", mock.Anything, "d1", "user-1").Return(errors.New("delete failed")).Once()

		err := uc.Delete(ctx, "d1", "user-1")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
