package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type Controller struct {
	ds    *Datastore
	queue *asynq.Client
}

type ScaleRequest struct {
	ServiceName string `json:"service_name"`
	Image       string `json:"image"`
	Replicas    int    `json:"replicas"`
	Port        int    `json:"port"`
}

func NewController(ds *Datastore, client *asynq.Client) *Controller {
	return &Controller{ds: ds, queue: client}
}

func (c *Controller) Scale(ctx context.Context, req ScaleRequest) error {
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

func (c *Controller) scaleUp(ctx context.Context, req ScaleRequest, count int) error {
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
		_ = c.ds.PutPod(ctx, Pod{
			ID:      podID,
			Service: req.ServiceName,
			Image:   req.Image,
			Port:    req.Port,
			Status:  PodPending,
		})
	}

	return nil
}

func (c *Controller) scaleDown(ctx context.Context, pods []Pod, count int) error {
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
