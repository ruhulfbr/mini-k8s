package responses

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

const (
	DefaultSuccessMessage = "Success"
	DefaultErrorMessage   = "Something went wrong"
)

type SuccessResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Errors  interface{} `json:"apperrors,omitempty"`
}

func Success(c echo.Context, statusCode int, message string, data ...interface{}) error {
	msg := message
	if msg == "" {
		msg = DefaultSuccessMessage
	}

	var response SuccessResponse
	response.Message = msg

	if len(data) > 0 {
		response.Data = data[0]
	}

	return c.JSON(statusCode, response)
}

func OK(c echo.Context, data interface{}) error {
	return Success(c, http.StatusOK, DefaultSuccessMessage, data)
}

func Created(c echo.Context, data interface{}) error {
	return Success(c, http.StatusCreated, "Created successfully", data)
}

func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

func Accepted(c echo.Context, data interface{}) error {
	return Success(c, http.StatusAccepted, "Accepted", data)
}

func Error(c echo.Context, statusCode int, message string, errors ...interface{}) error {
	msg := message
	if msg == "" {
		msg = DefaultErrorMessage
	}

	var response ErrorResponse
	response.Message = msg

	if len(errors) > 0 {
		response.Errors = errors[0]
	}

	return c.JSON(statusCode, response)
}

func BadRequest(c echo.Context, message string, errors ...interface{}) error {
	return Error(c, http.StatusBadRequest, message, errors...)
}

func ValidationError(c echo.Context, errors interface{}) error {
	return Error(c, http.StatusUnprocessableEntity, "Validation failed", errors)
}

func InternalServerError(c echo.Context) error {
	return Error(c, http.StatusInternalServerError, DefaultErrorMessage)
}
