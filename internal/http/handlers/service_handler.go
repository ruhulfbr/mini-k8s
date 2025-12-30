package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/http/requests"
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

	return c.JSON(http.StatusOK, serviceList)
}

func (h *ServiceHandler) Create(c echo.Context) error {
	req := new(requests.CreateServiceRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	s := &entities.Service{
		ApplicationID: req.ApplicationID,
		Name:          req.Name,
		IP:            req.IP,
		Port:          req.Port,
		ImageTag:      req.ImageTag,
		ContextPath:   req.ContextPath,
		Replicas:      req.Replicas,
		Resources:     req.Resources,
		Path:          req.Path,
		Type:          entities.ServiceType(req.Type),
	}

	if err := h.service.Create(s); err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, s)
}

func (h *ServiceHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	return h.service.Delete(id)
}
