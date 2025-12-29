package server

import (
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/datastore"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
	"github.com/ruhulfbr/mini-k8s/internal/worker"
)

func startWorker(
	cfg *config.Config,
	ds *datastore.Datastore,
) {
	podRepository := repositories.NewPodRepository(ds.DB)
	worker.NewWorker(cfg, podRepository).StartWorker()
}
