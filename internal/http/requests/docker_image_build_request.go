package requests

type DockerImageBuildRequest struct {
	Version string `json:"version" validate:"required"`
}
