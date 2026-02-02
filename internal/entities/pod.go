package entities

import (
	"time"
)

type PodStatus int

var (
	PodStatusPending PodStatus = 0
	PodStatusRunning PodStatus = 1
)

type Pod struct {
	Id            int64     `db:"id" json:"id"`
	ClusterId     int64     `db:"cluster_id" json:"clusterId"`
	ContainerId   string    `db:"container_id" json:"containerId"`
	ContainerName string    `db:"container_name" json:"containerName"`
	IpAddress     string    `db:"ip_address" json:"ipAddress"`
	Status        PodStatus `db:"status" json:"status"`
	CreatedAt     time.Time `db:"created_at" json:"createdAt"`
}
