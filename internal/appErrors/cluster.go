package appErrors

import "net/http"

var (
	ClusterAlreadyExist = &AppError{
		Code:    http.StatusConflict,
		Message: "A cluster with this name already exists",
	}

	NoClusterFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "No cluster found with this Id",
	}

	NoBuildConfigFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "No build config found with this cluster Id",
	}

	ClusterBuildInfoNotFound = &AppError{
		Code:    http.StatusNotFound,
		Message: "Cluster build info not found",
	}

	InvalidDeployMode = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Invalid deploy mode",
	}

	InvalidVersionText = &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Invalid version format. Expected vMAJOR.MM.PP (e.g. v1.10.00)",
	}
)
