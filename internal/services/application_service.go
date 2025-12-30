package services

import (
	"github.com/ruhulfbr/mini-k8s/internal/apperrors"
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
	if s.repo.ExistsByName(app.Name) {
		return apperrors.ApplicationAlreadyExist
	}

	return s.repo.Create(app)
}

func (s *ApplicationService) Update(app *entities.Application) error {
	if app.Id < 1 {
		return apperrors.InvalidApplicationId
	}

	return s.repo.Update(app)
}

func (s *ApplicationService) Delete(id int64) error {
	return s.repo.Delete(id)
}
