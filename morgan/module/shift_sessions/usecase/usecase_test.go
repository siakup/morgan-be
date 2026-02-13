package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/domain"
	"github.com/siakup/morgan-be/morgan/module/shift_sessions/usecase"
	"github.com/siakup/morgan-be/morgan/tests/mocks"
)

func TestUseCase_ShiftSessions(t *testing.T) {
	mockRepo := new(mocks.ShiftSessionsRepositoryMock)
	uc := usecase.NewUseCase(mockRepo)

	t.Run("FindAll", func(t *testing.T) {
		ctx := context.Background()
		filter := domain.ShiftSessionFilter{}
		shiftSessions := []*domain.ShiftSession{{Id: "ss1"}}
		count := int64(1)

		mockRepo.On("FindAll", mock.Anything, filter).Return(shiftSessions, count, nil).Once()

		res, c, err := uc.FindAll(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, shiftSessions, res)
		assert.Equal(t, count, c)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Get", func(t *testing.T) {
		ctx := context.Background()
		id := "ss1"
		shiftSession := &domain.ShiftSession{Id: id}

		mockRepo.On("FindByID", mock.Anything, id).Return(shiftSession, nil).Once()

		res, err := uc.Get(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, shiftSession, res)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Create", func(t *testing.T) {
		ctx := context.Background()
		shiftSession := &domain.ShiftSession{
			Name:   "Morning Shift",
			Start:  "08:00",
			End:    "16:00",
			Status: true,
		}

		mockRepo.On("Store", mock.Anything, mock.MatchedBy(func(ss *domain.ShiftSession) bool {
			return ss.Name == shiftSession.Name && ss.Start == shiftSession.Start
		})).Return(nil).Once()

		err := uc.Create(ctx, shiftSession)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Update", func(t *testing.T) {
		ctx := context.Background()
		shiftSession := &domain.ShiftSession{
			Id:     "ss1",
			Name:   "Updated Morning Shift",
			Start:  "08:00",
			End:    "16:00",
			Status: true,
		}

		mockRepo.On("FindByID", mock.Anything, "ss1").Return(shiftSession, nil).Once()
		mockRepo.On("Update", mock.Anything, shiftSession).Return(nil).Once()

		err := uc.Update(ctx, shiftSession)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete", func(t *testing.T) {
		ctx := context.Background()
		id := "ss1"

		deletedBy := "user-1"

		mockRepo.On("Delete", mock.Anything, id, deletedBy).Return(nil).Once()

		err := uc.Delete(ctx, id, deletedBy)
		assert.NoError(t, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Get_NotFound", func(t *testing.T) {
		ctx := context.Background()
		mockRepo.On("FindByID", mock.Anything, "ss1").Return((*domain.ShiftSession)(nil), pgx.ErrNoRows).Once()

		res, err := uc.Get(ctx, "ss1")
		assert.Error(t, err)
		assert.Nil(t, res)
		assert.Contains(t, err.Error(), "not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Update_NotFound", func(t *testing.T) {
		ctx := context.Background()
		shiftSession := &domain.ShiftSession{Id: "ss1"}

		mockRepo.On("FindByID", mock.Anything, "ss1").Return((*domain.ShiftSession)(nil), pgx.ErrNoRows).Once()

		err := uc.Update(ctx, shiftSession)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Create_Error", func(t *testing.T) {
		ctx := context.Background()
		shiftSession := &domain.ShiftSession{Name: "Test"}
		mockRepo.On("Store", mock.Anything, mock.Anything).Return(errors.New("store failed")).Once()

		err := uc.Create(ctx, shiftSession)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("FindAll_Error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo.On("FindAll", mock.Anything, mock.Anything).Return(([]*domain.ShiftSession)(nil), int64(0), errors.New("query failed")).Once()

		_, _, err := uc.FindAll(ctx, domain.ShiftSessionFilter{})
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Delete_Error", func(t *testing.T) {
		ctx := context.Background()
		mockRepo.On("Delete", mock.Anything, "ss1", "user-1").Return(errors.New("delete failed")).Once()

		err := uc.Delete(ctx, "ss1", "user-1")
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})
}
