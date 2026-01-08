package services

import (
	"errors"

	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type PodService struct {
	repo *repositories.PodRepository
}

func NewPodService(r *repositories.PodRepository) *PodService {
	return &PodService{repo: r}
}

func (s *PodService) ListByService(serviceID int64, status *string) ([]entities.Pod, error) {
	if serviceID == 0 {
		return nil, errors.New("invalid service id")
	}
	return s.repo.ListByService(serviceID, status)
}

func (s *PodService) Create(pod *entities.Pod) error {
	if pod.ApplicationId == 0 || pod.ServiceId == 0 {
		return errors.New("application_id and service_id are required")
	}
	if pod.Name == "" {
		return errors.New("pod name is required")
	}
	pod.Status = entities.PodPending

	return s.repo.Create(pod)
}

func (s *PodService) Update(pod *entities.Pod) error {
	if pod.Id == 0 {
		return errors.New("invalid pod id")
	}
	return s.repo.Update(pod)
}

func (s *PodService) Delete(id int64) error {
	return s.repo.Delete(id)
}
