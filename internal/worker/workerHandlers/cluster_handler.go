package workerHandlers

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/worker/workerServices"
)

type ClusterHandler struct {
	ClusterService *workerServices.ClusterService
}

func NewClusterHandler(clusterService *workerServices.ClusterService) *ClusterHandler {
	return &ClusterHandler{
		ClusterService: clusterService,
	}
}

func (ch *ClusterHandler) HandleBuildDockerImage() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		err := ch.ClusterService.BuildDockerImage(ctx, t)
		if err != nil {
			logger.Error(ctx, "[Worker] build docker image error: %v", err)
		}
		return nil
	}
}

func (ch *ClusterHandler) HandlePullDockerImage() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		err := ch.ClusterService.PullDockerImage(ctx, t)
		if err != nil {
			logger.Error(ctx, "[Worker] Pull docker image error: %v", err)
			return err
		}
		return nil
	}
}

func (ch *ClusterHandler) HandleDeployCluster() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		err := ch.ClusterService.Deploy(ctx, t)
		if err != nil {
			logger.Error(ctx, "[Worker] Deploy Cluster error: %v", err)
			return err
		}
		return nil
	}
}

func (ch *ClusterHandler) HandleRollingDeployCluster() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		err := ch.ClusterService.RollingDeploy(ctx, t)
		if err != nil {
			logger.Error(ctx, "[Worker] Rolling Deploy Cluster error: %v", err)
			return err
		}
		return nil
	}
}

func (ch *ClusterHandler) HandleScaleCluster() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		err := ch.ClusterService.HandleScale(ctx, t)
		if err != nil {
			logger.Error(ctx, "[Worker] Cluster scaling error: %v", err)
			return err
		}
		return nil
	}
}

func (ch *ClusterHandler) HandleDeleteCluster() asynq.HandlerFunc {
	return func(ctx context.Context, t *asynq.Task) error {
		err := ch.ClusterService.Delete(ctx, t)
		if err != nil {
			logger.Error(ctx, "[Worker] Delete Cluster error: %v", err)
			return err
		}
		return nil
	}
}
