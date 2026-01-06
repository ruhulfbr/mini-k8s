package requests

type CreateServiceRequest struct {
	Name        string `json:"name" validate:"required"`
	IP          string `json:"ip" validate:"required,ip"`
	Port        int    `json:"port"`
	ContextPath string `json:"context_path" validate:"required"`
	Replicas    int    `json:"replicas" validate:"gte=1"`
	CPU         int    `json:"cpu" validate:"required,min=1,max=4"`
	Memory      int    `json:"memory" validate:"required,min=512,max=4096"`
	Path        string `json:"path"`
	Type        string `json:"type" validate:"oneof=http worker"`
}

type UpdateServiceRequest = CreateServiceRequest
