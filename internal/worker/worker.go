package worker

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/datastore"
)

type Worker struct {
	cfg *config.Config
	ds  *datastore.Datastore
}

func NewWorker(cfg *config.Config, ds *datastore.Datastore) *Worker {
	return &Worker{
		cfg: cfg,
		ds:  ds,
	}
}

func (w *Worker) StartWorker() {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: w.cfg.Redis.Host},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.Handle("deploy", w.HandleDeploy(w.ds))
	mux.Handle("terminate", w.HandleTerminate(w.ds))

	log.Println("[Worker] starting...")
	if err := server.Run(mux); err != nil {
		log.Fatalf("[Worker] failed: %v", err)
	}
}
