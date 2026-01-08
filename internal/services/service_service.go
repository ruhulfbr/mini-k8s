package services

import (
	"context"
	"regexp"
	"sync"

	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type ServiceService struct {
	repo          *repositories.ServiceRepository
	appRepo       *repositories.ApplicationRepository
	buildRepo     *repositories.BuildHistoryRepository
	gitService    *GitService
	dockerService *DockerService
}

func NewServiceService(
	r *repositories.ServiceRepository,
	appRepo *repositories.ApplicationRepository,
	buildRepo *repositories.BuildHistoryRepository,
	gitService *GitService,
	dockerService *DockerService,
) *ServiceService {
	return &ServiceService{repo: r, appRepo: appRepo, buildRepo: buildRepo, gitService: gitService, dockerService: dockerService}
}

func (s *ServiceService) ListByApplication(appId int64, serviceType *string) ([]entities.Service, error) {
	if s.appRepo.ExistsById(appId) == false {
		return nil, appErrors.NoApplicationFound
	}

	return s.repo.ListByApplication(appId, serviceType)
}

func (s *ServiceService) GetByID(appId int64, id int64) (*entities.Service, error) {
	application, err := s.appRepo.GetByID(appId)
	if err != nil || application == nil {
		return nil, appErrors.NoApplicationFound
	}

	service, err := s.repo.GetById(appId, id)
	if err != nil {
		return nil, err
	}

	if service == nil {
		return nil, appErrors.NoServiceFound
	}

	return service, nil
}

func (s *ServiceService) Create(service *entities.Service) error {
	application, err := s.appRepo.GetByID(service.ApplicationId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	if err := s.dockerService.ValidateDockerContext(application.Name, service.ContextPath); err != nil {
		return err
	}

	if s.repo.ExistsByName(service.ApplicationId, service.Name) {
		return appErrors.ServiceAlreadyExist
	}

	if service.Type == "" {
		service.Type = entities.ServiceTypeHTTP
	}
	if service.Replicas < 1 {
		service.Replicas = 1
	}

	return s.repo.Create(service)
}

func (s *ServiceService) Update(service *entities.Service) error {
	application, err := s.appRepo.GetByID(service.ApplicationId)
	if err != nil || application == nil {
		return appErrors.NoApplicationFound
	}

	if !s.repo.IsExists(service.ApplicationId, service.Id) {
		return appErrors.NoServiceFound
	}

	if s.repo.ExistsByNameExceptId(service.ApplicationId, service.Name, service.Id) {
		return appErrors.ServiceAlreadyExist
	}

	if err := s.dockerService.ValidateDockerContext(application.Name, service.ContextPath); err != nil {
		return err
	}

	return s.repo.Update(service)
}

func (s *ServiceService) Delete(appId int64, id int64) error {
	_, err := s.GetByID(appId, id)
	if err != nil {
		return err
	}

	return s.repo.Delete(id)
}

func (s *ServiceService) GetBuildHistory(appId int64, id int64) ([]entities.BuildHistory, error) {
	_, err := s.GetByID(appId, id)
	if err != nil {
		return nil, err
	}

	return s.buildRepo.GetByService(id)
}

func (s *ServiceService) Build(appId int64, id int64, version string) (*entities.BuildHistory, error) {
	// Cheap validation first
	if err := validateVersionText(version); err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var (
		application *entities.Application
		service     *entities.Service
	)

	errCh := make(chan error, 1)
	wg := sync.WaitGroup{}
	wg.Add(3)

	// Fetch application
	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		default:
		}

		app, err := s.appRepo.GetByID(appId)
		if err != nil || app == nil {
			errCh <- appErrors.NoApplicationFound
			cancel()
			return
		}
		application = app
	}()

	// Fetch service
	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		default:
		}

		svc, err := s.repo.GetById(appId, id)
		if err != nil {
			errCh <- err
			cancel()
			return
		}
		if svc == nil {
			errCh <- appErrors.NoServiceFound
			cancel()
			return
		}
		service = svc
	}()

	// Check duplicate version
	go func() {
		defer wg.Done()

		select {
		case <-ctx.Done():
			return
		default:
		}

		if s.buildRepo.ExistsByVersion(id, version) {
			errCh <- appErrors.DockerDuplicateVersion
			cancel()
			return
		}
	}()

	// Wait and close error channel
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Return on first error
	for err := range errCh {
		if err != nil {
			return nil, err
		}
	}

	// ---- Sequential steps ----

	// Pull latest code
	if err := s.gitService.PullApplication(application.Name, application.GitBranch); err != nil {
		return nil, err
	}

	logger.Info(nil, "git pull application success",
		"application", application.Name,
		"service", service.Name,
		"branch", application.GitBranch,
	)

	// Build docker image
	imageTag, err := s.dockerService.BuildImage(service, application.Name)
	if err != nil {
		return nil, err
	}

	logger.Info(nil, "docker build image success",
		"application", application.Name,
		"service", service.Name,
		"branch", application.GitBranch,
	)

	buildHistory := &entities.BuildHistory{
		ApplicationId: application.Id,
		ServiceId:     service.Id,
		Version:       version,
		ImageTag:      imageTag,
	}

	if err := s.buildRepo.Create(buildHistory); err != nil {
		return nil, err
	}

	return buildHistory, nil
}

func validateVersionText(version string) error {
	var versionRegex = regexp.MustCompile(
		`^v(0|[1-9]\d*)\.\d{2}\.\d{2}$`,
	)

	if !versionRegex.MatchString(version) {
		return appErrors.InvalidVersionText
	}

	return nil
}
