package worker

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/loadBalancer"
	"github.com/ruhulfbr/mini-k8s/internal/tasks"
	"github.com/ruhulfbr/mini-k8s/internal/worker/workerServices"
)

type Worker struct {
	services    workerServices.Services
	redisConfig config.RedisConfig
}

func NewWorker(DB *sqlx.DB, asynqClient *asynq.Client, lb *loadBalancer.LoadBalancer) *Worker {
	return &Worker{
		services:    *workerServices.InitWorkerServices(DB, asynqClient, lb),
		redisConfig: config.GetRedisConfig(),
	}
}

func (w *Worker) StartWorker() {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: w.redisConfig.Host},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.Handle(tasks.BuildDockerImage, w.HandleBuildDockerImage())
	//mux.Handle(tasks.PullDockerImage, w.HandleDeploy())
	//
	//mux.Handle("deploy", w.HandleDeploy())
	//mux.Handle("terminate", w.HandleTerminate())

	log.Println("[Worker] starting...")

	if err := server.Run(mux); err != nil {
		log.Fatalf("[Worker] failed: %v", err)
	}
}

func (w *Worker) HandleBuildDockerImage() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		err := w.services.ClusterService.BuildDockerImage(ctx, t)
		if err != nil {
			logger.Error(ctx, "[Worker] build docker image error: %v", err)
			return err
		}
		return nil
	}
}
