package services

import (
	"time"

	"github.com/ruhulfbr/mini-k8s/internal/apperrors"
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
		return nil, apperrors.NoApplicationFound
	}

	return s.repo.ListByApplication(appID, serviceType)
}

func (s *ServiceService) Create(service *entities.Service) error {
	if s.repo.ExistsByName(service.ApplicationId, service.Name) == true {
		return apperrors.ServiceAlreadyExist
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
		return apperrors.NoServiceFound
	}

	return s.repo.Update(service)
}

func (s *ServiceService) Delete(id int64) error {
	if s.repo.ExistsById(id) == false {
		return apperrors.NoServiceFound
	}

	return s.repo.Delete(id)
}

func (s *ServiceService) MarkBuild(id int64) error {
	if s.repo.ExistsById(id) == false {
		return apperrors.NoServiceFound
	}

	now := time.Now()
	return s.repo.Update(&entities.Service{
		Id:          id,
		LastBuildAt: &now,
	})
}
