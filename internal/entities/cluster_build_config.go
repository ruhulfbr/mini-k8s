package entities

import (
	"time"

	"github.com/ruhulfbr/mini-k8s/internal/http/requests"
)

type ClusterBuildConfig struct {
	Id                int64     `db:"id" json:"id"`
	ClusterId         int64     `db:"cluster_id" json:"clusterId"`
	GitRepo           string    `db:"git_repo" json:"gitRepo"`
	GitBranch         string    `db:"git_branch" json:"gitBranch"`
	DockerContextPath string    `db:"docker_context_path" json:"dockerContextPath"`
	DockerfileName    string    `db:"dockerfile_name" json:"dockerfileName"`
	CreatedAt         time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time `db:"updated_at" json:"updatedAt"`
}

func ClusterBuildConfigFromRequest(req *requests.BuildConfigRequest, clusterId ...int64) *ClusterBuildConfig {
	cfg := &ClusterBuildConfig{
		GitRepo:           req.GitRepo,
		GitBranch:         req.GitBranch,
		DockerContextPath: req.DockerContextPath,
		DockerfileName:    req.DockerfileName,
	}

	if len(clusterId) > 0 {
		cfg.ClusterId = clusterId[0]
	}

	return cfg
}
