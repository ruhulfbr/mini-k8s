package workerServices

import (
	"context"
	"strconv"
	"time"

	"github.com/ruhulfbr/mini-k8s/internal/clients"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
	"github.com/ruhulfbr/mini-k8s/internal/services"
	"github.com/ruhulfbr/mini-k8s/internal/tasks"
	"github.com/ruhulfbr/mini-k8s/internal/utils"
	"github.com/ruhulfbr/mini-k8s/internal/worker/events"
)

type ClusterService struct {
	clusterRepo   *repositories.ClusterRepository
	buildRepo     *repositories.ClusterBuildRepository
	eventRepo     *repositories.ClusterEventRepository
	podRepo       *repositories.PodRepository
	gitService    *services.GitService
	dockerService *services.DockerService
	pusher        *clients.PusherClient
}

func NewClusterService(
	clusterRepo *repositories.ClusterRepository,
	buildRepo *repositories.ClusterBuildRepository,
	eventRepo *repositories.ClusterEventRepository,
	podRepo *repositories.PodRepository,
	gitService *services.GitService,
	dockerService *services.DockerService,
) *ClusterService {
	return &ClusterService{
		clusterRepo:   clusterRepo,
		buildRepo:     buildRepo,
		eventRepo:     eventRepo,
		podRepo:       podRepo,
		gitService:    gitService,
		dockerService: dockerService,
		pusher:        clients.NewPusherClient(),
	}
}

func (cs *ClusterService) Delete(ctx context.Context, payload tasks.DeleteClusterPayload) error {
	// Delete build Config
	// Delete build History
	// Pods
	// Deleted Images
	// Deleted Containers
	// Stop load balancer

	return cs.clusterRepo.Delete(payload.ClusterId)
}

func (cs *ClusterService) BuildDockerImage(ctx context.Context, payload tasks.BuildDockerImagePayload) error {
	if err := cs.gitService.PullApplication(payload.ApplicationName, payload.ClusterName, &payload.BuildConfig); err != nil {
		return err
	}

	imageTag, err := cs.dockerService.BuildImage(&payload.BuildConfig, payload.ApplicationName, payload.ClusterName)
	if err != nil {
		return err
	}

	clusterBuild := &entities.ClusterBuild{
		ClusterId: payload.ClusterId,
		Version:   payload.Version,
		ImageTag:  imageTag,
	}

	return cs.buildRepo.Create(clusterBuild)
}

func (cs *ClusterService) PullDockerImage(ctx context.Context, payload tasks.PullDockerImagePayload) error {
	newImageTag, err := cs.dockerService.PullImageWithTag(payload.ApplicationName, &payload.Cluster)
	if err != nil {
		return err
	}

	clusterBuild := &entities.ClusterBuild{
		ClusterId: payload.Cluster.Id,
		Version:   payload.Version,
		ImageTag:  newImageTag,
	}

	return cs.buildRepo.Create(clusterBuild)
}

func (cs *ClusterService) Deploy(ctx context.Context, payload tasks.DeployClusterPayload) error {
	cluster := &payload.Cluster
	clusterId := cluster.Id
	build := &payload.ClusterBuild

	existingPods, err := cs.podRepo.GetByClusterId(clusterId)
	if err != nil {
		return err
	}

	if len(existingPods) == 0 {
		logger.Info(nil, "No pods found for cluster start new deploy",
			"clusterId", clusterId,
		)
		if err := cs.scaleUp(cluster, build, cluster.Replicas); err != nil {
			return err
		}
	} else {
		if err := cs.recreateUpdate(cluster, build, existingPods); err != nil {
			return err
		}
	}

	return cs.updateMetadata(cluster, build)
}

func (cs *ClusterService) RollingDeploy(ctx context.Context, payload tasks.RollingDeployClusterPayload) error {
	cluster := &payload.Cluster
	build := &payload.ClusterBuild

	if err := cs.rollingUpdate(cluster, build); err != nil {
		return err
	}

	return cs.updateMetadata(cluster, build)
}

