package services

import (
	"regexp"
	"time"

	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type ClusterService struct {
	repo            *repositories.ClusterRepository
	appRepo         *repositories.ApplicationRepository
	buildConfigRepo *repositories.ClusterBuildConfigRepository
	buildRepo       *repositories.ClusterBuildRepository
	podRepo         *repositories.PodRepository
	gitService      *GitService
	dockerService   *DockerService
}

func NewClusterService(
	repo *repositories.ClusterRepository,
	buildConfigRepo *repositories.ClusterBuildConfigRepository,
	appRepo *repositories.ApplicationRepository,
	buildRepo *repositories.ClusterBuildRepository,
	podRepo *repositories.PodRepository,
	gitService *GitService,
	dockerService *DockerService,
) *ClusterService {
	return &ClusterService{
		repo:            repo,
		buildConfigRepo: buildConfigRepo,
		appRepo:         appRepo,
		buildRepo:       buildRepo,
		podRepo:         podRepo,
		gitService:      gitService,
		dockerService:   dockerService,
	}
}

func (s *ClusterService) ListByApplication(appId int64, clusterType *string) ([]entities.Cluster, error) {
	if s.appRepo.ExistsById(appId) == false {
		return nil, appErrors.NoApplicationFound
	}

	return s.repo.ListByApplication(appId, clusterType)
}

func (s *ClusterService) GetByID(appId int64, id int64) (*entities.Cluster, error) {
	application, err := s.appRepo.GetByID(appId)
	if err != nil || application == nil {
		return nil, appErrors.NoApplicationFound
	}

	cluster, err := s.repo.GetByAppAndId(appId, id)
	if err != nil {
		return nil, err
	}

	if cluster == nil {
		return nil, appErrors.NoClusterFound
	}

	return cluster, nil
}

func (s *ClusterService) Create(cluster *entities.Cluster, bCfg *entities.ClusterBuildConfig) error {
	if bCfg != nil {
		if err := s.gitService.ValidateRepoAndBranch(bCfg.GitRepo, bCfg.GitBranch); err != nil {
			return err
		}
	}

	application, err := s.appRepo.GetByID(cluster.ApplicationId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	if s.repo.ExistsByName(cluster.ApplicationId, cluster.Name) {
		return appErrors.ClusterAlreadyExist
	}

	if cluster.Type == "" {
		cluster.Type = entities.ClusterTypeHTTP
	}
	if cluster.Replicas < 1 {
		cluster.Replicas = 1
	}

	if err := s.repo.Create(cluster); err != nil {
		return err
	}

	if bCfg != nil {
		bCfg.ClusterId = cluster.Id
		if err := s.buildConfigRepo.Create(bCfg); err != nil {
			return err
		}
	}

	return nil
}

func (s *ClusterService) Update(cluster *entities.Cluster, bCfg *entities.ClusterBuildConfig) error {
	if bCfg != nil {
		if err := s.gitService.ValidateRepoAndBranch(bCfg.GitRepo, bCfg.GitBranch); err != nil {
			return err
		}
	}

	application, err := s.appRepo.GetByID(cluster.ApplicationId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	existing, err := s.repo.GetByAppAndId(cluster.ApplicationId, cluster.Id)
	if err != nil || existing == nil {
		return appErrors.NoClusterFound
	}

	if s.repo.ExistsByNameExceptId(cluster.ApplicationId, cluster.Name, cluster.Id) {
		return appErrors.ClusterAlreadyExist
	}

	if err := s.repo.Update(cluster); err != nil {
		return err
	}

	if bCfg != nil {
		if err := s.updateBuildConfig(bCfg); err != nil {
			logger.Error(nil, "Update build config error", err)
			return err
		}
	}

	if existing.DeployMode == entities.DeployModeBuild && cluster.DeployMode == entities.DeployModeImage {
		if err := s.buildConfigRepo.Delete(cluster.Id); err != nil {
			logger.Error(nil, "Delete build config error while updated the deploy mode", err)
			return err
		}
	}

	return nil
}

func (s *ClusterService) Delete(appId int64, id int64) error {
	_, err := s.GetByID(appId, id)
	if err != nil {
		return err
	}

	// Delete build Config
	// Delete build History
	// Pods
	// Deleted Images
	// Deleted Containers

	return s.repo.Delete(id)
}

func (s *ClusterService) GetBuildHistory(appId int64, id int64) ([]entities.ClusterBuild, error) {
	_, err := s.GetByID(appId, id)
	if err != nil {
		return nil, err
	}

	return s.buildRepo.GetByCluster(id)
}

func (s *ClusterService) BuildDockerImage(appId int64, clusterId int64, version string) (*entities.ClusterBuild, error) {
	if err := validateVersionText(version); err != nil {
		return nil, err
	}

	application, err := s.appRepo.GetByID(appId)
	if err != nil || application == nil {
		return nil, appErrors.NoApplicationFound
	}

	cluster, err := s.repo.GetByAppAndId(appId, clusterId)
	if err != nil || cluster == nil {
		return nil, appErrors.NoClusterFound
	}

	if cluster.DeployMode != entities.DeployModeBuild {
		return nil, appErrors.InvalidDeployMode
	}

	if s.buildRepo.ExistsByVersion(clusterId, version) {
		return nil, appErrors.DockerDuplicateVersion
	}

	buildConfig, err := s.buildConfigRepo.GetByClusterId(clusterId)
	if err != nil || buildConfig == nil {
		return nil, appErrors.NoBuildConfigFound
	}

	if err := s.gitService.PullApplication(application.Name, cluster.Name, buildConfig); err != nil {
		return nil, err
	}

	imageTag, err := s.dockerService.BuildImage(buildConfig, application.Name, cluster.Name)
	if err != nil {
		return nil, err
	}

	clusterBuild := &entities.ClusterBuild{
		ClusterId: cluster.Id,
		Version:   version,
		ImageTag:  imageTag,
	}

	if err := s.buildRepo.Create(clusterBuild); err != nil {
		return nil, err
	}

	return clusterBuild, nil
}

func (s *ClusterService) PullDockerImage(appId int64, clusterId int64, version string) (*entities.ClusterBuild, error) {
	if err := validateVersionText(version); err != nil {
		return nil, err
	}

	application, err := s.appRepo.GetByID(appId)
	if err != nil || application == nil {
		return nil, appErrors.NoApplicationFound
	}

	cluster, err := s.repo.GetByAppAndId(appId, clusterId)
	if err != nil || cluster == nil {
		return nil, appErrors.NoClusterFound
	}

	if cluster.DeployMode != entities.DeployModeImage {
		return nil, appErrors.InvalidDeployMode
	}

	if s.buildRepo.ExistsByVersion(clusterId, version) {
		return nil, appErrors.DockerDuplicateVersion
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

// ---------------------- Private Methods ------------------------------

func (s *ClusterService) fetchClusterAndBuild(clusterId int64) (*entities.Cluster, *entities.ClusterBuild, error) {
	cluster, err := s.repo.GetById(clusterId)
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
	for i := 0; i < cluster.Replicas; i++ {
		if err := s.scaleDown(cluster.Id, 1); err != nil {
			return err
		}

		if err := s.scaleUp(cluster, build, 1); err != nil {
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

	if err := s.repo.UpdateLatestVersion(cluster); err != nil {
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
