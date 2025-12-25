package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/ruhulfbr/mini-k8s/internal/datastore"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/loadbalancer"
)

type NodeHandler struct {
	ds      *datastore.Datastore
	queue   *asynq.Client
	lb      *loadbalancer.LoadBalancer
	counter uint64
}

type ScaleRequest struct {
	ServiceName string `json:"service_name"`
	Image       string `json:"image"`
	Replicas    int    `json:"replicas"`
	Port        int    `json:"port"`
}

type DeployPayload struct {
	ID      string `json:"id"`
	Service string `json:"service"`
	Image   string `json:"image"`
	Port    int    `json:"port"`
}

type TerminatePayload struct {
	ID      string `json:"id"`
	Service string `json:"service"`
}

func NewNodeHandler(
	ds *datastore.Datastore,
	client *asynq.Client,
	lb *loadbalancer.LoadBalancer,
) *NodeHandler {
	return &NodeHandler{
		ds:    ds,
		queue: client,
		lb:    lb,
	}
}

func (c *NodeHandler) HandleDeploy(ctx echo.Context) error {
	var req ScaleRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if err := c.createNode(ctx.Request().Context(), req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusAccepted, echo.Map{
		"message": "Deploying operation queued",
	})
}

func (c *NodeHandler) HandleScale(ctx echo.Context) error {
	var req ScaleRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	if err := c.scale(ctx.Request().Context(), req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{
			"error": err.Error(),
		})
	}

	return ctx.JSON(http.StatusAccepted, echo.Map{
		"message": "Scaling operation queued",
	})
}

func (c *NodeHandler) createNode(ctx context.Context, req ScaleRequest) error {
	pods, err := c.ds.ListPodsByService(ctx, req.ServiceName)
	if err != nil {
		log.Printf("scheduler: failed to list Node: %v", err)
		return err
	}

	if len(pods) > 0 {
		log.Printf("scheduler: Node already exists")
		return nil
	}

	err = c.scaleUp(ctx, req, req.Replicas)
	if err != nil {
		log.Printf("Something went wrong while creating node: %v", err)
		return err
	}

	// Start Load Balancer for newly created node
	go c.lb.RunForNode(req.ServiceName, req.Port)

	return nil
}

func (c *NodeHandler) scale(ctx context.Context, req ScaleRequest) error {
	pods, err := c.ds.ListPodsByService(ctx, req.ServiceName)
	if err != nil {
		log.Printf("scheduler: failed to list pods: %v", err)
		return err
	}

	current := len(pods)
	desired := req.Replicas

	switch {
	case current < desired:
		return c.scaleUp(ctx, req, desired-current)

	case current > desired:
		return c.scaleDown(ctx, pods, current-desired)
	}

	return nil
}

func (c *NodeHandler) scaleUp(ctx context.Context, req ScaleRequest, count int) error {
	for i := 0; i < count; i++ {
		podID := uuid.NewString()

		fmt.Println("Adding New Task For scale Up: ")

		payload, _ := json.Marshal(DeployPayload{
			ID:      podID,
			Service: req.ServiceName,
			Image:   req.Image,
			Port:    req.Port,
		})

		task := asynq.NewTask("deploy", payload)

		if _, err := c.queue.EnqueueContext(ctx, task); err != nil {
			log.Printf("scheduler: enqueue deploy failed: %v", err)
			continue
		}

		// Optimistic write (like K8s does)
		_ = c.ds.PutPod(ctx, entities.Pod{
			ID:      podID,
			Service: req.ServiceName,
			Image:   req.Image,
			Port:    req.Port,
			Status:  entities.PodPending,
		})
	}

	return nil
}

func (c *NodeHandler) scaleDown(ctx context.Context, pods []entities.Pod, count int) error {
	for i := 0; i < count && i < len(pods); i++ {
		pod := pods[i]

		fmt.Println("Adding New Task For scale DOWN: ")

		payload, _ := json.Marshal(TerminatePayload{
			ID:      pod.ID,
			Service: pod.Service,
		})

		task := asynq.NewTask("terminate", payload)

		if _, err := c.queue.EnqueueContext(ctx, task); err != nil {
			log.Printf("scheduler: enqueue terminate failed: %v", err)
		}
	}

	return nil
}
