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

type ApplicationHandler struct {
	service *services.ApplicationService
}

func NewApplicationHandler(s *services.ApplicationService) *ApplicationHandler {
	return &ApplicationHandler{service: s}
}

func (h *ApplicationHandler) List(c echo.Context) error {
	name := c.QueryParam("name")
	var filter *string
	if name != "" {
		filter = &name
	}

	apps, err := h.service.List(filter)
	if err != nil {
		return err
	}

	return responses.OK(c, apps)
}

func (h *ApplicationHandler) Show(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	app, err := h.service.GetByID(id)
	if err != nil {
		return err
	}

	return responses.OK(c, app)
}

func (h *ApplicationHandler) Create(c echo.Context) error {
	req := new(requests.CreateApplicationRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	app := &entities.Application{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := h.service.Create(app); err != nil {
		return err
	}

	return responses.Created(c, app)
}

func (h *ApplicationHandler) Update(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := new(requests.UpdateApplicationRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	app := &entities.Application{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := h.service.Update(app); err != nil {
		return err
	}

	return responses.OK(c, app)
}

func (h *ApplicationHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.service.Delete(id)
	if err != nil {
		return err
	}

	return responses.NoContent(c)
}
