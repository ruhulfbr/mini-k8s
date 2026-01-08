package services

import (
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type ServiceService struct {
	repo          *repositories.ServiceRepository
	appRepo       *repositories.ApplicationRepository
	buildRepo     *repositories.BuildHistoryRepository
	gitService    *GitService
	dockerService *DockerService
}

func NewServiceService(
	r *repositories.ServiceRepository,
	appRepo *repositories.ApplicationRepository,
	buildRepo *repositories.BuildHistoryRepository,
	gitService *GitService,
	dockerService *DockerService,
) *ServiceService {
	return &ServiceService{repo: r, appRepo: appRepo, buildRepo: buildRepo, gitService: gitService, dockerService: dockerService}
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

func (s *ServiceService) Create(service *entities.Service) error {
	application, err := s.appRepo.GetByID(service.ApplicationId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	if err := s.dockerService.ValidateDockerContext(application.Name, service.ContextPath); err != nil {
		return err
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

	return s.repo.Create(service)
}

func (s *ServiceService) Update(service *entities.Service) error {
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

	if err := s.dockerService.ValidateDockerContext(application.Name, service.ContextPath); err != nil {
		return err
	}

	return s.repo.Update(service)
}

func (s *ServiceService) Delete(appId int64, id int64) error {
	_, err := s.GetByID(appId, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *ServiceService) GetBuildHistory(appId int64, id int64) ([]entities.BuildHistory, error) {
	_, err := s.GetByID(appId, id)
	if err != nil {
		return nil, err
	}

	return s.buildRepo.GetByService(id)
}

func (s *ServiceService) Build(appId int64, id int64, version string) (*entities.BuildHistory, error) {
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

	if s.buildRepo.ExistsByVersion(id, version) {
		return nil, appErrors.DockerDuplicateVersion
	}

	// Pull latest code from git
	err = s.gitService.PullApplication(application.Name, application.GitBranch)
	if err != nil {
		return nil, err
	}

	logger.Info(nil, "git pull application success",
		"application", application.Name,
		"service", service.Name,
		"branch", application.GitBranch,
	)

	// Build image from new codebase
	imageTag, err := s.dockerService.BuildImage(service, application.Name)
	if err != nil {
		return nil, err
	}

	logger.Info(nil, "docker build image success",
		"application", application.Name,
		"service", service.Name,
		"branch", application.GitBranch,
	)

	buildHistory := &entities.BuildHistory{
		ApplicationId: application.Id,
		ServiceId:     service.Id,
		Version:       version,
		ImageTag:      imageTag,
	}

	err = s.buildRepo.Create(buildHistory)
	if err != nil {
		return nil, err
	}

	return buildHistory, nil
}
