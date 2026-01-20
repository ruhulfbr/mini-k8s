package workerServices

import (
	"context"
	"encoding/json"
	"regexp"
	"time"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
	"github.com/ruhulfbr/mini-k8s/internal/services"
	"github.com/ruhulfbr/mini-k8s/internal/tasks"
)

type ClusterService struct {
	clusterRepo     *repositories.ClusterRepository
	applicationRepo *repositories.ApplicationRepository
	buildConfigRepo *repositories.ClusterBuildConfigRepository
	buildRepo       *repositories.ClusterBuildRepository
	podRepo         *repositories.PodRepository
	gitService      *services.GitService
	dockerService   *services.DockerService
}

func NewClusterService(
	clusterRepo *repositories.ClusterRepository,
	buildConfigRepo *repositories.ClusterBuildConfigRepository,
	applicationRepo *repositories.ApplicationRepository,
	buildRepo *repositories.ClusterBuildRepository,
	podRepo *repositories.PodRepository,
	gitService *services.GitService,
	dockerService *services.DockerService,
) *ClusterService {
	return &ClusterService{
		clusterRepo:     clusterRepo,
		buildConfigRepo: buildConfigRepo,
		applicationRepo: applicationRepo,
		buildRepo:       buildRepo,
		podRepo:         podRepo,
		gitService:      gitService,
		dockerService:   dockerService,
	}
}

func (s *ClusterService) Delete(appId int64, id int64) error {
	//_, err := s.GetByID(appId, id)
	//if err != nil {
	//	return err
	//}

	// Delete build Config
	// Delete build History
	// Pods
	// Deleted Images
	// Deleted Containers
	// Stop load balancer

	return s.clusterRepo.Delete(id)
}

func (s *ClusterService) BuildDockerImage(ctx context.Context, t *asynq.Task) error {
	var payload tasks.BuildDockerImagePayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
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

	if err := s.buildRepo.Create(clusterBuild); err != nil {
		return err
	}

	return nil
}

func (s *ClusterService) PullDockerImage(appId int64, clusterId int64, version string) (*entities.ClusterBuild, error) {
	if err := validateVersionText(version); err != nil {
		return nil, err
	}

	application, err := s.applicationRepo.GetByID(appId)
	if err != nil || application == nil {
		return nil, appErrors.NoApplicationFound
	}

	cluster, err := s.clusterRepo.GetByAppAndId(appId, clusterId)
	if err != nil || cluster == nil {
		return nil, appErrors.NoClusterFound
	}

	if cluster.DeployMode != entities.DeployModeImage {
		return nil, appErrors.InvalidDeployMode
	}

	if s.buildRepo.ExistsByVersion(clusterId, version) {
		return nil, appErrors.DuplicateBuildVersion
	}

	newImageTag, err := s.dockerService.PullImageWithTag(application.Name, cluster)
	if err != nil {
		return nil, err
	}

	clusterBuild := &entities.ClusterBuild{
		ClusterId: cluster.Id,
		Version:   version,
		ImageTag:  newImageTag,
	}

	if err := s.buildRepo.Create(clusterBuild); err != nil {
		return nil, err
	}

	return clusterBuild, nil
}

func (s *ClusterService) Deploy(clusterId int64) error {
	cluster, build, err := s.fetchClusterAndBuild(clusterId)
	if err != nil {
		return err
	}

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

func (s *ClusterService) RollingDeploy(clusterId int64) error {
	cluster, build, err := s.fetchClusterAndBuild(clusterId)
	if err != nil {
		return err
	}

	if err := s.rollingUpdate(cluster, build); err != nil {
		return err
	}

	return s.updateMetadata(cluster, build)
}

func (s *ClusterService) HandleScale(clusterId int64, replicas int) error {
	cluster, build, err := s.fetchClusterAndBuild(clusterId)
	if err != nil {
		return err
	}

	currentPods, err := s.podRepo.GetByClusterId(clusterId)
	if err != nil {
		return err
	}

	currentReplicas := len(currentPods)
	desiredReplicas := replicas
	delta := desiredReplicas - currentReplicas

	switch {
	case delta > 0:
		if err := s.scaleUp(cluster, build, delta); err != nil {
			return err
		}
	case delta < 0:
		if err := s.scaleDown(clusterId, -delta); err != nil {
			return err
		}
	default:
		return nil
	}

	cluster.Replicas = replicas

	return s.clusterRepo.Update(cluster)
}

// ---------------------- Private Methods ------------------------------

func (s *ClusterService) fetchClusterAndBuild(clusterId int64) (*entities.Cluster, *entities.ClusterBuild, error) {
	cluster, err := s.clusterRepo.GetById(clusterId)
	if err != nil || cluster == nil {
		return nil, nil, appErrors.NoClusterFound
	}

	build, err := s.buildRepo.GetLatestBuild(clusterId)
	if err != nil {
		return nil, nil, appErrors.ClusterBuildInfoNotFound
	}

	return cluster, build, nil
}

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

func (s *ClusterService) updateBuildConfig(cfg *entities.ClusterBuildConfig) error {
	exist, err := s.buildConfigRepo.GetByClusterId(cfg.ClusterId)

	if err != nil || exist == nil {
		return s.buildConfigRepo.Create(cfg)
	}

	return s.buildConfigRepo.Update(cfg)
}

func validateVersionText(version string) error {
	var versionRegex = regexp.MustCompile(
		`^v(0|[1-9]\d*)\.\d{2}\.\d{2}$`,
	)

	if !versionRegex.MatchString(version) {
		return appErrors.InvalidVersionText
	}

	return nil
}
