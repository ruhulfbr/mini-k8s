package worker

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type Worker struct {
	cfg     *config.Config
	podRepo *repositories.PodRepository
}

func NewWorker(cfg *config.Config, podRepo *repositories.PodRepository) *Worker {
	return &Worker{
		cfg:     cfg,
		podRepo: podRepo,
	}
}

func (w *Worker) StartWorker() {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: w.cfg.Redis.Host},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.Handle("deploy", w.HandleDeploy())
	mux.Handle("terminate", w.HandleTerminate())

	log.Println("[Worker] starting...")
	if err := server.Run(mux); err != nil {
		log.Fatalf("[Worker] failed: %v", err)
	}
}
