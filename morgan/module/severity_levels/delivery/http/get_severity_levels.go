package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/libraries/types"
	"github.com/siakup/morgan-be/morgan/module/severity_levels/domain"
)

func (h *SeverityLevelHandler) GetSeverityLevels(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	size := c.QueryInt("size", 10)

	filter := domain.SeverityLevelFilter{
		Pagination: types.Pagination{
			Page: page,
			Size: size,
		},
		Search: c.Query("search"),
	}

	severityLevels, total, err := h.useCase.FindAll(c.Context(), filter)
	if err != nil {
		return h.handleError(c, err)
	}

	meta := &responses.Meta{
		Page:       page,
		Size:       size,
		Total:      total,
		TotalPages: (int(total) + size - 1) / size,
	}

	return c.JSON(responses.SuccessWithMeta(severityLevels, "Severity Levels retrieved", meta))
}
