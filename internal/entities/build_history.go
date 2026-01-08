package entities

import "time"

type BuildHistory struct {
	Id            int64     `db:"id"          json:"id"`
	ApplicationId int64     `db:"application_id" json:"application_id"`
	ServiceId     int64     `db:"service_id"  json:"service_id"`
	Version       string    `db:"version"     json:"version"`
	ImageTag      string    `db:"image_tag"   json:"image_tag"`
	CreatedAt     time.Time `db:"created_at"  json:"created_at"`
}
