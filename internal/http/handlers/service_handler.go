package handlers

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ruhulfbr/mini-k8s/internal/appErrors"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/http/requests"
	"github.com/ruhulfbr/mini-k8s/internal/http/responses"
	"github.com/ruhulfbr/mini-k8s/internal/services"
)

type ServiceHandler struct {
	service *services.ServiceService
}

func NewServiceHandler(s *services.ServiceService) *ServiceHandler {
	return &ServiceHandler{service: s}
}

func (h *ServiceHandler) ListByApplication(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	typ := c.QueryParam("type")
	var filter *string
	if typ != "" {
		filter = &typ
	}

	serviceList, err := h.service.ListByApplication(appId, filter)

	if err != nil {
		return err
	}

	return responses.OK(c, serviceList)
}

func (h *ServiceHandler) Show(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	service, err := h.service.GetByID(appId, id)
	if err != nil {
		return err
	}

	return responses.OK(c, service)
}

func (h *ServiceHandler) Create(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)

	req := new(requests.CreateServiceRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	s := &entities.Service{
		ApplicationId: appId,
		Name:          req.Name,
		IP:            req.IP,
		Port:          req.Port,
		Replicas:      req.Replicas,
		CPU:           req.CPU,
		Memory:        req.Memory,
		Path:          req.Path,
		Type:          entities.ServiceType(req.Type),
		DeployMode:    entities.DeployMode(req.DeployMode),
	}

	if s.DeployMode == entities.DeployModeImage {
		s.CurrentImageTag = req.Image
	}

	var cfg *entities.ServiceBuildConfig
	if s.DeployMode == entities.DeployModeBuild {
		cfg = &entities.ServiceBuildConfig{
			GitRepo:           req.Build.GitRepo,
			GitBranch:         req.Build.GitBranch,
			DockerContextPath: req.Build.DockerContextPath,
			DockerfileName:    req.Build.DockerfileName,
		}
	}

	if err := h.service.Create(s, cfg); err != nil {
		return err
	}

	return responses.Created(c, s)
}

func (h *ServiceHandler) Update(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := new(requests.UpdateServiceRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	s := &entities.Service{
		Id:            id,
		ApplicationId: appId,
		Name:          req.Name,
		IP:            req.IP,
		Port:          req.Port,
		Replicas:      req.Replicas,
		CPU:           req.CPU,
		Memory:        req.Memory,
		Path:          req.Path,
		Type:          entities.ServiceType(req.Type),
		DeployMode:    entities.DeployMode(req.DeployMode),
	}

	if s.DeployMode == entities.DeployModeImage {
		s.CurrentImageTag = req.Image
	}
	var cfg *entities.ServiceBuildConfig
	if s.DeployMode == entities.DeployModeBuild {
		cfg = &entities.ServiceBuildConfig{
			ServiceId:         s.Id,
			GitRepo:           req.Build.GitRepo,
			GitBranch:         req.Build.GitBranch,
			DockerContextPath: req.Build.DockerContextPath,
			DockerfileName:    req.Build.DockerfileName,
		}
	}

	if err := h.service.Update(s, cfg); err != nil {
		return err
	}

	return responses.OK(c, s)
}

func (h *ServiceHandler) Delete(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.service.Delete(appId, id)

	if err != nil {
		return err
	}

	return responses.NoContent(c)
}

func (h *ServiceHandler) BuildHistory(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	buildHistories, err := h.service.GetBuildHistory(appId, id)

	if err != nil {
		return err
	}

	return responses.OK(c, buildHistories)
}

func (h *ServiceHandler) Build(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := new(requests.BuildRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	service, err := h.service.BuildDockerImage(appId, id, req.Version)
	if err != nil {
		return err
	}

	return responses.OK(c, service)
}
