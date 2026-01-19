package appErrors

import (
	"net/http"
)

var (
	DockerInvalidContextPath = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid docker context path",
	}

	DockerFileNotFound = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Dockerfile file not found",
	}

	DockerDuplicateImageTag = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Docker image tag already exists",
	}

	DockerDuplicateVersion = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Docker build version already exists",
	}

	DockerFailedToBuildImage = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to build Docker image",
	}

	DockerFailedToPullImage = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to pull Docker image",
	}

	DockerContainerHasNoIP = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Docker container has no IP",
	}

	DockerFailedReadPullResponse = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to read Docker pull response",
	}

	DockerFailedAddNewTagToImage = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to add new tag to Docker image",
	}

	DockerFailedCreateContainer = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to create Docker container",
	}

	DockerFailedRunContainer = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to run Docker container",
	}

	DockerFailedRetrieveContainerIP = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to retrieve Docker container IP",
	}

	DockerFailedDeleteContainer = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to delete Docker container",
	}
)
