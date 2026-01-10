package requests

type BuildRequest struct {
	Version string `json:"version" validate:"required"`
}
