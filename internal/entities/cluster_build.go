package entities

import "time"

type ClusterBuild struct {
	Id         int64      `db:"id"          json:"id"`
	ClusterId  int64      `db:"cluster_id" json:"cluster_id"`
	Version    string     `db:"version"     json:"version"`
	ImageTag   string     `db:"image_tag"   json:"image_tag"`
	DeployedAt *time.Time `db:"deployed_at"  json:"deployed_at"`
	CreatedAt  time.Time  `db:"created_at"  json:"created_at"`
}
