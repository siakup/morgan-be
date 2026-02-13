package http

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/middleware"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/libraries/types"
	"github.com/siakup/morgan-be/morgan/module/users/domain"
)

type (
	GetUserResponse struct {
		Id              string         `json:"id"`
		InstitutionId   string         `json:"institution_id"`
		ExternalSubject string         `json:"external_subject"`
		Status          string         `json:"status"`
		Metadata        map[string]any `json:"metadata"`
	}
)

// GetUsers handles GET /users
func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	ctx := c.UserContext()

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	filter := domain.UserFilter{
		Pagination: types.Pagination{
			Page: page,
			Size: pageSize,
		},
		InstitutionId: c.Locals(middleware.XInstitutionId).(string),
		Status:        c.Query("status"),
		Search:        c.Query("search"),
	}

	users, total, err := h.useCase.FindAll(ctx, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	result := make([]GetUserResponse, len(users))
	for i, u := range users {
		result[i] = GetUserResponse{
			Id:              u.Id,
			InstitutionId:   u.InstitutionId,
			ExternalSubject: u.ExternalSubject,
			Status:          u.Status,
			Metadata:        u.Metadata,
		}
	}

	meta := &responses.Meta{
		Page:       page,
		Size:       pageSize,
		Total:      total,
		TotalPages: (int(total) + pageSize - 1) / pageSize,
	}

	return c.Status(http.StatusOK).JSON(responses.SuccessWithMeta(result, "Users retrieved", meta))
}
