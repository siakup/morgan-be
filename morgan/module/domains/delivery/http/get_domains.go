package http

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/responses"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/lib/types"
	"yuhuu.universitaspertamina.ac.id/siak/siakup/backend/morgan/module/domains/domain"
)

type (
	GetDomainsResponse struct {
		Id     string `json:"id"`
		Name   string `json:"name"`
		Status bool   `json:"status"`
	}
)

// GetDomains handles GET /domains
func (h *DomainHandler) GetDomains(c *fiber.Ctx) error {
	ctx := c.UserContext()

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	filter := domain.DomainFilter{
		Pagination: types.Pagination{
			Page: page,
			Size: pageSize,
		},
		Search: c.Query("search"),
	}

	domains, total, err := h.useCase.FindAll(ctx, filter)
	if err != nil {
		return h.handleError(c, err)
	}

	result := make([]GetDomainsResponse, len(domains))
	for i, d := range domains {
		result[i] = GetDomainsResponse{
			Id:     d.Id,
			Name:   d.Name,
			Status: d.Status,
		}
	}

	meta := &responses.Meta{
		Page:       page,
		Size:       pageSize,
		Total:      total,
		TotalPages: (int(total) + pageSize - 1) / pageSize,
	}

	return c.Status(http.StatusOK).JSON(responses.SuccessWithMeta(result, "Domains retrieved", meta))
}
