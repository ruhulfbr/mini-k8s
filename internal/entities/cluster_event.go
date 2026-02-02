package entities

import "time"

type ClusterEvent struct {
	ID        int64     `db:"id" json:"id"`
	ClusterId int64     `db:"cluster_id" json:"clusterId"`
	PodId     *int64    `db:"pod_id" json:"podId,omitempty"`
	Event     string    `db:"event" json:"event"`
	Action    string    `db:"action" json:"action"`
	Message   string    `db:"message" json:"message"`
	Metadata  string    `db:"metadata" json:"metadata"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
}
