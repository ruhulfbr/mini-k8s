package handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/http/requests"
	"github.com/ruhulfbr/mini-k8s/internal/services"
)

type PodHandler struct {
	service *services.PodService
}

func NewPodHandler(s *services.PodService) *PodHandler {
	return &PodHandler{s}
}

func (h *PodHandler) ListByService(c echo.Context) error {
	id, _ := strconv.ParseInt(c.Param("serviceId"), 10, 64)
	status := c.QueryParam("status")
	var filter *string
	if status != "" {
		filter = &status
	}

	pods, err := h.service.ListByService(id, filter)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, pods)
}

func (h *PodHandler) Create(c echo.Context) error {
	req := new(requests.CreatePodRequest)
	if err := c.Bind(req); err != nil {
		return err
	}
	if err := c.Validate(req); err != nil {
		return err
	}

	p := &entities.Pod{
		ApplicationID: req.ApplicationID,
		ServiceID:     req.ServiceID,
		Name:          req.Name,
		Status:        entities.PodPending,
	}

	return h.service.Create(p)
}

func (h *PodHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid pod id")
	}

	if err := h.service.Delete(id); err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
