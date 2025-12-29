package middleware

import (
	"errors"
	"net/http"

	"github.com/go-ozzo/ozzo-validation/v4"
	"github.com/labstack/echo/v4"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func EchoHTTPErrorHandler(err error, c echo.Context) {
	// Handle validation errors
	var ve validation.Errors
	if errors.As(err, &ve) {
		var errorsArr []FieldError
		for field, e := range ve {
			errorsArr = append(errorsArr, FieldError{
				Field:   field,
				Message: e.Error(),
			})
		}

		err := c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"code":    http.StatusUnprocessableEntity,
			"message": "Validation Error",
			"errors":  errorsArr,
		})
		if err != nil {
			return
		}
		return
	}

	var he *echo.HTTPError
	if errors.As(err, &he) {
		err := c.JSON(he.Code, map[string]interface{}{
			"code":  he.Code,
			"error": he.Message,
		})
		if err != nil {
			return
		}
		return
	}

	err = c.JSON(http.StatusInternalServerError, map[string]interface{}{
		"code":  500,
		"error": err.Error(),
	})
	if err != nil {
		return
	}
}
