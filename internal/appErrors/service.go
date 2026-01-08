package appErrors

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

	InvalidVersionText = &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Invalid version format. Expected vMAJOR.MM.PP (e.g. v1.10.00)",
	}
)
