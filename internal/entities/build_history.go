package entities

import "time"

type BuildHistory struct {
	Id            int64
	ApplicationId int64
	ServiceId     int64
	Version       string
	ImageTag      string
	CreatedAt     time.Time
}
