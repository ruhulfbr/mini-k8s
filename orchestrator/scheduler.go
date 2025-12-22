package orchestrator

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type Scheduler struct {
	ds    *Datastore
	queue *asynq.Client
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

func NewScheduler(ds *Datastore, q *asynq.Client) *Scheduler {
	return &Scheduler{
		ds:    ds,
		queue: q,
	}
}

func (s *Scheduler) Schedule(ctx context.Context, svc ServiceConfig) {
	pods, err := s.ds.ListPods(ctx, svc.ServiceName)
	if err != nil {
		log.Printf("scheduler: failed to list pods: %v", err)
		return
	}

	fmt.Println("Listing pods: ", pods)

	//actual := len(pods)
	actual := 0
	desired := svc.Replicas

	fmt.Println("Desired pods: ", desired)

	type ScaleUpPayload struct {
		ID      string `json:"id"`
		Service string `json:"service"`
	}

	// Scale UP
	for i := actual; i < desired; i++ {
		podID := uuid.NewString()

		fmt.Println("Adding New Task For scale Up: ")

		payload, _ := json.Marshal(DeployPayload{
			ID:      podID,
			Service: svc.ServiceName,
			Image:   svc.Image,
			Port:    svc.Port,
		})

		task := asynq.NewTask("deploy", payload)

		if _, err := s.queue.EnqueueContext(ctx, task); err != nil {
			log.Printf("scheduler: enqueue deploy failed: %v", err)
			continue
		}

		// Optimistic write (like K8s does)
		_ = s.ds.PutPod(ctx, Pod{
			ID:      podID,
			Service: svc.ServiceName,
			Image:   svc.Image,
			Port:    svc.Port,
			Status:  PodPending,
		})
	}

	// Scale DOWN
	for i := desired; i < actual; i++ {
		pod := pods[i]

		fmt.Println("Adding New Task For scale DOWN: ")

		payload, _ := json.Marshal(TerminatePayload{
			ID:      pod.ID,
			Service: pod.Service,
		})

		task := asynq.NewTask("terminate", payload)

		if _, err := s.queue.EnqueueContext(ctx, task); err != nil {
			log.Printf("scheduler: enqueue terminate failed: %v", err)
		}
	}
}
