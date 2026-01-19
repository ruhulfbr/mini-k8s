package handlers

import (
	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/database"
	"github.com/ruhulfbr/mini-k8s/internal/loadBalancer"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
	"github.com/ruhulfbr/mini-k8s/internal/services"
)

type Handlers struct {
	ApplicationHandler *ApplicationHandler
	ClusterHandler     *ClusterHandler
	NodeHandler        *NodeHandler
}

func InitHandlers(
	ds *database.Database,
	asynqClient *asynq.Client,
	lb *loadBalancer.LoadBalancer,
) *Handlers {

	repos := repositories.InitRepositories(ds.DB)
	appServices, _ := services.InitServices(repos, asynqClient, lb)

	return &Handlers{
		ApplicationHandler: NewApplicationHandler(appServices.ApplicationService),
		ClusterHandler:     NewClusterHandler(appServices.ClusterService),
		NodeHandler:        NewNodeHandler(appServices.NodeService),
	}
}
