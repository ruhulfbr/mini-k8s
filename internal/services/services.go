package services

import (
	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/loadBalancer"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type Services struct {
	ApplicationService  *ApplicationService
	ServiceService      *ServiceService
	BuildHistoryService *BuildHistoryService
	PodService          *PodService
	GitService          *GitService
	NodeService         *NodeService
}

func InitServices(
	repos *repositories.Repositories,
	cfg *config.Config,
	asynqClient *asynq.Client,
	lb *loadBalancer.LoadBalancer,
) (*Services, error) {
	gitService := NewGitService(cfg)

	sv := &Services{
		ApplicationService:  NewApplicationService(repos.ApplicationRepository, gitService),
		ServiceService:      NewServiceService(cfg, repos.ServiceRepository, repos.ApplicationRepository),
		BuildHistoryService: NewBuildHistoryService(repos.BuildHistoryRepository),
		PodService:          NewPodService(repos.PodRepository),
		GitService:          gitService,
		NodeService:         NewNodeService(repos.PodRepository, asynqClient, lb),
	}

	return sv, nil
}
