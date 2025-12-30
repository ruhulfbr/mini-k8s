package services

import (
	"errors"

	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type ApplicationService struct {
	repo *repositories.ApplicationRepository
}

func NewApplicationService(r *repositories.ApplicationRepository) *ApplicationService {
	return &ApplicationService{repo: r}
}

func (s *ApplicationService) List(name *string) ([]entities.Application, error) {
	return s.repo.List(name)
}

func (s *ApplicationService) GetByID(id int64) (*entities.Application, error) {
	return s.repo.GetByID(id)
}

func (s *ApplicationService) Create(app *entities.Application) error {
	if app.Name == "" || app.GitRepo == "" {
		return errors.New("name and git_repo are required")
	}
	return s.repo.Create(app)
}

func (s *ApplicationService) Update(app *entities.Application) error {
	if app.ID == 0 {
		return errors.New("invalid application id")
	}
	return s.repo.Update(app)
}

func (s *ApplicationService) Delete(id int64) error {
	return s.repo.Delete(id)
}
