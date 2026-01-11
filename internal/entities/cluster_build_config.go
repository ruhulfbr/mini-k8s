package entities

import "time"

type ClusterBuildConfig struct {
	Id                int64     `db:"id" json:"id"`
	ClusterId         int64     `db:"cluster_id" json:"cluster_id"`
	GitRepo           string    `db:"git_repo" json:"git_repo"`
	GitBranch         string    `db:"git_branch" json:"git_branch"`
	DockerContextPath string    `db:"docker_context_path" json:"docker_context_path"`
	DockerfileName    string    `db:"dockerfile_name" json:"dockerfile_name"`
	CreatedAt         time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"  json:"updated_at"`
}
