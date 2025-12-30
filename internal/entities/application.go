package entities

import "time"

type Application struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	GitRepo   string    `json:"git_repo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
