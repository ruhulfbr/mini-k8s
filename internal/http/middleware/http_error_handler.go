package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func EchoHTTPErrorHandler(err error, c echo.Context) {
	var ve validator.ValidationErrors

	if errors.As(err, &ve) {
		errorsArr := make([]FieldError, 0, len(ve))

		for _, fe := range ve {
			errorsArr = append(errorsArr, FieldError{
				Field:   fe.Field(),
				Message: validationMessage(fe),
			})
		}

		_ = c.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
			"message":   "Validation failed",
			"apperrors": errorsArr,
		})
		return
	}

	// App errors
	var appErr *appErrors.AppError
	if errors.As(err, &appErr) {
		_ = c.JSON(appErr.Code, map[string]interface{}{
			"message": appErr.Message,
		})
		return
	}

	// --------------------------------------------------
	// 2. Handle Echo HTTP apperrors
	// --------------------------------------------------
	var he *echo.HTTPError
	if errors.As(err, &he) {
		_ = c.JSON(he.Code, map[string]interface{}{
			"message": he.Message,
		})
		return
	}

	// --------------------------------------------------
	// 3. Fallback: Bad Request
	// --------------------------------------------------
	_ = c.JSON(http.StatusBadRequest, map[string]interface{}{
		"message": appErrors.SomethingWentWrong,
	})
}

func validationMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email address"
	case "min":
		return "Must be at least " + fe.Param()
	case "max":
		return "Must be at most " + fe.Param()
	case "gte":
		return "Must be greater than or equal to " + fe.Param()
	case "lte":
		return "Must be less than or equal to " + fe.Param()
	case "url":
		return "The field value must be valid url"
	case "unique":
		return "Duplicate values are not allowed"
	case "oneof":
		return "Must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
	default:
		return "Invalid value"
	}
}
