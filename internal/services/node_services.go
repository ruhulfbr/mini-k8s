package services

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/entities"
	"github.com/ruhulfbr/mini-k8s/internal/http/requests"
	"github.com/ruhulfbr/mini-k8s/internal/loadBalancer"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

type NodeService struct {
	repo  *repositories.PodRepository
	queue *asynq.Client
	lb    *loadBalancer.LoadBalancer
}

func NewNodeService(
	repo *repositories.PodRepository,
	queue *asynq.Client,
	lb *loadBalancer.LoadBalancer,
) *NodeService {
	return &NodeService{repo: repo, queue: queue, lb: lb}
}

func (s *NodeService) CreateNode(ctx context.Context, req requests.ScaleRequest) error {
	//pods, err := s.clusterRepo.ListPodsByService(ctx, req.ServiceName)
	//if err != nil {
	//	return err
	//}
	//
	//if len(pods) > 0 {
	//	log.Println("node already exists")
	//	return nil
	//}
	//
	//if err := s.scaleUp(ctx, req, req.Replicas); err != nil {
	//	return err
	//}
	//
	//go s.lb.StartServiceListener(req.ServiceName, req.Port)

	return nil
}

func (s *NodeService) Scale(ctx context.Context, req requests.ScaleRequest) error {
	//currentPods, err := s.clusterRepo.ListPodsByService(ctx, req.ServiceName)
	//if err != nil {
	//	return err
	//}
	//
	//currentReplicas := len(currentPods)
	//desiredReplicas := req.Replicas
	//delta := desiredReplicas - currentReplicas
	//
	//// Scale up or down depending on delta
	//switch {
	//case delta > 0:
	//	return s.scaleUp(ctx, req, delta)
	//case delta < 0:
	//	return s.scaleDown(ctx, req.ServiceName, currentPods, -delta)
	//default:
	//	return nil
	//}

	return nil
}

func (s *NodeService) scaleUp(ctx context.Context, req requests.ScaleRequest, count int) error {
	//for i := 0; i < count; i++ {
	//	pod := entities.NewPendingPod(req)
	//
	//	task := jobs.NewDeployTask(pod)
	//	if _, err := s.queue.EnqueueContext(ctx, task); err != nil {
	//		log.Printf("enqueue deploy failed: %v", err)
	//		continue
	//	}
	//
	//	_ = s.clusterRepo.PutPod(ctx, pod)
	//}
	return nil
}

func (s *NodeService) scaleDown(ctx context.Context, service string, pods []entities.Pod, count int) error {
	//for i := 0; i < count && i < len(pods); i++ {
	//	task := jobs.NewTerminateTask(pods[i])
	//	if _, err := s.queue.EnqueueContext(ctx, task); err != nil {
	//		log.Printf("enqueue terminate failed: %v", err)
	//	}
	//}
	//
	//if len(pods) <= count {
	//	go s.lb.StopServiceListener(service)
	//}

	return nil
}
