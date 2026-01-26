package services

import (
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type ApplicationService struct {
	repo        *repositories.ApplicationRepository
	clusterRepo *repositories.ClusterRepository
}

func NewApplicationService(r *repositories.ApplicationRepository, cr *repositories.ClusterRepository) *ApplicationService {
	return &ApplicationService{repo: r, clusterRepo: cr}
}

func (s *ApplicationService) List(name *string) ([]entities.Application, error) {
	apps, err := s.repo.List(name)
	if err != nil {
		return nil, err
	}

	for i := range apps {
		clusters, err := s.clusterRepo.ListByApplication(apps[i].Id, nil)
		if err != nil {
			return nil, err
		}
		apps[i].Clusters = clusters
	}

	return apps, nil
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

	return s.repo.Create(app)
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
