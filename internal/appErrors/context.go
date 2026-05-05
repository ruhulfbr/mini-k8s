package appErrors

import (
	"net/http"
)

var (
	ContextAlreadyExist = &AppError{
		Code:    http.StatusConflict,
		Message: "A context with this name already exists",
	}

	NoContextFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "No context found with this Id",
	}
)
