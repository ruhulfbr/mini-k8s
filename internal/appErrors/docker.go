package appErrors

import (
	"net/http"
)

var (
	DockerContextFileNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "Docker context file not found",
	}

	DockerDuplicateVersion = &AppError{
		Code:    http.StatusNotFound,
		Message: "Docker build version already exists",
	}
)
