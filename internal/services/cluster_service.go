package services

import (
	"context"
	"regexp"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
	"github.com/ruhulfbr/mini-k8s/internal/tasks"
)

type ClusterService struct {
	clusterRepo     *repositories.ClusterRepository
	applicationRepo *repositories.ApplicationRepository
	buildConfigRepo *repositories.ClusterBuildConfigRepository
	buildRepo       *repositories.ClusterBuildRepository
	podRepo         *repositories.PodRepository
	gitService      *GitService
	asynqClient     *asynq.Client
}

func NewClusterService(
	clusterRepo *repositories.ClusterRepository,
	buildConfigRepo *repositories.ClusterBuildConfigRepository,
	applicationRepo *repositories.ApplicationRepository,
	buildRepo *repositories.ClusterBuildRepository,
	podRepo *repositories.PodRepository,
	gitService *GitService,
	asynqClient *asynq.Client,
) *ClusterService {
	return &ClusterService{
		clusterRepo:     clusterRepo,
		buildConfigRepo: buildConfigRepo,
		applicationRepo: applicationRepo,
		buildRepo:       buildRepo,
		podRepo:         podRepo,
		gitService:      gitService,
		asynqClient:     asynqClient,
	}
}

func (s *ClusterService) ListByApplication(appId int64, clusterType *string) ([]entities.Cluster, error) {
	if s.applicationRepo.ExistsById(appId) == false {
		return nil, appErrors.NoApplicationFound
	}

	return s.clusterRepo.ListByApplication(appId, clusterType)
}

