package services

import (
	"errors"
	"time"

	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type ServiceService struct {
	repo *repositories.ServiceRepository
}

func NewServiceService(r *repositories.ServiceRepository) *ServiceService {
	return &ServiceService{repo: r}
}

func (s *ServiceService) ListByApplication(appID int64, serviceType *string) ([]entities.Service, error) {
	if appID == 0 {
		return nil, errors.New("invalid application id")
	}
	return s.repo.ListByApplication(appID, serviceType)
}

func (s *ServiceService) Create(service *entities.Service) error {
	if service.ApplicationID == 0 {
		return errors.New("application_id is required")
	}
	if service.Name == "" {
		return errors.New("service name is required")
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
	if service.ID == 0 {
		return errors.New("invalid service id")
	}
	return s.repo.Update(service)
}

func (s *ServiceService) Delete(id int64) error {
	return s.repo.Delete(id)
}

func (s *ServiceService) MarkBuild(serviceID int64) error {
	now := time.Now()
	return s.repo.Update(&entities.Service{
		ID:          serviceID,
		LastBuildAt: &now,
	})
}
