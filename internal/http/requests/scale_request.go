package requests

type ScaleRequest struct {
	ServiceName string `json:"service_name"`
	Image       string `json:"image"`
	Replicas    *int   `json:"replicas" validate:"required,min=0"`
	Port        int    `json:"port"`
}
