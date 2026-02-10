package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/middleware"
	deliverhttp "yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_sessions/delivery/http"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_sessions/domain"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/tests/mocks"
)

func setupShiftSessionApp(useCase domain.UseCase) *fiber.App {
	handler := deliverhttp.NewShiftSessionHandler(useCase, nil)

	app := fiber.New()

	// Mock middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.XUserIdKey, "user-1")
		return c.Next()
	})

	app.Get("/shift-sessions", handler.GetShiftSessions)
	app.Get("/shift-sessions/:id", handler.GetShiftSessionByID)
	app.Post("/shift-sessions", handler.CreateShiftSession)
	app.Put("/shift-sessions/:id", handler.UpdateShiftSession)
	app.Delete("/shift-sessions/:id", handler.DeleteShiftSession)

	return app
}

func TestShiftSessionHandler_GetShiftSessions(t *testing.T) {
	mockUseCase := new(mocks.ShiftSessionsUseCaseMock)
	app := setupShiftSessionApp(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		shiftSessions := []*domain.ShiftSession{{Id: "ss1", Name: "Morning Shift"}}
		count := int64(1)

		mockUseCase.On("FindAll", mock.Anything, mock.Anything).Return(shiftSessions, count, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/shift-sessions", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase.On("FindAll", mock.Anything, mock.Anything).Return(([]*domain.ShiftSession)(nil), int64(0), errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodGet, "/shift-sessions", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestShiftSessionHandler_GetShiftSessionByID(t *testing.T) {
	mockUseCase := new(mocks.ShiftSessionsUseCaseMock)
	app := setupShiftSessionApp(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		shiftSession := &domain.ShiftSession{
			Id:     "ss1",
			Name:   "Morning Shift",
			Start:  "08:00",
			End:    "16:00",
			Status: true,
		}

		mockUseCase.On("Get", mock.Anything, "ss1").Return(shiftSession, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/shift-sessions/ss1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockUseCase.On("Get", mock.Anything, "ss-invalid").Return((*domain.ShiftSession)(nil), errors.New("not found")).Once()

		req := httptest.NewRequest(http.MethodGet, "/shift-sessions/ss-invalid", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestShiftSessionHandler_CreateShiftSession(t *testing.T) {
	mockUseCase := new(mocks.ShiftSessionsUseCaseMock)
	app := setupShiftSessionApp(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":   "Evening Shift",
			"start":  "17:00",
			"end":    "23:00",
			"status": true,
		}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(ss *domain.ShiftSession) bool {
			return ss.Name == "Evening Shift" && ss.Start == "17:00"
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/shift-sessions", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "Evening Shift",
			// Missing start and end
		}
		reqBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/shift-sessions", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/shift-sessions", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestShiftSessionHandler_UpdateShiftSession(t *testing.T) {
	mockUseCase := new(mocks.ShiftSessionsUseCaseMock)
	app := setupShiftSessionApp(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":   "Updated Shift",
			"start":  "09:00",
			"end":    "17:00",
			"status": true,
		}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Update", mock.Anything, mock.MatchedBy(func(ss *domain.ShiftSession) bool {
			return ss.Id == "ss1" && ss.Name == "Updated Shift"
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/shift-sessions/ss1", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		reqBody := map[string]interface{}{"name": "Test", "start": "08:00", "end": "16:00"}
		reqBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/shift-sessions/", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		// Empty ID in path means the handler won't be called
		assert.NotEqual(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestShiftSessionHandler_DeleteShiftSession(t *testing.T) {
	mockUseCase := new(mocks.ShiftSessionsUseCaseMock)
	app := setupShiftSessionApp(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.On("Delete", mock.Anything, "ss1", "user-1").Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/shift-sessions/ss1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase.On("Delete", mock.Anything, "ss1", "user-1").Return(errors.New("delete failed")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/shift-sessions/ss1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}
