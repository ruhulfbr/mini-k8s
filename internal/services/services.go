package services

import (
	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/datastore"
	"github.com/ruhulfbr/mini-k8s/internal/loadbalancer"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type Services struct {
	NodeService *NodeService
}

func InitServices(
	repos *repositories.Repositories,
	cfg *config.Config,
	ds *datastore.Datastore,
	asynqClient *asynq.Client,
	lb *loadbalancer.LoadBalancer,
) (*Services, error) {
	sv := &Services{
		NodeService: NewNodeService(cfg, ds, asynqClient, lb),
	}

	return sv, nil
}
