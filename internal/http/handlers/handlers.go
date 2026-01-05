package handlers

import (
	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/database"
	"github.com/ruhulfbr/mini-k8s/internal/loadBalancer"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
	"github.com/ruhulfbr/mini-k8s/internal/services"
)

type Handlers struct {
	ApplicationHandler *ApplicationHandler
	ServiceHandler     *ServiceHandler
	PodHandler         *PodHandler
	NodeHandler        *NodeHandler
}

func InitHandlers(
	cfg *config.Config,
	ds *database.Database,
	asynqClient *asynq.Client,
	lb *loadBalancer.LoadBalancer,
) *Handlers {

	repos := repositories.InitRepositories(ds.DB)
	appServices, _ := services.InitServices(repos, cfg, asynqClient, lb)

	return &Handlers{
		ApplicationHandler: NewApplicationHandler(appServices.ApplicationService),
		ServiceHandler:     NewServiceHandler(appServices.ServiceService),
		PodHandler:         NewPodHandler(appServices.PodService),
		NodeHandler:        NewNodeHandler(appServices.NodeService),
	}
}
