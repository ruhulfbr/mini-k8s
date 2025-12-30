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
	ID            int64
	ApplicationID int64
	ServiceID     int64
	Name          string
	Status        PodStatus
	CreatedAt     time.Time
}
