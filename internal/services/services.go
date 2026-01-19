package services

import (
	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/loadBalancer"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type Services struct {
	ApplicationService *ApplicationService
	ClusterService     *ClusterService
	GitService         *GitService
	DockerService      *DockerService
	NodeService        *NodeService
}

func InitServices(
	repos *repositories.Repositories,
	asynqClient *asynq.Client,
	lb *loadBalancer.LoadBalancer,
) (*Services, error) {
	gitService := NewGitService()
	dockerService := NewDockerService()

	sv := &Services{
		ApplicationService: NewApplicationService(repos.ApplicationRepository, gitService),
		ClusterService: NewClusterService(
			repos.ClusterRepository,
			repos.ClusterBuildConfigRepository,
			repos.ApplicationRepository,
			repos.ClusterBuildRepository,
			repos.PodRepository,
			gitService, dockerService,
		),
		GitService:  gitService,
		NodeService: NewNodeService(repos.PodRepository, asynqClient, lb),
	}

	return sv, nil
}
