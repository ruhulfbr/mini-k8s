package appErrors

import "net/http"

var (
	GitRepoNotFound = &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: "Invalid or unreachable git repository",
	}

	GitBranchNotExist = &AppError{
		Code:    http.StatusUnprocessableEntity,
		Message: "branch does not exist in repository",
	}

	GitFailedToClone = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to clone repository",
	}

	GitRepoFailedToOpen = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Git repository failed to open",
	}

	GitFailedToOpenWorkTree = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to open git repository worktree",
	}

	GitFailedToCheckout = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to checkout repository branch",
	}

	GitFailedToPull = &AppError{
		Code:    http.StatusBadRequest,
		Message: "Failed to pull from repository",
	}
)
