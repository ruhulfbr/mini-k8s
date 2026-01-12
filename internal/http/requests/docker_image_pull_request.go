package requests

type DockerImagePullRequest struct {
	Version string `json:"version" validate:"required"`
}
