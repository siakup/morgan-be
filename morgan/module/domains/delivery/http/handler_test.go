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
	"github.com/siakup/morgan-be/libraries/middleware"
	deliverhttp "github.com/siakup/morgan-be/morgan/module/domains/delivery/http"
	"github.com/siakup/morgan-be/morgan/module/domains/domain"
	"github.com/siakup/morgan-be/morgan/tests/mocks"
)

func setupDomainApp(useCase domain.UseCase) *fiber.App {
	handler := deliverhttp.NewDomaiHandler(useCase, nil)

	app := fiber.New()

	// Mock middleware
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.XUserIdKey, "user-1")
		return c.Next()
	})

	app.Get("/domains", handler.GetDomains)
	app.Get("/domains/:id", handler.GetDomainByID)
	app.Post("/domains", handler.CreateDomain)
	app.Put("/domains/:id", handler.UpdateDomain)
	app.Delete("/domains/:id", handler.DeleteDomain)

	return app
}

func TestDomainHandler_GetDomains(t *testing.T) {
	mockUseCase := new(mocks.DomainsUseCaseMock)
	app := setupDomainApp(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		domains := []*domain.Domain{{Id: "d1", Name: "Domain 1"}}
		count := int64(1)

		mockUseCase.On("FindAll", mock.Anything, mock.Anything).Return(domains, count, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/domains", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase.On("FindAll", mock.Anything, mock.Anything).Return(([]*domain.Domain)(nil), int64(0), errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodGet, "/domains", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestDomainHandler_GetDomainByID(t *testing.T) {
	mockUseCase := new(mocks.DomainsUseCaseMock)
	app := setupDomainApp(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		domain := &domain.Domain{
			Id:     "d1",
			Name:   "Domain 1",
			Status: true,
		}

		mockUseCase.On("Get", mock.Anything, "d1").Return(domain, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/domains/d1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockUseCase.On("Get", mock.Anything, "d-invalid").Return((*domain.Domain)(nil), errors.New("not found")).Once()

		req := httptest.NewRequest(http.MethodGet, "/domains/d-invalid", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestDomainHandler_CreateDomain(t *testing.T) {
	mockUseCase := new(mocks.DomainsUseCaseMock)
	app := setupDomainApp(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "New Domain",
		}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(d *domain.Domain) bool {
			return d.Name == "New Domain"
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/domains", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		reqBody := map[string]interface{}{
			// Missing name
		}
		reqBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/domains", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/domains", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDomainHandler_UpdateDomain(t *testing.T) {
	mockUseCase := new(mocks.DomainsUseCaseMock)
	app := setupDomainApp(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name":   "Updated Domain",
			"status": true,
		}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Update", mock.Anything, mock.MatchedBy(func(d *domain.Domain) bool {
			return d.Id == "d1" && d.Name == "Updated Domain"
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/domains/d1", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("MissingRequiredField", func(t *testing.T) {
		reqBody := map[string]interface{}{
			"name": "",
		}
		reqBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/domains/d1", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestDomainHandler_DeleteDomain(t *testing.T) {
	mockUseCase := new(mocks.DomainsUseCaseMock)
	app := setupDomainApp(mockUseCase)

	t.Run("Success", func(t *testing.T) {
		mockUseCase.On("Delete", mock.Anything, "d1", "user-1").Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/domains/d1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase.On("Delete", mock.Anything, "d1", "user-1").Return(errors.New("delete failed")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/domains/d1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}
