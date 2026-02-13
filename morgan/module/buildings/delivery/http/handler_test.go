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
	deliverhttp "github.com/siakup/morgan-be/morgan/module/buildings/delivery/http"
	"github.com/siakup/morgan-be/morgan/module/buildings/domain"
	"github.com/siakup/morgan-be/morgan/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupBuildingApp(useCase domain.UseCase) *fiber.App {
	// Mock auth middleware - passing nil as we bypass checks in tests by manually registering or using trusted middleware
	handler := deliverhttp.NewBuildingHandler(useCase, nil)

	app := fiber.New()

	// Mock middleware context variables
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.XInstitutionId, "inst-1")
		c.Locals(middleware.XUserIdKey, "admin-user")
		return c.Next()
	})

	app.Get("/buildings", handler.GetBuildings)
	app.Get("/buildings/:id", handler.GetBuildingByID)
	app.Post("/buildings", handler.CreateBuilding)
	app.Put("/buildings/:id", handler.UpdateBuilding)
	app.Delete("/buildings/:id", handler.DeleteBuilding)

	return app
}

func TestBuildingHandler_GetBuildings(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.BuildingUseCaseMock)
		app := setupBuildingApp(mockUseCase)

		buildings := []*domain.Building{{Id: "bldg1", Name: "Building A"}}
		count := int64(1)

		mockUseCase.On("GetAll", mock.Anything, mock.Anything).Return(buildings, count, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/buildings", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase := new(mocks.BuildingUseCaseMock)
		app := setupBuildingApp(mockUseCase)

		mockUseCase.On("GetAll", mock.Anything, mock.Anything).Return(([]*domain.Building)(nil), int64(0), errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodGet, "/buildings", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestBuildingHandler_GetBuildingByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.BuildingUseCaseMock)
		app := setupBuildingApp(mockUseCase)

		bldg := &domain.Building{Id: "bldg1", Name: "Building A"}
		mockUseCase.On("Get", mock.Anything, "bldg1").Return(bldg, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/buildings/bldg1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase := new(mocks.BuildingUseCaseMock)
		app := setupBuildingApp(mockUseCase)

		mockUseCase.On("Get", mock.Anything, "bldg1").Return((*domain.Building)(nil), errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodGet, "/buildings/bldg1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestBuildingHandler_CreateBuilding(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.BuildingUseCaseMock)
		app := setupBuildingApp(mockUseCase)

		reqBody := map[string]interface{}{
			"name": "Building B",
		}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(b *domain.Building) bool {
			return b.Name == "Building B"
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/buildings", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockUseCase := new(mocks.BuildingUseCaseMock)
		app := setupBuildingApp(mockUseCase)

		req := httptest.NewRequest(http.MethodPost, "/buildings", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UseCaseError", func(t *testing.T) {
		mockUseCase := new(mocks.BuildingUseCaseMock)
		app := setupBuildingApp(mockUseCase)

		reqBody := map[string]interface{}{"name": "Error Building"}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Create", mock.Anything, mock.Anything).Return(errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodPost, "/buildings", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestBuildingHandler_UpdateBuilding(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.BuildingUseCaseMock)
		app := setupBuildingApp(mockUseCase)

		reqBody := map[string]interface{}{
			"name": "Updated Building",
		}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Update", mock.Anything, mock.MatchedBy(func(b *domain.Building) bool {
			return b.Id == "bldg1" && b.Name == "Updated Building"
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/buildings/bldg1", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockUseCase := new(mocks.BuildingUseCaseMock)
		app := setupBuildingApp(mockUseCase)

		req := httptest.NewRequest(http.MethodPut, "/buildings/bldg1", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestBuildingHandler_DeleteBuilding(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.BuildingUseCaseMock)
		app := setupBuildingApp(mockUseCase)

		mockUseCase.On("Delete", mock.Anything, "bldg1", "admin-user").Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/buildings/bldg1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("UseCaseError", func(t *testing.T) {
		mockUseCase := new(mocks.BuildingUseCaseMock)
		app := setupBuildingApp(mockUseCase)

		mockUseCase.On("Delete", mock.Anything, "bldg1", "admin-user").Return(errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/buildings/bldg1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}
