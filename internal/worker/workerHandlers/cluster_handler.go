package workerHandlers

import (
	"context"
	"encoding/json"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/tasks"
	"github.com/ruhulfbr/mini-k8s/internal/worker/events"
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
	return func(ctx context.Context, task *asynq.Task) error {
		var payload tasks.BuildDockerImagePayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			logger.Error(ctx, "[Worker] [BuildDockerImage] invalid payload", err)
			return nil
		}

		ch.ClusterService.EmitClusterEvent(ctx, payload.ClusterId, events.EventBuildDockerImage, events.ActionStarted, nil)

		err := ch.ClusterService.BuildDockerImage(ctx, payload)
		if err != nil {
			ch.ClusterService.EmitClusterEvent(ctx, payload.ClusterId, events.EventBuildDockerImage, events.ActionFailed, err)
			logger.Error(ctx, "[Worker] build docker image error: %v", err)
			return nil
		}

		ch.ClusterService.EmitClusterEvent(ctx, payload.ClusterId, events.EventBuildDockerImage, events.ActionSuccess, nil)

		return nil
	}
}

func (ch *ClusterHandler) HandlePullDockerImage() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload tasks.PullDockerImagePayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			logger.Error(ctx, "[Worker] [PullDockerImage] invalid payload", err)
			return nil
		}

		ch.ClusterService.EmitClusterEvent(ctx, payload.Cluster.Id, events.EventPullDockerImage, events.ActionStarted, nil)

		err := ch.ClusterService.PullDockerImage(ctx, payload)
		if err != nil {
			ch.ClusterService.EmitClusterEvent(ctx, payload.Cluster.Id, events.EventPullDockerImage, events.ActionFailed, err)
			logger.Error(ctx, "[Worker] Pull docker image error: %v", err)
			return nil
		}

		ch.ClusterService.EmitClusterEvent(ctx, payload.Cluster.Id, events.EventPullDockerImage, events.ActionSuccess, nil)

		return nil
	}
}

func (ch *ClusterHandler) HandleDeployCluster() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload tasks.DeployClusterPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			logger.Error(ctx, "[Worker] [DeployCluster] invalid payload", err)
			return nil
		}

		ch.ClusterService.EmitClusterEvent(ctx, payload.Cluster.Id, events.EventDeploy, events.ActionStarted, nil)

		err := ch.ClusterService.Deploy(ctx, payload)
		if err != nil {
			ch.ClusterService.EmitClusterEvent(ctx, payload.Cluster.Id, events.EventDeploy, events.ActionFailed, err)
			logger.Error(ctx, "[Worker] Deploy Cluster error: %v", err)
			return nil
		}

		ch.ClusterService.EmitClusterEvent(ctx, payload.Cluster.Id, events.EventDeploy, events.ActionSuccess, nil)

		return nil
	}
}

func (ch *ClusterHandler) HandleRollingDeployCluster() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload tasks.RollingDeployClusterPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			logger.Error(ctx, "[Worker] [RollingDeployCluster] invalid payload", err)
			return nil
		}

		ch.ClusterService.EmitClusterEvent(ctx, payload.Cluster.Id, events.EventRollingDeploy, events.ActionStarted, nil)

		err := ch.ClusterService.RollingDeploy(ctx, payload)
		if err != nil {
			ch.ClusterService.EmitClusterEvent(ctx, payload.Cluster.Id, events.EventRollingDeploy, events.ActionFailed, err)
			logger.Error(ctx, "[Worker] Rolling Deploy Cluster error: %v", err)
			return nil
		}

		ch.ClusterService.EmitClusterEvent(ctx, payload.Cluster.Id, events.EventRollingDeploy, events.ActionFailed, nil)

		return nil
	}
}

func (ch *ClusterHandler) HandleScaleCluster() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload tasks.ScaleClusterPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			logger.Error(ctx, "[Worker] [ScaleCluster] invalid payload", err)
			return nil
		}

		ch.ClusterService.EmitClusterEvent(ctx, payload.Cluster.Id, events.EventScaling, events.ActionStarted, nil)

		err := ch.ClusterService.HandleScale(ctx, payload)
		if err != nil {
			ch.ClusterService.EmitClusterEvent(ctx, payload.Cluster.Id, events.EventScaling, events.ActionFailed, err)
			logger.Error(ctx, "[Worker] Cluster scaling error: %v", err)
			return nil
		}

		ch.ClusterService.EmitClusterEvent(ctx, payload.Cluster.Id, events.EventScaling, events.ActionSuccess, nil)

		return nil
	}
}

func (ch *ClusterHandler) HandleDeleteCluster() asynq.HandlerFunc {
	return func(ctx context.Context, task *asynq.Task) error {
		var payload tasks.DeleteClusterPayload
		if err := json.Unmarshal(task.Payload(), &payload); err != nil {
			logger.Error(ctx, "[Worker] [DeleteCluster] invalid payload", err)
			return nil
		}

		err := ch.ClusterService.Delete(ctx, payload)
		if err != nil {
			logger.Error(ctx, "[Worker] Delete Cluster error: %v", err)
			return nil
		}
		return nil
	}
}
