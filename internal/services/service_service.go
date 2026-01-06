package services

import (
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
	"github.com/ruhulfbr/mini-k8s/internal/utils/fsUtils"
)

type ServiceService struct {
	repo    *repositories.ServiceRepository
	appRepo *repositories.ApplicationRepository
}

func NewServiceService(
	r *repositories.ServiceRepository,
	appRepo *repositories.ApplicationRepository,
) *ServiceService {
	return &ServiceService{repo: r, appRepo: appRepo}
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

	app, err := s.repo.GetById(appId, id)
	if err != nil {
		return nil, err
	}

	if app == nil {
		return nil, appErrors.NotFound
	}

	return app, nil
}

func (s *ServiceService) Create(service *entities.Service) error {
	application, err := s.appRepo.GetByID(service.ApplicationId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	if err := s.validateDockerContext(application.Name, service.ContextPath); err != nil {
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

	if err := s.validateDockerContext(application.Name, service.ContextPath); err != nil {
		return err
	}

	return s.repo.Update(service)
}

func (s *ServiceService) Delete(appId int64, id int64) error {
	application, err := s.appRepo.GetByID(appId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	if s.repo.IsExists(appId, id) == false {
		return appErrors.NoServiceFound
	}

	return s.repo.Delete(id)
}

func (s *ServiceService) MarkBuild(appId int64, id int64) error {
	if s.repo.IsExists(appId, id) == false {
		return appErrors.NoServiceFound
	}

	return s.repo.TouchLastBuild(id)
}

func (s *ServiceService) validateDockerContext(appName string, contextPath string) error {
	appPath := fsUtils.Join(config.GetDockerConfig().ApplicationPath, appName)
	if !fsUtils.DirExists(appPath) {
		return appErrors.GirApplicationNotClonedYet
	}

	if !fsUtils.FileExists(fsUtils.Join(appPath, contextPath)) {
		return appErrors.DockerContextFileNotFound
	}

	return nil
}
