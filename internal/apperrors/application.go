package apperrors

import (
	"net/http"
)

var (
	ApplicationAlreadyExist = &AppError{
		Code:    http.StatusConflict,
		Message: "An application with this name already exists",
	}

	InvalidApplicationId = &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Invalid application Id",
	}
)
