package entities

import "time"

type ApplicationStatus int

const (
	ApplicationPending ApplicationStatus = 0
	ApplicationRunning ApplicationStatus = 1
)

type Application struct {
	Id          int64             `json:"id" db:"id"`
	Name        string            `json:"name" db:"name"`
	Description *string           `json:"description" db:"description"`
	Status      ApplicationStatus `json:"status" db:"status"`
	CreatedAt   time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at" db:"updated_at"`
}
