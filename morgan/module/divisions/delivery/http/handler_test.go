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
	deliverhttp "github.com/siakup/morgan-be/morgan/module/divisions/delivery/http"
	"github.com/siakup/morgan-be/morgan/module/divisions/domain"
	"github.com/siakup/morgan-be/morgan/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupDivisionApp(useCase domain.UseCase) *fiber.App {
	// Mock auth middleware - passing nil as we bypass checks in tests by manually registering or using trusted middleware
	handler := deliverhttp.NewDivisionHandler(useCase, nil)

	app := fiber.New()

	// Mock middleware context variables
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.XInstitutionId, "inst-1")
		c.Locals(middleware.XUserIdKey, "admin-user")
		return c.Next()
	})

	app.Get("/divisions", handler.GetDivisions)
	app.Get("/divisions/:id", handler.GetDivisionByID)
	app.Post("/divisions", handler.CreateDivision)
	app.Put("/divisions/:id", handler.UpdateDivision)
	app.Delete("/divisions/:id", handler.DeleteDivision)

	return app
}

func TestDivisionHandler_GetDivisions(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.DivisionUseCaseMock)
		app := setupDivisionApp(mockUseCase)

		divisions := []*domain.Division{{Id: "div1", Name: "IT Division"}}
		count := int64(1)

		mockUseCase.On("GetAll", mock.Anything, mock.Anything).Return(divisions, count, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/divisions", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase := new(mocks.DivisionUseCaseMock)
		app := setupDivisionApp(mockUseCase)

		mockUseCase.On("GetAll", mock.Anything, mock.Anything).Return(([]*domain.Division)(nil), int64(0), errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodGet, "/divisions", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestDivisionHandler_GetDivisionByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.DivisionUseCaseMock)
		app := setupDivisionApp(mockUseCase)

		div := &domain.Division{Id: "div1", Name: "IT Division"}
		mockUseCase.On("Get", mock.Anything, "div1").Return(div, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/divisions/div1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase := new(mocks.DivisionUseCaseMock)
		app := setupDivisionApp(mockUseCase)

		mockUseCase.On("Get", mock.Anything, "div1").Return((*domain.Division)(nil), errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodGet, "/divisions/div1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestDivisionHandler_CreateDivision(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.DivisionUseCaseMock)
		app := setupDivisionApp(mockUseCase)

		reqBody := map[string]interface{}{
			"name": "HR Division",
		}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(d *domain.Division) bool {
			return d.Name == "HR Division"
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/divisions", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockUseCase := new(mocks.DivisionUseCaseMock)
		app := setupDivisionApp(mockUseCase)

		req := httptest.NewRequest(http.MethodPost, "/divisions", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UseCaseError", func(t *testing.T) {
		mockUseCase := new(mocks.DivisionUseCaseMock)
		app := setupDivisionApp(mockUseCase)

		reqBody := map[string]interface{}{"name": "Error Division"}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Create", mock.Anything, mock.Anything).Return(errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodPost, "/divisions", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestDivisionHandler_UpdateDivision(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.DivisionUseCaseMock)
		app := setupDivisionApp(mockUseCase)

		reqBody := map[string]interface{}{
			"name": "Updated Division",
		}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Update", mock.Anything, mock.MatchedBy(func(d *domain.Division) bool {
			return d.Id == "div1" && d.Name == "Updated Division"
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/divisions/div1", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockUseCase := new(mocks.DivisionUseCaseMock)
		app := setupDivisionApp(mockUseCase)

		req := httptest.NewRequest(http.MethodPut, "/divisions/div1", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDivisionHandler_DeleteDivision(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.DivisionUseCaseMock)
		app := setupDivisionApp(mockUseCase)

		mockUseCase.On("Delete", mock.Anything, "div1", "admin-user").Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/divisions/div1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("UseCaseError", func(t *testing.T) {
		mockUseCase := new(mocks.DivisionUseCaseMock)
		app := setupDivisionApp(mockUseCase)

		mockUseCase.On("Delete", mock.Anything, "div1", "admin-user").Return(errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/divisions/div1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}
