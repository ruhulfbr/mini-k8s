package requests

type CreatePodRequest struct {
	ApplicationID int64  `json:"application_id" validate:"required"`
	ServiceId     int64  `json:"service_id" validate:"required"`
	Name          string `json:"name" validate:"required"`
	Status        string `json:"status"`
}
