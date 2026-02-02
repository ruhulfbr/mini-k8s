package entities

import "time"

type AppStatus int

const (
	AppPending AppStatus = 0
	AppRunning AppStatus = 1
)

type Application struct {
	Id          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description *string   `json:"description" db:"description"`
	Status      AppStatus `json:"status" db:"status"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
	Clusters    []Cluster `json:"clusters,omitempty"`
}
