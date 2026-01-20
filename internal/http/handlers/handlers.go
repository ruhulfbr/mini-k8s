package handlers

import (
	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/database"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
	"github.com/ruhulfbr/mini-k8s/internal/services"
)

type Handlers struct {
	ApplicationHandler *ApplicationHandler
	ClusterHandler     *ClusterHandler
}

func InitHandlers(
	ds *database.Database,
	asynqClient *asynq.Client,
) *Handlers {

	repos := repositories.InitRepositories(ds.DB)
	appServices, _ := services.InitServices(repos, asynqClient)

	return &Handlers{
		ApplicationHandler: NewApplicationHandler(appServices.ApplicationService),
		ClusterHandler:     NewClusterHandler(appServices.ClusterService),
	}
}
