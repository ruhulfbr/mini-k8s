package workerHandlers

import (
	"github.com/ruhulfbr/mini-k8s/internal/worker/workerServices"
)

type Handlers struct {
	ClusterHandler *ClusterHandler
}

func InitWorkerHandlers(services *workerServices.Services) *Handlers {
	return &Handlers{
		ClusterHandler: NewClusterHandler(services.ClusterService),
	}
}
