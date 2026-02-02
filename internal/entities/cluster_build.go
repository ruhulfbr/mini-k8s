package entities

import "time"

type ClusterBuild struct {
	Id         int64      `db:"id"          json:"id"`
	ClusterId  int64      `db:"cluster_id"  json:"clusterId"`
	Version    string     `db:"version"     json:"version"`
	ImageTag   string     `db:"image_tag"   json:"imageTag"`
	DeployedAt *time.Time `db:"deployed_at" json:"deployedAt"`
	CreatedAt  time.Time  `db:"created_at"  json:"createdAt"`
}
