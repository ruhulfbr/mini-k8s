package entities

import "time"

type BuildHistory struct {
	ID            int64
	ApplicationID int64
	ServiceID     int64
	Tag           string
	CreatedAt     time.Time
}
