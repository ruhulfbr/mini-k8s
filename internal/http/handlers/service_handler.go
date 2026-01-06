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
	appID, _ := strconv.ParseInt(c.Param("appId"), 10, 64)
	typ := c.QueryParam("type")
	var filter *string
	if typ != "" {
		filter = &typ
	}

	serviceList, err := h.service.ListByApplication(appID, filter)

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
		ContextPath:   req.ContextPath,
		Replicas:      req.Replicas,
		CPU:           req.CPU,
		Memory:        req.Memory,
		Path:          req.Path,
		Type:          entities.ServiceType(req.Type),
	}

	if err := h.service.Create(s); err != nil {
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
		ContextPath:   req.ContextPath,
		Replicas:      req.Replicas,
		CPU:           req.CPU,
		Memory:        req.Memory,
		Path:          req.Path,
		Type:          entities.ServiceType(req.Type),
	}

	if err := h.service.Update(s); err != nil {
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