func (cs *ClusterService) HandleScale(ctx context.Context, payload tasks.ScaleClusterPayload) error {
	cluster := &payload.Cluster
	build := &payload.ClusterBuild
	desiredReplicas := payload.Replicas

	currentPods, err := cs.podRepo.GetByClusterId(cluster.Id)
	if err != nil {
		return err
	}

	currentReplicas := len(currentPods)
	delta := desiredReplicas - currentReplicas

	switch {
	case delta > 0:
		if err := cs.scaleUp(cluster, build, delta); err != nil {
			return err
		}
	case delta < 0:
		if err := cs.scaleDown(cluster.Id, -delta); err != nil {
			return err
		}
	default:
		return nil
	}

	cluster.Replicas = desiredReplicas

	return cs.clusterRepo.Update(cluster)
}

// ---------------------- Private Methods ------------------------------
func (cs *ClusterService) recreateUpdate(cluster *entities.Cluster, build *entities.ClusterBuild, pods []entities.Pod) error {
	if err := cs.terminateAllPods(pods); err != nil {
		return err
	}

	return cs.scaleUp(cluster, build, cluster.Replicas)
}

func (cs *ClusterService) rollingUpdate(cluster *entities.Cluster, build *entities.ClusterBuild) error {
	existingPods, err := cs.podRepo.GetByClusterId(cluster.Id)
	if err != nil {
		return err
	}

	for i := 0; i < cluster.Replicas; i++ {
		if err := cs.deletePod(&existingPods[i]); err != nil {
			return err
		}

		if err := cs.createPod(cluster, build); err != nil {
			return err
		}
	}
	return nil
}

func (cs *ClusterService) scaleUp(cluster *entities.Cluster, build *entities.ClusterBuild, count int) error {
	for i := 0; i < count; i++ {
		if err := cs.createPod(cluster, build); err != nil {
			return err
		}
	}
	return nil
}

func (cs *ClusterService) scaleDown(clusterId int64, count int) error {
	pods, err := cs.podRepo.GetByClusterId(clusterId)
	if err != nil {
		return err
	}

	for i := 0; i < count && i < len(pods); i++ {
		if err := cs.deletePod(&pods[i]); err != nil {
			return err
		}
	}
	return nil
}

func (cs *ClusterService) terminateAllPods(pods []entities.Pod) error {
	for _, pod := range pods {
		err := cs.deletePod(&pod)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cs *ClusterService) createPod(cluster *entities.Cluster, build *entities.ClusterBuild) error {
	info, err := cs.dockerService.DeployImage(cluster, build)
	if err != nil {
		return err
	}

	pod := &entities.Pod{
		ClusterId:     cluster.Id,
		ContainerId:   info.ContainerID,
		ContainerName: info.ContainerName,
		IpAddress:     info.IP,
		Status:        entities.PodStatusPending,
	}

	if err := cs.podRepo.Create(pod); err != nil {
		_ = cs.dockerService.DeleteContainer(*pod)
		return err
	}

	return nil
}

func (cs *ClusterService) deletePod(pod *entities.Pod) error {
	if err := cs.dockerService.DeleteContainer(*pod); err != nil {
		return err
	}
	if err := cs.podRepo.Delete(pod.Id); err != nil {
		return err
	}
	return nil
}

func (cs *ClusterService) updateMetadata(cluster *entities.Cluster, build *entities.ClusterBuild) error {
	now := time.Now()
	imageTag := build.ImageTag
	version := build.Version

	cluster.CurrentImageTag = &imageTag
	cluster.CurrentVersion = &version
	cluster.LastDeployedAt = &now

	if err := cs.clusterRepo.UpdateLatestVersion(cluster); err != nil {
		return err
	}

	build.DeployedAt = &now
	if err := cs.buildRepo.Update(build); err != nil {
		return err
	}
	return nil
}

func (cs *ClusterService) EmitClusterEvent(
	ctx context.Context,
	clusterId int64,
	event events.ClusterEvent,
	action events.ClusterEventAction,
	metadata any,
) {
	eventType := events.GetEventMessage(event, action)
	metaJSON := utils.NormalizeEventMetadata(metadata)

	data := map[string]string{
		"clusterId": strconv.FormatInt(clusterId, 10),
		"event":     string(event),
		"action":    string(action),
		"message":   eventType,
	}

	cs.pusher.TriggerClusterEvent(data)
	if err := cs.eventRepo.LogEvent(&entities.ClusterEvent{
		ClusterId: clusterId,
		Type:      eventType,
		Metadata:  metaJSON,
	}); err != nil {
		logger.Error(ctx, "Failed to emit cluster event", err,
			"clusterId", clusterId,
			"eventType", eventType,
		)
	}
}
