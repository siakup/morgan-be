package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/siakup/morgan-be/libraries/errors"
)

var validate = validator.New()

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func ValidateStruct(payload interface{}) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors
}

func BindAndValidate(c *fiber.Ctx, payload interface{}) error {
	if err := c.BodyParser(payload); err != nil {
		return errors.BadRequest("Invalid request body")
	}

	if validationErrors := ValidateStruct(payload); len(validationErrors) > 0 {
		var errMsgs []string
		for _, err := range validationErrors {
			errMsgs = append(errMsgs, fmt.Sprintf("[%s]: %s", err.Field, err.Tag))
		}
		return errors.BadRequest(strings.Join(errMsgs, "; "))
	}

	return nil
}
