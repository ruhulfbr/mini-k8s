package workerServices

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
	"github.com/ruhulfbr/mini-k8s/internal/services"
	"github.com/ruhulfbr/mini-k8s/internal/tasks"
)

type ClusterService struct {
	clusterRepo   *repositories.ClusterRepository
	buildRepo     *repositories.ClusterBuildRepository
	podRepo       *repositories.PodRepository
	gitService    *services.GitService
	dockerService *services.DockerService
}

func NewClusterService(
	clusterRepo *repositories.ClusterRepository,
	buildRepo *repositories.ClusterBuildRepository,
	podRepo *repositories.PodRepository,
	gitService *services.GitService,
	dockerService *services.DockerService,
) *ClusterService {
	return &ClusterService{
		clusterRepo:   clusterRepo,
		buildRepo:     buildRepo,
		podRepo:       podRepo,
		gitService:    gitService,
		dockerService: dockerService,
	}
}

func (s *ClusterService) Delete(ctx context.Context, task *asynq.Task) error {
	var payload tasks.DeleteClusterPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	// Delete build Config
	// Delete build History
	// Pods
	// Deleted Images
	// Deleted Containers
	// Stop load balancer

	return s.clusterRepo.Delete(payload.ClusterId)
}

func (s *ClusterService) BuildDockerImage(ctx context.Context, task *asynq.Task) error {
	var payload tasks.BuildDockerImagePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	if err := s.gitService.PullApplication(payload.ApplicationName, payload.ClusterName, &payload.BuildConfig); err != nil {
		return err
	}

	imageTag, err := s.dockerService.BuildImage(&payload.BuildConfig, payload.ApplicationName, payload.ClusterName)
	if err != nil {
		return err
	}

	clusterBuild := &entities.ClusterBuild{
		ClusterId: payload.ClusterId,
		Version:   payload.Version,
		ImageTag:  imageTag,
	}

	return s.buildRepo.Create(clusterBuild)
}

func (s *ClusterService) PullDockerImage(ctx context.Context, task *asynq.Task) error {
	var payload tasks.PullDockerImagePayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	newImageTag, err := s.dockerService.PullImageWithTag(payload.ApplicationName, &payload.Cluster)
	if err != nil {
		return err
	}

	clusterBuild := &entities.ClusterBuild{
		ClusterId: payload.Cluster.Id,
		Version:   payload.Version,
		ImageTag:  newImageTag,
	}

	return s.buildRepo.Create(clusterBuild)
}

func (s *ClusterService) Deploy(ctx context.Context, task *asynq.Task) error {
	var payload tasks.DeployClusterPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	cluster := &payload.Cluster
	clusterId := cluster.Id
	build := &payload.ClusterBuild

	existingPods, err := s.podRepo.GetByClusterId(clusterId)
	if err != nil {
		return err
	}

	if len(existingPods) == 0 {
		logger.Info(nil, "No pods found for cluster start new deploy",
			"clusterId", clusterId,
		)
		if err := s.scaleUp(cluster, build, cluster.Replicas); err != nil {
			return err
		}
	} else {
		if err := s.recreateUpdate(cluster, build, existingPods); err != nil {
			return err
		}
	}

	return s.updateMetadata(cluster, build)
}

func (s *ClusterService) RollingDeploy(ctx context.Context, task *asynq.Task) error {
	var payload tasks.RollingDeployClusterPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	cluster := &payload.Cluster
	build := &payload.ClusterBuild

	if err := s.rollingUpdate(cluster, build); err != nil {
		return err
	}

	return s.updateMetadata(cluster, build)
}

func (s *ClusterService) HandleScale(ctx context.Context, task *asynq.Task) error {
	var payload tasks.ScaleClusterPayload
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return err
	}

	cluster := &payload.Cluster
	build := &payload.ClusterBuild
	desiredReplicas := payload.Replicas

	currentPods, err := s.podRepo.GetByClusterId(cluster.Id)
	if err != nil {
		return err
	}

	currentReplicas := len(currentPods)
	delta := desiredReplicas - currentReplicas

	switch {
	case delta > 0:
		if err := s.scaleUp(cluster, build, delta); err != nil {
			return err
		}
	case delta < 0:
		if err := s.scaleDown(cluster.Id, -delta); err != nil {
			return err
		}
	default:
		return nil
	}

	cluster.Replicas = desiredReplicas

	return s.clusterRepo.Update(cluster)
}

// ---------------------- Private Methods ------------------------------
func (s *ClusterService) recreateUpdate(cluster *entities.Cluster, build *entities.ClusterBuild, pods []entities.Pod) error {
	if err := s.terminateAllPods(pods); err != nil {
		return err
	}

	return s.scaleUp(cluster, build, cluster.Replicas)
}

func (s *ClusterService) rollingUpdate(cluster *entities.Cluster, build *entities.ClusterBuild) error {
	existingPods, err := s.podRepo.GetByClusterId(cluster.Id)
	if err != nil {
		return err
	}

	for i := 0; i < cluster.Replicas; i++ {
		if err := s.deletePod(&existingPods[i]); err != nil {
			return err
		}

		if err := s.createPod(cluster, build); err != nil {
			return err
		}
	}
	return nil
}

func (s *ClusterService) scaleUp(cluster *entities.Cluster, build *entities.ClusterBuild, count int) error {
	for i := 0; i < count; i++ {
		if err := s.createPod(cluster, build); err != nil {
			return err
		}
	}
	return nil
}

func (s *ClusterService) scaleDown(clusterId int64, count int) error {
	pods, err := s.podRepo.GetByClusterId(clusterId)
	if err != nil {
		return err
	}

	for i := 0; i < count && i < len(pods); i++ {
		if err := s.deletePod(&pods[i]); err != nil {
			return err
		}
	}
	return nil
}

func (s *ClusterService) terminateAllPods(pods []entities.Pod) error {
	for _, pod := range pods {
		err := s.deletePod(&pod)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *ClusterService) createPod(cluster *entities.Cluster, build *entities.ClusterBuild) error {
	info, err := s.dockerService.DeployImage(cluster, build)
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

	if err := s.podRepo.Create(pod); err != nil {
		_ = s.dockerService.DeleteContainer(*pod)
		return err
	}

	return nil
}

func (s *ClusterService) deletePod(pod *entities.Pod) error {
	if err := s.dockerService.DeleteContainer(*pod); err != nil {
		return err
	}
	if err := s.podRepo.Delete(pod.Id); err != nil {
		return err
	}
	return nil
}

func (s *ClusterService) updateMetadata(cluster *entities.Cluster, build *entities.ClusterBuild) error {
	now := time.Now()
	imageTag := build.ImageTag
	version := build.Version

	cluster.CurrentImageTag = &imageTag
	cluster.CurrentVersion = &version
	cluster.LastDeployedAt = &now

	if err := s.clusterRepo.UpdateLatestVersion(cluster); err != nil {
		return err
	}

	build.DeployedAt = &now
	if err := s.buildRepo.Update(build); err != nil {
		return err
	}
	return nil
}
