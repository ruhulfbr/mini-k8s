package services

import (
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type ContextService struct {
	repo        *repositories.ContextRepository
	clusterRepo *repositories.ClusterRepository
}

func NewContextService(r *repositories.ContextRepository, cr *repositories.ClusterRepository) *ContextService {
	return &ContextService{repo: r, clusterRepo: cr}
}

func (s *ContextService) List(name *string) ([]entities.Context, error) {
	ctxs, err := s.repo.List(name)
	if err != nil {
		return nil, err
	}

	for i := range ctxs {
		clusters, err := s.clusterRepo.ListByContext(ctxs[i].Id, nil)
		if err != nil {
			return nil, err
		}
		ctxs[i].Clusters = clusters
	}

	return ctxs, nil
}

func (s *ContextService) GetByID(id int64) (*entities.Context, error) {
	ctx, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if ctx == nil {
		return nil, appErrors.NoContextFound
	}

	return ctx, nil
}

func (s *ContextService) Create(ctx *entities.Context) error {
	if s.repo.ExistsByName(ctx.Name) {
		return appErrors.ContextAlreadyExist
	}

	return s.repo.Create(ctx)
}

func (s *ContextService) Update(ctx *entities.Context) error {
	if s.repo.ExistsById(ctx.Id) == false {
		return appErrors.NoContextFound
	}

	if s.repo.ExistsByNameExceptId(ctx.Name, ctx.Id) {
		return appErrors.ContextAlreadyExist
	}

	return s.repo.Update(ctx)
}

func (s *ContextService) Delete(id int64) error {
	if s.repo.ExistsById(id) == false {
		return appErrors.NoContextFound
	}

	return s.repo.Delete(id)
}
