package server

import (
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/database"
	"github.com/ruhulfbr/mini-k8s/internal/loadBalancer"
	"github.com/ruhulfbr/mini-k8s/internal/repositories"
)

func InitLoadBalancer(ds *database.Database) *loadBalancer.LoadBalancer {
	podRepository := repositories.NewPodRepository(ds.DB)

	lb := loadBalancer.NewLoadBalancer(podRepository)
	// lb.Start()

	return lb
}
