package apperrors

import "net/http"

var (
	SomethingWentWrong = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Something went wrong",
	}

	InvalidRequestBody = &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Invalid request body",
	}
)
