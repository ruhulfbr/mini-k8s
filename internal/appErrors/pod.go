package appErrors

import "net/http"

var (
	PodCreationFailed = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to create pod",
	}
)
