package services

import (
	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type Services struct {
	ContextService *ContextService
	ClusterService *ClusterService
	GitService     *GitService
	DockerService  *DockerService
}

func InitServices(
	repos *repositories.Repositories,
	asynqClient *asynq.Client,
) (*Services, error) {
	gitService := NewGitService()
	clusterService := NewClusterService(
		repos.ClusterRepository,
		repos.ClusterBuildConfigRepository,
		repos.ContextRepository,
		repos.ClusterBuildRepository,
		repos.PodRepository,
		gitService, asynqClient,
	)

	sv := &Services{
		ContextService: NewContextService(repos.ContextRepository, repos.ClusterRepository),
		ClusterService: clusterService,
		GitService:     gitService,
	}

	return sv, nil
}
