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

type ContextHandler struct {
	service *services.ContextService
}

func NewContextHandler(s *services.ContextService) *ContextHandler {
	return &ContextHandler{service: s}
}

func (h *ContextHandler) List(c echo.Context) error {
	name := c.QueryParam("name")
	var filter *string
	if name != "" {
		filter = &name
	}

	ctxs, err := h.service.List(filter)
	if err != nil {
		return err
	}

	return responses.OK(c, ctxs)
}

func (h *ContextHandler) Show(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	ctx, err := h.service.GetByID(id)
	if err != nil {
		return err
	}

	return responses.OK(c, ctx)
}

func (h *ContextHandler) Create(c echo.Context) error {
	req := new(requests.CreateContextRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}

	if err := c.Validate(req); err != nil {
		return err
	}

	ctx := &entities.Context{
		Name:        req.Name,
		Description: req.Description,
	}
	if err := h.service.Create(ctx); err != nil {
		return err
	}

	return responses.Created(c, ctx)
}

func (h *ContextHandler) Update(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	req := new(requests.UpdateContextRequest)
	if err := c.Bind(req); err != nil {
		return appErrors.InvalidRequestBody
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	ctx := &entities.Context{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
	}
	if err := h.service.Update(ctx); err != nil {
		return err
	}

	return responses.OK(c, ctx)
}

func (h *ContextHandler) Delete(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("id"), 10, 64)

	err := h.service.Delete(id)
	if err != nil {
		return err
	}

	return responses.NoContent(c)
}
