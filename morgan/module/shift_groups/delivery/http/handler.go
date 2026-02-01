package http

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/errors"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/middleware"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/types"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/shift_groups/domain"
)

// ShiftGroupHandler handles HTTP requests for shift groups module.
type ShiftGroupHandler struct {
	useCase domain.UseCase
	auth    *middleware.AuthorizationMiddleware
}

// NewShiftGroupHandler creates a new ShiftGroupHandler.
func NewShiftGroupHandler(useCase domain.UseCase, auth *middleware.AuthorizationMiddleware) *ShiftGroupHandler {
	return &ShiftGroupHandler{
		useCase: useCase,
		auth:    auth,
	}
}

// RegisterRoutes registers the routes for the shift groups module.
func (h *ShiftGroupHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/shift-groups", middleware.TraceMiddleware)

	// Adjust permissions as needed
	group.Get("/", h.auth.Authenticate("shift_groups.view"), h.FindAll)
	group.Get("/:id", h.auth.Authenticate("shift_groups.view"), h.FindByID)
	group.Post("/", h.auth.Authenticate("shift_groups.edit"), h.Create)
	group.Put("/:id", h.auth.Authenticate("shift_groups.edit"), h.Update)
	group.Delete("/:id", h.auth.Authenticate("shift_groups.edit"), h.Delete)
}

func (h *ShiftGroupHandler) FindAll(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	size := c.QueryInt("size", 10)

	filter := domain.ShiftGroupFilter{
		Pagination: types.Pagination{
			Page: page,
			Size: size,
		},
		Search: c.Query("search"),
	}

	shiftGroups, total, err := h.useCase.FindAll(c.Context(), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	meta := &responses.Meta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: (int(total) + size - 1) / size,
	}

	return c.JSON(responses.SuccessWithMeta(shiftGroups, "Shift Groups retrieved", meta))
}

func (h *ShiftGroupHandler) FindByID(c *fiber.Ctx) error {
	id := c.Params("id")
	shiftGroup, err := h.useCase.FindByID(c.Context(), id)
	if err != nil {
		return h.handleError(c, err)
	}
	if shiftGroup == nil {
		return c.Status(http.StatusNotFound).JSON(responses.Fail("NOT_FOUND", "Shift Group not found"))
	}

	return c.JSON(responses.Success(shiftGroup, "Shift Group retrieved"))
}

type (
	CreateShiftGroupRequest struct {
		Name   string `json:"name" validate:"required"`
		Status bool   `json:"status"`
	}
	UpdateShiftGroupRequest struct {
		Name   string `json:"name" validate:"required"`
		Status bool   `json:"status"`
	}
)

func (h *ShiftGroupHandler) Create(c *fiber.Ctx) error {
	var req CreateShiftGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Fail("BAD_REQUEST", "Invalid request body"))
	}

	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	shiftGroup := domain.ShiftGroup{
		Name:      req.Name,
		Status:    req.Status,
		CreatedBy: &userId,
		UpdatedBy: &userId,
	}

	err := h.useCase.Create(c.Context(), &shiftGroup)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusCreated).JSON(responses.Success(shiftGroup, "Shift Group created"))
}

func (h *ShiftGroupHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req UpdateShiftGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Fail("BAD_REQUEST", "Invalid request body"))
	}

	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	shiftGroup := domain.ShiftGroup{
		Id:        id,
		Name:      req.Name,
		Status:    req.Status,
		UpdatedBy: &userId,
	}

	err := h.useCase.Update(c.Context(), &shiftGroup)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(responses.Success(shiftGroup, "Shift Group updated"))
}

func (h *ShiftGroupHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	userId, _ := c.Locals(middleware.XUserIdKey).(string)

	err := h.useCase.Delete(c.Context(), id, userId)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.Status(http.StatusOK).JSON(responses.Success[any](nil, "Shift Group deleted"))
}

// handleError handles errors by mapping them to standardized responses.
func (h *ShiftGroupHandler) handleError(c *fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return c.Status(appErr.Code).JSON(responses.Fail(string(appErr.Type), appErr.Message))
	}

	return c.Status(http.StatusInternalServerError).JSON(responses.Fail("SYSTEM_ERROR", err.Error()))
}
