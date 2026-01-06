package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"

	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Message   string       `json:"message"`
	AppErrors []FieldError `json:"apperrors,omitempty"`
}

func NewEchoHTTPErrorHandler() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {

		switch {
		case handleValidationError(c, err):
			return
		case handleAppError(c, err):
			return
		case handleHTTPError(c, err):
			return
		default:
			logError(c, "unhandled error", err)
			_ = c.JSON(http.StatusBadRequest, ErrorResponse{
				Message: appErrors.SomethingWentWrong.Message,
			})
		}
	}
}

func handleValidationError(c echo.Context, err error) bool {
	var ve validator.ValidationErrors
	if !errors.As(err, &ve) {
		return false
	}

	errorsArr := make([]FieldError, 0, len(ve))
	for _, fe := range ve {
		errorsArr = append(errorsArr, FieldError{
			Field:   fe.Field(),
			Message: validationMessage(fe),
		})
	}

	logError(c, "validation failed", nil, "errors", errorsArr)

	_ = c.JSON(http.StatusUnprocessableEntity, ErrorResponse{
		Message:   "Validation failed",
		AppErrors: errorsArr,
	})
	return true
}

func handleAppError(c echo.Context, err error) bool {
	var appErr *appErrors.AppError
	if !errors.As(err, &appErr) {
		return false
	}

	logError(c, "application error", appErr)

	_ = c.JSON(appErr.Code, ErrorResponse{
		Message: appErr.Message,
	})
	return true
}

func handleHTTPError(c echo.Context, err error) bool {
	var he *echo.HTTPError
	if !errors.As(err, &he) {
		return false
	}

	logError(c, "http error", he,
		"code", he.Code,
		"message", he.Message,
	)

	_ = c.JSON(he.Code, ErrorResponse{
		Message: he.Message.(string),
	})
	return true
}

func logError(c echo.Context, msg string, err error, attrs ...any) {
	if config.GetLoggerConfig().EnableRequestLog {
		return
	}

	logger.Error(c, msg, err, attrs...)
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
		return "Must be a valid URL"
	case "unique":
		return "Duplicate values are not allowed"
	case "oneof":
		return "Must be one of: " + strings.ReplaceAll(fe.Param(), " ", ", ")
	default:
		return "Invalid value"
	}
}
