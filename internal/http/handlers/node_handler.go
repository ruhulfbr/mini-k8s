package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ruhulfbr/mini-k8s/internal/http/requests"
	"github.com/ruhulfbr/mini-k8s/internal/services"
)

type NodeHandler struct {
	nodeService *services.NodeService
}

func NewNodeHandler(nodeService *services.NodeService) *NodeHandler {
	return &NodeHandler{nodeService: nodeService}
}

func (h *NodeHandler) HandleDeploy(ctx echo.Context) error {
	var req requests.ScaleRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.nodeService.CreateNode(ctx.Request().Context(), req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusAccepted, echo.Map{
		"message": "Deploying operation queued",
	})
}

func (h *NodeHandler) HandleScale(ctx echo.Context) error {
	var req requests.ScaleRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	if err := h.nodeService.Scale(ctx.Request().Context(), req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return ctx.JSON(http.StatusAccepted, echo.Map{
		"message": "Scaling operation queued",
	})
}
