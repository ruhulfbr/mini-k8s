package appErrors

import (
	"net/http"
)

var (
	DockerContextFileNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "Docker context file not found",
	}
)
