package entities

import "time"

type ServiceBuildConfig struct {
	Id                int64     `db:"id" json:"id"`
	ServiceId         int64     `db:"service_id" json:"service_id"`
	GitRepo           string    `db:"git_repo" json:"git_repo"`
	GitBranch         string    `db:"git_branch" json:"git_branch"`
	DockerContextPath string    `db:"docker_context_path" json:"docker_context_path"`
	DockerfileName    string    `db:"dockerfile_name" json:"dockerfile_name"`
	CreatedAt         time.Time `db:"created_at"  json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at"  json:"updated_at"`
}
