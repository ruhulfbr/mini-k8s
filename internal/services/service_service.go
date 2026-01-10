package services

import (
	"regexp"

	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type ServiceService struct {
	repo            *repositories.ServiceRepository
	appRepo         *repositories.ApplicationRepository
	buildConfigRepo *repositories.ServiceBuildConfigRepository
	buildRepo       *repositories.BuildHistoryRepository
	gitService      *GitService
	dockerService   *DockerService
}

func NewServiceService(
	r *repositories.ServiceRepository,
	buildConfigRepo *repositories.ServiceBuildConfigRepository,
	appRepo *repositories.ApplicationRepository,
	buildRepo *repositories.BuildHistoryRepository,
	gitService *GitService,
	dockerService *DockerService,
) *ServiceService {
	return &ServiceService{
		repo:            r,
		buildConfigRepo: buildConfigRepo,
		appRepo:         appRepo,
		buildRepo:       buildRepo,
		gitService:      gitService,
		dockerService:   dockerService,
	}
}

func (s *ServiceService) ListByApplication(appId int64, serviceType *string) ([]entities.Service, error) {
	if s.appRepo.ExistsById(appId) == false {
		return nil, appErrors.NoApplicationFound
	}

	return s.repo.ListByApplication(appId, serviceType)
}

func (s *ServiceService) GetByID(appId int64, id int64) (*entities.Service, error) {
	application, err := s.appRepo.GetByID(appId)
	if err != nil || application == nil {
		return nil, appErrors.NoApplicationFound
	}

	service, err := s.repo.GetById(appId, id)
	if err != nil {
		return nil, err
	}

	if service == nil {
		return nil, appErrors.NoServiceFound
	}

	return service, nil
}

func (s *ServiceService) Create(service *entities.Service, bCfg *entities.ServiceBuildConfig) error {
	application, err := s.appRepo.GetByID(service.ApplicationId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	if s.repo.ExistsByName(service.ApplicationId, service.Name) {
		return appErrors.ServiceAlreadyExist
	}

	if service.Type == "" {
		service.Type = entities.ServiceTypeHTTP
	}
	if service.Replicas < 1 {
		service.Replicas = 1
	}

	if err := s.repo.Create(service); err != nil {
		return err
	}

	if bCfg != nil {
		if err := s.buildConfigRepo.Create(bCfg); err != nil {
			return err
		}

		// Enqueue New Job to clone application and validate docker context + docker file
	}

	return nil
}

func (s *ServiceService) Update(service *entities.Service, bCfg *entities.ServiceBuildConfig) error {
	application, err := s.appRepo.GetByID(service.ApplicationId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	if !s.repo.IsExists(service.ApplicationId, service.Id) {
		return appErrors.NoServiceFound
	}

	if s.repo.ExistsByNameExceptId(service.ApplicationId, service.Name, service.Id) {
		return appErrors.ServiceAlreadyExist
	}

	if err := s.repo.Update(service); err != nil {
		return err
	}

	if bCfg != nil {
		if err := s.buildConfigRepo.Update(bCfg); err != nil {
			return err
		}

		// Enqueue New Job to clone application and validate docker context + docker file
	}

	return nil
}

func (s *ServiceService) Delete(appId int64, id int64) error {
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

func (s *ServiceService) GetBuildHistory(appId int64, id int64) ([]entities.BuildHistory, error) {
	_, err := s.GetByID(appId, id)
	if err != nil {
		return nil, err
	}

	return s.buildRepo.GetByService(id)
}

func (s *ServiceService) BuildDockerImage(appId int64, serviceId int64, version string) (*entities.BuildHistory, error) {
	if err := validateVersionText(version); err != nil {
		return nil, err
	}

	application, err := s.appRepo.GetByID(appId)
	if err != nil || application == nil {
		return nil, appErrors.NoApplicationFound
	}

	service, err := s.repo.GetById(appId, serviceId)
	if err != nil || service == nil {
		return nil, appErrors.NoServiceFound
	}

	if service.DeployMode != entities.DeployModeBuild {
		return nil, appErrors.InvalidDeployMode
	}

	if s.buildRepo.ExistsByVersion(serviceId, version) {
		return nil, appErrors.DockerDuplicateVersion
	}

	buildConfig, err := s.buildConfigRepo.GetByServiceId(serviceId)
	if err != nil || buildConfig == nil {
		return nil, appErrors.NoBuildConfigFound
	}

	if err := s.gitService.PullApplication(application.Name, buildConfig.GitBranch); err != nil {
		return nil, err
	}

	imageTag, err := s.dockerService.BuildImage(buildConfig, application.Name, service.Name)
	if err != nil {
		return nil, err
	}

	buildHistory := &entities.BuildHistory{
		ApplicationId: application.Id,
		ServiceId:     service.Id,
		Version:       version,
		ImageTag:      imageTag,
	}

	if err := s.buildRepo.Create(buildHistory); err != nil {
		return nil, err
	}

	return buildHistory, nil
}

func (s *ServiceService) PullDockerImage(appId int64, serviceId int64, version string, image string) (*entities.BuildHistory, error) {
	if err := validateVersionText(version); err != nil {
		return nil, err
	}

	application, err := s.appRepo.GetByID(appId)
	if err != nil || application == nil {
		return nil, appErrors.NoApplicationFound
	}

	service, err := s.repo.GetById(appId, serviceId)
	if err != nil || service == nil {
		return nil, appErrors.NoServiceFound
	}

	if service.DeployMode != entities.DeployModeImage {
		return nil, appErrors.InvalidDeployMode
	}

	if s.buildRepo.ExistsByImage(serviceId, image) {
		return nil, appErrors.DockerDuplicateImageTag
	}

	if s.buildRepo.ExistsByVersion(serviceId, version) {
		return nil, appErrors.DockerDuplicateVersion
	}

	err = s.dockerService.PullImage(image)
	if err != nil {
		return nil, err
	}

	buildHistory := &entities.BuildHistory{
		ApplicationId: application.Id,
		ServiceId:     service.Id,
		Version:       version,
		ImageTag:      image,
	}

	if err := s.buildRepo.Create(buildHistory); err != nil {
		return nil, err
	}

	return buildHistory, nil
}

// ---------------------- Private Methods ------------------------------

func (s *ServiceService) updateBuildConfig(cfg *entities.ServiceBuildConfig) error {
	exist, err := s.buildConfigRepo.GetByServiceId(cfg.ServiceId)

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