func (s *ClusterService) GetByID(appId int64, id int64) (*entities.Cluster, error) {
	application, err := s.applicationRepo.GetByID(appId)
	if err != nil || application == nil {
		return nil, appErrors.NoApplicationFound
	}

	cluster, err := s.clusterRepo.GetByAppAndId(appId, id)
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

	application, err := s.applicationRepo.GetByID(cluster.ApplicationId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	if s.clusterRepo.ExistsByName(cluster.ApplicationId, cluster.Name) {
		return appErrors.ClusterAlreadyExist
	}

	if cluster.Type == "" {
		cluster.Type = entities.ClusterTypeHTTP
	}
	if cluster.Replicas < 1 {
		cluster.Replicas = 1
	}

	if err := s.clusterRepo.Create(cluster); err != nil {
		return err
	}

	if bCfg != nil {
		bCfg.ClusterId = cluster.Id
		if err := s.buildConfigRepo.Create(bCfg); err != nil {
			return err
		}

		cluster.BuildConfig = &entities.BuildConfig{
			GitRepo:           bCfg.GitRepo,
			GitBranch:         bCfg.GitBranch,
			DockerContextPath: bCfg.DockerContextPath,
			DockerfileName:    bCfg.DockerfileName,
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

	application, err := s.applicationRepo.GetByID(cluster.ApplicationId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	existing, err := s.clusterRepo.GetByAppAndId(cluster.ApplicationId, cluster.Id)
	if err != nil || existing == nil {
		return appErrors.NoClusterFound
	}

	if s.clusterRepo.ExistsByNameExceptId(cluster.ApplicationId, cluster.Name, cluster.Id) {
		return appErrors.ClusterAlreadyExist
	}

	if err := s.clusterRepo.Update(cluster); err != nil {
		return err
	}

	if bCfg != nil {
		if err := s.updateBuildConfig(bCfg); err != nil {
			logger.Error(nil, "Update build config error", err)
			return err
		}

		cluster.BuildConfig = &entities.BuildConfig{
			GitRepo:           bCfg.GitRepo,
			GitBranch:         bCfg.GitBranch,
			DockerContextPath: bCfg.DockerContextPath,
			DockerfileName:    bCfg.DockerfileName,
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

	ctx := context.Background()
	task := tasks.DeleteClusterTask(&tasks.DeleteClusterPayload{
		ClusterId: id,
	})
	if _, err := s.asynqClient.EnqueueContext(ctx, task); err != nil {
		logger.Error(ctx, "Delete cluster task error", err,
			"ClusterId", id,
		)
		return appErrors.SomethingWentWrong
	}
	return nil
}

func (s *ClusterService) GetBuildHistory(appId int64, id int64) ([]entities.ClusterBuild, error) {
	_, err := s.GetByID(appId, id)
	if err != nil {
		return nil, err
	}

	return s.buildRepo.GetByCluster(id)
}

func (s *ClusterService) BuildDockerImage(appId int64, clusterId int64, version string) error {
	if err := validateVersionText(version); err != nil {
		return err
	}

	application, err := s.applicationRepo.GetByID(appId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	cluster, err := s.clusterRepo.GetByAppAndId(appId, clusterId)
	if err != nil || cluster == nil {
		return appErrors.NoClusterFound
	}

	if cluster.DeployMode != entities.DeployModeBuild {
		return appErrors.InvalidDeployMode
	}

	if s.buildRepo.ExistsByVersion(clusterId, version) {
		return appErrors.DuplicateBuildVersion
	}

	buildConfig, err := s.buildConfigRepo.GetByClusterId(clusterId)
	if err != nil || buildConfig == nil {
		return appErrors.NoBuildConfigFound
	}

	ctx := context.Background()
	task := tasks.BuildDockerImageTask(&tasks.BuildDockerImagePayload{
		ApplicationId:   appId,
		ApplicationName: application.Name,
		ClusterId:       cluster.Id,
		ClusterName:     cluster.Name,
		Version:         version,
		BuildConfig:     *buildConfig,
	})
	if _, err := s.asynqClient.EnqueueContext(ctx, task); err != nil {
		logger.Error(ctx, "Enqueue build docker image task error", err,
			"ApplicationId", appId,
			"ClusterId", clusterId,
			"Version", version,
		)
		return appErrors.SomethingWentWrong
	}

	return nil
}

func (s *ClusterService) PullDockerImage(appId int64, clusterId int64, version string) error {
	if err := validateVersionText(version); err != nil {
		return err
	}

	application, err := s.applicationRepo.GetByID(appId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	cluster, err := s.clusterRepo.GetByAppAndId(appId, clusterId)
	if err != nil || cluster == nil {
		return appErrors.NoClusterFound
	}

	if cluster.DeployMode != entities.DeployModeImage {
		return appErrors.InvalidDeployMode
	}

	if s.buildRepo.ExistsByVersion(clusterId, version) {
		return appErrors.DuplicateBuildVersion
	}

	ctx := context.Background()
	task := tasks.PullDockerImageTask(&tasks.PullDockerImagePayload{
		ApplicationName: application.Name,
		Cluster:         *cluster,
		Version:         version,
	})
	if _, err := s.asynqClient.EnqueueContext(ctx, task); err != nil {
		logger.Error(ctx, "Enqueue pull docker image task error", err,
			"ApplicationId", appId,
			"ClusterId", clusterId,
			"Version", version,
		)
		return appErrors.SomethingWentWrong
	}

	return nil
}

func (s *ClusterService) Deploy(clusterId int64) error {
	cluster, build, err := s.fetchClusterAndBuild(clusterId)
	if err != nil {
		return err
	}

	ctx := context.Background()
	task := tasks.DeployClusterTask(&tasks.DeployClusterPayload{
		Cluster:      *cluster,
		ClusterBuild: *build,
	})
	if _, err := s.asynqClient.EnqueueContext(ctx, task); err != nil {
		logger.Error(ctx, "Enqueue deploy cluster task error", err,
			"ApplicationId", cluster.ApplicationId,
			"ClusterId", clusterId,
			"Replicas", cluster.Replicas,
		)
		return appErrors.SomethingWentWrong
	}

	return nil
}

func (s *ClusterService) RollingDeploy(clusterId int64) error {
	cluster, build, err := s.fetchClusterAndBuild(clusterId)
	if err != nil {
		return err
	}

	ctx := context.Background()
	task := tasks.RollingDeployClusterTask(&tasks.RollingDeployClusterPayload{
		Cluster:      *cluster,
		ClusterBuild: *build,
	})
	if _, err := s.asynqClient.EnqueueContext(ctx, task); err != nil {
		logger.Error(ctx, "Enqueue rolling deploy cluster task error", err,
			"ApplicationId", cluster.ApplicationId,
			"ClusterId", clusterId,
			"Replicas", cluster.Replicas,
		)
		return appErrors.SomethingWentWrong
	}

	return nil
}

func (s *ClusterService) HandleScale(clusterId int64, replicas int) error {
	cluster, build, err := s.fetchClusterAndBuild(clusterId)
	if err != nil {
		return err
	}

	ctx := context.Background()
	task := tasks.ScaleClusterTask(&tasks.ScaleClusterPayload{
		Cluster:      *cluster,
		ClusterBuild: *build,
		Replicas:     replicas,
	})
	if _, err := s.asynqClient.EnqueueContext(ctx, task); err != nil {
		logger.Error(ctx, "Enqueue Scale cluster task error", err,
			"ApplicationId", cluster.ApplicationId,
			"ClusterId", clusterId,
			"Replicas", cluster.Replicas,
		)
		return appErrors.SomethingWentWrong
	}

	return nil
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
