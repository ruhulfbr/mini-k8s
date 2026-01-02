package apperrors

import "net/http"

var (
	ServiceAlreadyExist = &AppError{
		Code:    http.StatusConflict,
		Message: "An service with this name already exists",
	}

	NoServiceFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "No service found with this Id",
	}
)
