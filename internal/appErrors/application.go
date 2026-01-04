package appErrors

import (
	"net/http"
)

var (
	ApplicationAlreadyExist = &AppError{
		Code:    http.StatusConflict,
		Message: "An application with this name already exists",
	}

	NoApplicationFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "No application found with this Id",
	}
)
