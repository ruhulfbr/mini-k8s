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

type ClusterHandler struct {
	service *services.ClusterService
}

func NewClusterHandler(s *services.ClusterService) *ClusterHandler {
	return &ClusterHandler{service: s}
}

func (h *ClusterHandler) ListByApplication(c echo.Context) error {
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

func (h *ClusterHandler) Show(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	service, err := h.service.GetByID(appId, id)
	if err != nil {
		return err
	}

	return responses.OK(c, service)
}

func (h *ClusterHandler) Create(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)

	req := new(requests.CreateClusterRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	s := &entities.Cluster{
		ApplicationId: appId,
		Name:          req.Name,
		IP:            req.IP,
		Port:          req.Port,
		Replicas:      req.Replicas,
		CPU:           req.CPU,
		Memory:        req.Memory,
		Path:          req.Path,
		Type:          entities.ClusterType(req.Type),
		DeployMode:    entities.DeployMode(req.DeployMode),
	}

	if s.DeployMode == entities.DeployModeImage {
		s.Image = req.Image
	}

	var cfg *entities.ClusterBuildConfig
	if s.DeployMode == entities.DeployModeBuild {
		cfg = &entities.ClusterBuildConfig{
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

func (h *ClusterHandler) Update(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := new(requests.UpdateClusterRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	s := &entities.Cluster{
		Id:            id,
		ApplicationId: appId,
		Name:          req.Name,
		IP:            req.IP,
		Port:          req.Port,
		Replicas:      req.Replicas,
		CPU:           req.CPU,
		Memory:        req.Memory,
		Path:          req.Path,
		Type:          entities.ClusterType(req.Type),
		DeployMode:    entities.DeployMode(req.DeployMode),
	}

	if s.DeployMode == entities.DeployModeImage {
		s.Image = req.Image
	}
	var cfg *entities.ClusterBuildConfig
	if s.DeployMode == entities.DeployModeBuild {
		cfg = &entities.ClusterBuildConfig{
			ClusterId:         s.Id,
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

func (h *ClusterHandler) Delete(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.service.Delete(appId, id)

	if err != nil {
		return err
	}

	return responses.NoContent(c)
}

func (h *ClusterHandler) BuildHistory(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	buildHistories, err := h.service.GetBuildHistory(appId, id)

	if err != nil {
		return err
	}

	return responses.OK(c, buildHistories)
}

func (h *ClusterHandler) BuildImage(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := new(requests.DockerImageBuildRequest)
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

func (h *ClusterHandler) PullImage(c echo.Context) error {
	appId, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := new(requests.DockerImagePullRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	service, err := h.service.PullDockerImage(appId, id, req.Version)
	if err != nil {
		return err
	}

	return responses.OK(c, service)
}
