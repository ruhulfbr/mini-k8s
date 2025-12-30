package services

import (
	"errors"

	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type BuildHistoryService struct {
	repo *repositories.BuildHistoryRepository
}

func NewBuildHistoryService(r *repositories.BuildHistoryRepository) *BuildHistoryService {
	return &BuildHistoryService{repo: r}
}

func (s *BuildHistoryService) Create(history *entities.BuildHistory) error {
	if history.ApplicationID == 0 || history.ServiceID == 0 {
		return errors.New("application_id and service_id are required")
	}
	if history.Tag == "" {
		return errors.New("image tag is required")
	}
	return s.repo.Create(history)
}

func (s *BuildHistoryService) GetByApplication(appID int64) ([]entities.BuildHistory, error) {
	return s.repo.GetByApplication(appID)
}

func (s *BuildHistoryService) GetByService(serviceID int64) ([]entities.BuildHistory, error) {
	return s.repo.GetByService(serviceID)
}

func (s *BuildHistoryService) GetRollbackImage(serviceID int64) (*string, error) {
	return s.repo.GetRollbackImage(serviceID)
}

func (s *BuildHistoryService) Delete(id int64) error {
	return s.repo.Delete(id)
}
