package requests

type CreateApplicationRequest struct {
	Name    string `json:"name" validate:"required"`
	GitRepo string `json:"git_repo" validate:"required,url"`
}

type UpdateApplicationRequest struct {
	Name    string `json:"name" validate:"required"`
	GitRepo string `json:"git_repo" validate:"required,url"`
}
