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

func (h *ClusterHandler) ListByContext(c echo.Context) error {
	ctxId, _ := strconv.ParseInt(c.Param("ctxId"), 10, 64)
	typ := c.QueryParam("type")
	var filter *string
	if typ != "" {
		filter = &typ
	}

	serviceList, err := h.service.ListByContext(ctxId, filter)

	if err != nil {
		return err
	}

	return responses.OK(c, serviceList)
}

func (h *ClusterHandler) Show(c echo.Context) error {
	ctxId, _ := strconv.ParseInt(c.Param("ctxId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	service, err := h.service.GetByID(ctxId, id)
	if err != nil {
		return err
	}

	return responses.OK(c, service)
}

func (h *ClusterHandler) Create(c echo.Context) error {
	ctxId, _ := strconv.ParseInt(c.Param("ctxId"), 10, 64)

	req := new(requests.CreateClusterRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	cl, err := entities.ClusterFromRequest(ctxId, req)
	if err != nil {
		return err
	}

	var cfg *entities.ClusterBuildConfig
	if cl.DeployMode == entities.DeployModeBuild {
		cfg = entities.ClusterBuildConfigFromRequest(req.Build)
	}

	if err := h.service.Create(cl, cfg); err != nil {
		return err
	}

	return responses.Created(c, cl)
}

func (h *ClusterHandler) Update(c echo.Context) error {
	ctxId, _ := strconv.ParseInt(c.Param("ctxId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := new(requests.UpdateClusterRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	cl, err := entities.ClusterFromRequest(ctxId, req, id)
	if err != nil {
		return err
	}
	var cfg *entities.ClusterBuildConfig
	if cl.DeployMode == entities.DeployModeBuild {
		cfg = entities.ClusterBuildConfigFromRequest(req.Build, id)
	}

	if err := h.service.Update(cl, cfg); err != nil {
		return err
	}

	return responses.OK(c, cl)
}

func (h *ClusterHandler) Delete(c echo.Context) error {
	ctxId, _ := strconv.ParseInt(c.Param("ctxId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.service.Delete(ctxId, id)

	if err != nil {
		return err
	}

	return responses.NoContent(c)
}

func (h *ClusterHandler) BuildHistory(c echo.Context) error {
	ctxId, _ := strconv.ParseInt(c.Param("ctxId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	buildHistories, err := h.service.GetBuildHistory(ctxId, id)

	if err != nil {
		return err
	}

	return responses.OK(c, buildHistories)
}

func (h *ClusterHandler) BuildImage(c echo.Context) error {
	ctxId, _ := strconv.ParseInt(c.Param("ctxId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := new(requests.DockerImageBuildRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	err := h.service.BuildDockerImage(ctxId, id, req.Version)
	if err != nil {
		return err
	}

	return responses.OK(c, nil)
}

func (h *ClusterHandler) PullImage(c echo.Context) error {
	ctxId, _ := strconv.ParseInt(c.Param("ctxId"), 10, 64)
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := new(requests.DockerImagePullRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	err := h.service.PullDockerImage(ctxId, id, req.Version)
	if err != nil {
		return err
	}

	return responses.OK(c, nil)
}

func (h *ClusterHandler) Deploy(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.service.Deploy(id)
	if err != nil {
		return err
	}

	return responses.OK(c, nil)
}

func (h *ClusterHandler) RollingDeploy(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.service.RollingDeploy(id)
	if err != nil {
		return err
	}

	return responses.OK(c, nil)
}

func (h *ClusterHandler) HandleScale(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := new(requests.ScaleRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	err := h.service.HandleScale(id, *req.Replicas)
	if err != nil {
		return err
	}

	return responses.OK(c, nil)
}
