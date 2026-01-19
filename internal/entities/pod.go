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
	Id            int64     `db:"id"         json:"id"`
	ClusterId     int64     `db:"cluster_id" json:"cluster_id"`
	ContainerId   string    `db:"container_id" json:"container_id"`
	ContainerName string    `db:"container_name" json:"container_name"`
	IpAddress     string    `db:"ip_address" json:"ip_address"`
	CreatedAt     time.Time `db:"created_at"  json:"created_at"`
}
