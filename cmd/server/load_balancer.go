package server

import (
	"github.com/ruhulfbr/mini-k8s/internal/database"
	"github.com/ruhulfbr/mini-k8s/internal/loadbalancer"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

func InitLoadBalancer(ds *database.Database) *loadbalancer.LoadBalancer {
	podRepository := repositories.NewPodRepository(ds.DB)

	lb := loadbalancer.NewLoadBalancer(podRepository)
	// lb.Start()

	return lb
}
