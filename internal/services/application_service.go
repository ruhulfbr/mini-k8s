package services

import (
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type ApplicationService struct {
	repo       *repositories.ApplicationRepository
	GitService *GitService
}

func NewApplicationService(r *repositories.ApplicationRepository, gs *GitService) *ApplicationService {
	return &ApplicationService{repo: r, GitService: gs}
}

func (s *ApplicationService) List(name *string) ([]entities.Application, error) {
	return s.repo.List(name)
}

func (s *ApplicationService) GetByID(id int64) (*entities.Application, error) {
	app, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if app == nil {
		return nil, appErrors.NotFound
	}

	return app, nil
}

func (s *ApplicationService) Create(app *entities.Application) error {
	if s.repo.ExistsByName(app.Name) {
		return appErrors.ApplicationAlreadyExist
	}

	err := s.repo.Create(app)
	if err != nil {
		_ = s.GitService.RemoveApplicationDir(app.Name)
		return err
	}

	return nil
}

func (s *ApplicationService) Update(app *entities.Application) error {
	if s.repo.ExistsById(app.Id) == false {
		return appErrors.NotFound
	}

	if s.repo.ExistsByNameExceptId(app.Name, app.Id) {
		return appErrors.ApplicationAlreadyExist
	}

	return s.repo.Update(app)
}

func (s *ApplicationService) Delete(id int64) error {
	if s.repo.ExistsById(id) == false {
		return appErrors.NotFound
	}

	return s.repo.Delete(id)
}
