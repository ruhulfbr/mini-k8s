package handlers

import (
	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/database"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
	"github.com/ruhulfbr/mini-k8s/internal/services"
)

type Handlers struct {
	ContextHandler *ContextHandler
	ClusterHandler *ClusterHandler
}

func InitHandlers(
	ds *database.Database,
	asynqClient *asynq.Client,
) *Handlers {

	repos := repositories.InitRepositories(ds.DB)
	appServices, _ := services.InitServices(repos, asynqClient)

	return &Handlers{
		ContextHandler: NewContextHandler(appServices.ContextService),
		ClusterHandler: NewClusterHandler(appServices.ClusterService),
	}
}
