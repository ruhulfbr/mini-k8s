package services

import (
	"context"
	"log"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/datastore"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/http/requests"
	"github.com/ruhulfbr/mini-k8s/internal/jobs"
	"github.com/ruhulfbr/mini-k8s/internal/loadbalancer"
)

type NodeService struct {
	ds    *datastore.Datastore
	queue *asynq.Client
	lb    *loadbalancer.LoadBalancer
}

func NewNodeService(
	ds *datastore.Datastore,
	queue *asynq.Client,
	lb *loadbalancer.LoadBalancer,
) *NodeService {
	return &NodeService{ds: ds, queue: queue, lb: lb}
}

func (s *NodeService) CreateNode(ctx context.Context, req requests.ScaleRequest) error {
	pods, err := s.ds.ListPodsByService(ctx, req.ServiceName)
	if err != nil {
		return err
	}

	if len(pods) > 0 {
		log.Println("node already exists")
		return nil
	}

	if err := s.scaleUp(ctx, req, req.Replicas); err != nil {
		return err
	}

	go s.lb.StartNodeListener(req.ServiceName, req.Port)
	return nil
}

func (s *NodeService) Scale(ctx context.Context, req requests.ScaleRequest) error {
	pods, err := s.ds.ListPodsByService(ctx, req.ServiceName)
	if err != nil {
		return err
	}

	diff := req.Replicas - len(pods)

	switch {
	case diff > 0:
		return s.scaleUp(ctx, req, diff)
	case diff < 0:
		return s.scaleDown(ctx, pods, -diff)
	}

	return nil
}

func (s *NodeService) scaleUp(ctx context.Context, req requests.ScaleRequest, count int) error {
	for i := 0; i < count; i++ {
		pod := entities.NewPendingPod(req)

		task := jobs.NewDeployTask(pod)
		if _, err := s.queue.EnqueueContext(ctx, task); err != nil {
			log.Printf("enqueue deploy failed: %v", err)
			continue
		}

		_ = s.ds.PutPod(ctx, pod)
	}
	return nil
}

func (s *NodeService) scaleDown(ctx context.Context, pods []entities.Pod, count int) error {
	for i := 0; i < count && i < len(pods); i++ {
		task := jobs.NewTerminateTask(pods[i])
		if _, err := s.queue.EnqueueContext(ctx, task); err != nil {
			log.Printf("enqueue terminate failed: %v", err)
		}
	}
	return nil
}
