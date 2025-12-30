package entities

import "time"

type Application struct {
	ID        int64
	Name      string
	GitRepo   string
	CreatedAt time.Time
	UpdatedAt time.Time
}
