package services

import (
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
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

func (s *ServiceService) ListByApplication(appID int64, serviceType *string) ([]entities.Service, error) {
	if s.appRepo.ExistsById(appID) == false {
		return nil, appErrors.NoApplicationFound
	}

	return s.repo.ListByApplication(appID, serviceType)
}

func (s *ServiceService) Create(service *entities.Service) error {
	if s.appRepo.ExistsById(service.ApplicationId) == false {
		return appErrors.NoApplicationFound
	}

	if s.repo.ExistsByName(service.ApplicationId, service.Name) {
		return appErrors.ServiceAlreadyExist
	}
	
	if service.Type == "" {
		service.Type = entities.ServiceTypeHTTP
	}
	if service.Replicas <= 0 {
		service.Replicas = 1
	}

	return s.repo.Create(service)
}

func (s *ServiceService) Update(service *entities.Service) error {
	if s.repo.ExistsById(service.Id) == false {
		return appErrors.NoServiceFound
	}

	if s.repo.ExistsByNameExceptId(service.ApplicationId, service.Name, service.Id) {
		return appErrors.ServiceAlreadyExist
	}

	return s.repo.Update(service)
}

func (s *ServiceService) Delete(id int64) error {
	if s.repo.ExistsById(id) == false {
		return appErrors.NoServiceFound
	}

	return s.repo.Delete(id)
}

func (s *ServiceService) MarkBuild(id int64) error {
	if s.repo.ExistsById(id) == false {
		return appErrors.NoServiceFound
	}

	return s.repo.TouchLastBuild(id)
}
