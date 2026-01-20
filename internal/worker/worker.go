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
	mux.Handle(tasks.PullDockerImage, w.HandlePullDockerImage())
	mux.Handle(tasks.DeployCluster, w.HandleDeployCluster())
	mux.Handle(tasks.RollingDeployCluster, w.HandleRollingDeployCluster())
	mux.Handle(tasks.ScaleCluster, w.HandleScaleCluster())
	mux.Handle(tasks.DeleteCluster, w.HandleDeleteCluster())

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

func (w *Worker) HandlePullDockerImage() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		err := w.services.ClusterService.PullDockerImage(ctx, t)
		if err != nil {
			logger.Error(ctx, "[Worker] Pull docker image error: %v", err)
			return err
		}
		return nil
	}
}

func (w *Worker) HandleDeployCluster() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		err := w.services.ClusterService.Deploy(ctx, t)
		if err != nil {
			logger.Error(ctx, "[Worker] Deploy Cluster error: %v", err)
			return err
		}
		return nil
	}
}

func (w *Worker) HandleRollingDeployCluster() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		err := w.services.ClusterService.RollingDeploy(ctx, t)
		if err != nil {
			logger.Error(ctx, "[Worker] Rolling Deploy Cluster error: %v", err)
			return err
		}
		return nil
	}
}

func (w *Worker) HandleScaleCluster() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		err := w.services.ClusterService.HandleScale(ctx, t)
		if err != nil {
			logger.Error(ctx, "[Worker] Cluster scaling error: %v", err)
			return err
		}
		return nil
	}
}

func (w *Worker) HandleDeleteCluster() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		err := w.services.ClusterService.Delete(ctx, t)
		if err != nil {
			logger.Error(ctx, "[Worker] Delete Cluster error: %v", err)
			return err
		}
		return nil
	}
}
