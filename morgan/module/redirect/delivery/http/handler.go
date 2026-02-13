package http

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/responses"
	"github.com/siakup/morgan-be/morgan/module/redirect/domain"
)

type RedirectHandler struct {
	useCase domain.RedirectUseCase
}

func NewHandler(useCase domain.RedirectUseCase) *RedirectHandler {
	return &RedirectHandler{
		useCase: useCase,
	}
}

func (h *RedirectHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/redirect/:institution_id", h.Redirect)
}

func (h *RedirectHandler) Redirect(c *fiber.Ctx) error {
	ctx := c.UserContext()
	institutionId := c.Params("institution_id")
	token := c.Query("token")

	if institutionId == "" || token == "" {
		return c.Status(http.StatusBadRequest).JSON(responses.Fail(strconv.Itoa(http.StatusBadRequest), "institution_id and token are required"))
	}

	redirectUrl, sessionId, err := h.useCase.Redirect(ctx, institutionId, token)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(http.StatusNotFound).JSON(responses.Fail(strconv.Itoa(http.StatusNotFound), err.Error()))
		}
		if strings.Contains(err.Error(), "invalid token") {
			return c.Status(http.StatusForbidden).JSON(responses.Fail(strconv.Itoa(http.StatusForbidden), err.Error()))
		}

		return c.Status(http.StatusInternalServerError).JSON(responses.Fail(strconv.Itoa(http.StatusInternalServerError), err.Error()))
	}

	// Set Cookie
	c.Cookie(&fiber.Cookie{
		Name:     "session_id",
		Value:    sessionId,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	})

	return c.Redirect(redirectUrl)
}
