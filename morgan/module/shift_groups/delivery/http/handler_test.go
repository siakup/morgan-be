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
	deliverhttp "github.com/siakup/morgan-be/morgan/module/shift_groups/delivery/http"
	"github.com/siakup/morgan-be/morgan/module/shift_groups/domain"
	"github.com/siakup/morgan-be/morgan/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupShiftGroupApp(useCase domain.UseCase) *fiber.App {
	// Mock auth middleware
	// We pass nil for auth middleware because we manually register routes below,
	// bypassing the RegisterRoutes method which would require a valid auth middleware.

	handler := deliverhttp.NewShiftGroupHandler(useCase, nil)

	app := fiber.New()

	// Mock middleware context variables if needed
	app.Use(func(c *fiber.Ctx) error {
		c.Locals(middleware.XInstitutionId, "inst-1")
		c.Locals(middleware.XUserIdKey, "admin-user")
		return c.Next()
	})

	app.Get("/shift-groups", handler.GetShiftGroups)
	app.Get("/shift-groups/:id", handler.GetShiftGroupByID)
	app.Post("/shift-groups", handler.CreateShiftGroup)
	app.Put("/shift-groups/:id", handler.UpdateShiftGroup)
	app.Delete("/shift-groups/:id", handler.DeleteShiftGroup)

	return app
}

func TestShiftGroupHandler_GetShiftGroups(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.ShiftGroupsUseCaseMock)
		app := setupShiftGroupApp(mockUseCase)

		shiftGroups := []*domain.ShiftGroup{{Id: "sg1", Name: "Morning Shift"}}
		count := int64(1)

		mockUseCase.On("FindAll", mock.Anything, mock.Anything).Return(shiftGroups, count, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/shift-groups", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase := new(mocks.ShiftGroupsUseCaseMock)
		app := setupShiftGroupApp(mockUseCase)

		mockUseCase.On("FindAll", mock.Anything, mock.Anything).Return(([]*domain.ShiftGroup)(nil), int64(0), errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodGet, "/shift-groups", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestShiftGroupHandler_GetShiftGroupByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.ShiftGroupsUseCaseMock)
		app := setupShiftGroupApp(mockUseCase)

		sg := &domain.ShiftGroup{Id: "sg1", Name: "Morning Shift"}
		mockUseCase.On("Get", mock.Anything, "sg1").Return(sg, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/shift-groups/sg1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockUseCase := new(mocks.ShiftGroupsUseCaseMock)
		app := setupShiftGroupApp(mockUseCase)

		mockUseCase.On("Get", mock.Anything, "sg1").Return((*domain.ShiftGroup)(nil), errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodGet, "/shift-groups/sg1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestShiftGroupHandler_CreateShiftGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.ShiftGroupsUseCaseMock)
		app := setupShiftGroupApp(mockUseCase)

		reqBody := map[string]interface{}{
			"name": "Evening Shift",
		}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Create", mock.Anything, mock.MatchedBy(func(sg *domain.ShiftGroup) bool {
			return sg.Name == "Evening Shift"
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/shift-groups", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockUseCase := new(mocks.ShiftGroupsUseCaseMock)
		app := setupShiftGroupApp(mockUseCase)

		req := httptest.NewRequest(http.MethodPost, "/shift-groups", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UseCaseError", func(t *testing.T) {
		mockUseCase := new(mocks.ShiftGroupsUseCaseMock)
		app := setupShiftGroupApp(mockUseCase)

		reqBody := map[string]interface{}{"name": "Error Shift"}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Create", mock.Anything, mock.Anything).Return(errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodPost, "/shift-groups", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}

func TestShiftGroupHandler_UpdateShiftGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.ShiftGroupsUseCaseMock)
		app := setupShiftGroupApp(mockUseCase)

		reqBody := map[string]interface{}{
			"name": "Updated Shift",
		}
		reqBytes, _ := json.Marshal(reqBody)

		mockUseCase.On("Update", mock.Anything, mock.MatchedBy(func(sg *domain.ShiftGroup) bool {
			return sg.Id == "sg1" && sg.Name == "Updated Shift"
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/shift-groups/sg1", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("BadRequest", func(t *testing.T) {
		mockUseCase := new(mocks.ShiftGroupsUseCaseMock)
		app := setupShiftGroupApp(mockUseCase)

		req := httptest.NewRequest(http.MethodPut, "/shift-groups/sg1", bytes.NewReader([]byte("invalid")))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestShiftGroupHandler_DeleteShiftGroup(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockUseCase := new(mocks.ShiftGroupsUseCaseMock)
		app := setupShiftGroupApp(mockUseCase)

		mockUseCase.On("Delete", mock.Anything, "sg1", "admin-user").Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/shift-groups/sg1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})

	t.Run("UseCaseError", func(t *testing.T) {
		mockUseCase := new(mocks.ShiftGroupsUseCaseMock)
		app := setupShiftGroupApp(mockUseCase)

		mockUseCase.On("Delete", mock.Anything, "sg1", "admin-user").Return(errors.New("fail")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/shift-groups/sg1", nil)
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockUseCase.AssertExpectations(t)
	})
}
