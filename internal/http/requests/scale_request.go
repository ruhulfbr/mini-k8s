package requests

type ScaleRequest struct {
	Replicas *int `json:"replicas" validate:"required,min=0"`
}
