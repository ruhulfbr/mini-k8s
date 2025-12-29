package handlers

import (
	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/datastore"
	"github.com/ruhulfbr/mini-k8s/internal/loadbalancer"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
	"github.com/ruhulfbr/mini-k8s/internal/services"
)

type Handlers struct {
	NodeHandler *NodeHandler
}

func InitHandlers(
	cfg *config.Config,
	ds *datastore.Datastore,
	asynqClient *asynq.Client,
	lb *loadbalancer.LoadBalancer,
) *Handlers {

	repos := repositories.InitRepositories(ds)
	appServices, _ := services.InitServices(repos, cfg, ds, asynqClient, lb)

	return &Handlers{
		NodeHandler: NewNodeHandler(appServices.NodeService),
	}
}
