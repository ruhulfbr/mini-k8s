package entities

import (
	"time"
)

type PodStatus int

const (
	PodPending PodStatus = 0
	PodRunning PodStatus = 1
)

type Pod struct {
	Id        int64     `db:"id"         json:"id"`
	ServiceId int64     `db:"service_id"  json:"service_id"`
	Name      string    `json:"name"`
	Status    PodStatus `json:"status"`
	CreatedAt time.Time `db:"created_at"  json:"created_at"`
}
