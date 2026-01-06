package requests

type CreateServiceRequest struct {
	ApplicationId int64  `json:"application_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
	IP            string `json:"ip" validate:"required,ip"`
	Port          int    `json:"port"`
	ImageTag      string `json:"image_tag" validate:"required"`
	ContextPath   string `json:"context_path" validate:"required"`
	Replicas      int    `json:"replicas" validate:"gte=1"`
	Resources     string `json:"resources" validate:"required"`
	Path          string `json:"path"`
	Type          string `json:"type" validate:"oneof=http worker"`
}

type UpdateServiceRequest = CreateServiceRequest
