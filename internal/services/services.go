package services

import (
	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type Services struct {
	ApplicationService *ApplicationService
	ClusterService     *ClusterService
	GitService         *GitService
	DockerService      *DockerService
}

func InitServices(
	repos *repositories.Repositories,
	asynqClient *asynq.Client,
) (*Services, error) {
	gitService := NewGitService()

	sv := &Services{
		ApplicationService: NewApplicationService(repos.ApplicationRepository, gitService),
		ClusterService: NewClusterService(
			repos.ClusterRepository,
			repos.ClusterBuildConfigRepository,
			repos.ApplicationRepository,
			repos.ClusterBuildRepository,
			repos.PodRepository,
			gitService, asynqClient,
		),
		GitService: gitService,
	}

	return sv, nil
}
