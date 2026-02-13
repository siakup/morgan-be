package http_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	deliverhttp "github.com/siakup/morgan-be/morgan/module/severity_levels/delivery/http"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
	"github.com/siakup/morgan-be/morgan/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupSeverityLevelApp(useCase domain.UseCase) *fiber.App {
	handler := deliverhttp.NewSeverityLevelHandler(useCase, nil)

	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.XUserIdKey, "admin-user")
		return c.Next()
	})

	app.Get("/severity-levels", handler.GetSeverityLevels)
	app.Get("/severity-levels/:id", handler.GetSeverityLevelByID)
	app.Post("/severity-levels", handler.CreateSeverityLevel)
	app.Put("/severity-levels/:id", handler.UpdateSeverityLevel)
	app.Delete("/severity-levels/:id", handler.DeleteSeverityLevel)

	return app
}

func TestSeverityLevelHandler_GetSeverityLevels(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.SeverityLevelsUseCaseMock)
		app := setupSeverityLevelApp(mockUseCase)

		severityLevels := []*domain.SeverityLevel{{Id: "sl1", Name: "Low"}}
		count := int64(1)

		mockUseCase.On("FindAll", mock.Anything, mock.Anything).Return(severityLevels, count, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/severity-levels", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase := new(mocks.SeverityLevelsUseCaseMock)
		app := setupSeverityLevelApp(mockUseCase)

		mockUseCase.On("FindAll", mock.Anything, mock.Anything).Return(([]*domain.SeverityLevel)(nil), int64(0), errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodGet, "/severity-levels", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestSeverityLevelHandler_GetSeverityLevelByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.SeverityLevelsUseCaseMock)
		app := setupSeverityLevelApp(mockUseCase)

		sl := &domain.SeverityLevel{Id: "sl1", Name: "Low"}
		mockUseCase.On("Get", mock.Anything, "sl1").Return(sl, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/severity-levels/sl1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase := new(mocks.SeverityLevelsUseCaseMock)
		app := setupSeverityLevelApp(mockUseCase)

		mockUseCase.On("Get", mock.Anything, "sl1").Return((*domain.SeverityLevel)(nil), errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodGet, "/severity-levels/sl1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestSeverityLevelHandler_CreateSeverityLevel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.SeverityLevelsUseCaseMock)
		app := setupSeverityLevelApp(mockUseCase)

		reqBody := map[string]interface{}{
			"name": "High",
		}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(sl *domain.SeverityLevel) bool {
			return sl.Name == "High"
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/severity-levels", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockUseCase := new(mocks.SeverityLevelsUseCaseMock)
		app := setupSeverityLevelApp(mockUseCase)

		req := httptest.NewRequest(http.MethodPost, "/severity-levels", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UseCaseError", func(t *testing.T) {
		mockUseCase := new(mocks.SeverityLevelsUseCaseMock)
		app := setupSeverityLevelApp(mockUseCase)

		reqBody := map[string]interface{}{"name": "Error"}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Create", mock.Anything, mock.Anything).Return(errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodPost, "/severity-levels", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestSeverityLevelHandler_UpdateSeverityLevel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.SeverityLevelsUseCaseMock)
		app := setupSeverityLevelApp(mockUseCase)

		reqBody := map[string]interface{}{
			"name": "Updated Low",
		}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Update", mock.Anything, mock.MatchedBy(func(sl *domain.SeverityLevel) bool {
			return sl.Id == "sl1" && sl.Name == "Updated Low"
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/severity-levels/sl1", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockUseCase := new(mocks.SeverityLevelsUseCaseMock)
		app := setupSeverityLevelApp(mockUseCase)

		req := httptest.NewRequest(http.MethodPut, "/severity-levels/sl1", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestSeverityLevelHandler_DeleteSeverityLevel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.SeverityLevelsUseCaseMock)
		app := setupSeverityLevelApp(mockUseCase)

		mockUseCase.On("Delete", mock.Anything, "sl1", "admin-user").Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/severity-levels/sl1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("UseCaseError", func(t *testing.T) {
		mockUseCase := new(mocks.SeverityLevelsUseCaseMock)
		app := setupSeverityLevelApp(mockUseCase)

		mockUseCase.On("Delete", mock.Anything, "sl1", "admin-user").Return(errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/severity-levels/sl1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}
