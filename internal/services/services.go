package services

import (
	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/loadbalancer"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type Services struct {
	ApplicationService  *ApplicationService
	ServiceService      *ServiceService
	BuildHistoryService *BuildHistoryService
	PodService          *PodService
	NodeService         *NodeService
}

func InitServices(
	repos *repositories.Repositories,
	cfg *config.Config,
	asynqClient *asynq.Client,
	lb *loadbalancer.LoadBalancer,
) (*Services, error) {
	sv := &Services{
		ApplicationService:  NewApplicationService(repos.ApplicationRepository),
		ServiceService:      NewServiceService(repos.ServiceRepository, repos.ApplicationRepository),
		BuildHistoryService: NewBuildHistoryService(repos.BuildHistoryRepository),
		PodService:          NewPodService(repos.PodRepository),
		NodeService:         NewNodeService(repos.PodRepository, asynqClient, lb),
	}

	return sv, nil
}
