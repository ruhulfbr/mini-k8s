package requests

type CreateContextRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=32"`
	Description *string `json:"description,omitempty"`
}

type UpdateContextRequest = CreateContextRequest
