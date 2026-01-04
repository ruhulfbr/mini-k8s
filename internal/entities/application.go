package entities

import "time"

type Application struct {
	Id        int64     `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	GitRepo   string    `json:"git_repo" db:"git_repo"`
	GitBranch string    `json:"git_branch" db:"git_branch"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
