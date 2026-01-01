package requests

type CreateApplicationRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=32"`
	GitRepo   string `json:"git_repo" validate:"required,url"`
	GitBranch string `json:"git_branch" validate:"required,min=1,max=128"`
}

type UpdateApplicationRequest struct {
	Name      string `json:"name" validate:"required"`
	GitRepo   string `json:"git_repo" validate:"required,url"`
	GitBranch string `json:"git_branch" validate:"required,min=1,max=128"`
}
