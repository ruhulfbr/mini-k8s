package services

import (
	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/loadBalancer"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type Services struct {
	ApplicationService *ApplicationService
	ServiceService     *ServiceService
	PodService         *PodService
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
		ServiceService:     NewServiceService(repos.ServiceRepository, repos.ApplicationRepository, repos.BuildHistoryRepository, gitService, dockerService),
		PodService:         NewPodService(repos.PodRepository),
		GitService:         gitService,
		NodeService:        NewNodeService(repos.PodRepository, asynqClient, lb),
	}

	return sv, nil
}
