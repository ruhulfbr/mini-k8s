package orchestrator

import (
	"context"
	"time"
)

type Controller struct {
	cfg       *ClusterConfig
	scheduler *Scheduler
}

func NewController(cfg *ClusterConfig, s *Scheduler) *Controller {
	return &Controller{
		cfg:       cfg,
		scheduler: s,
	}
}

func (c *Controller) Run(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			for _, svc := range c.cfg.GetServices() {
				c.scheduler.Schedule(ctx, svc)
			}
		case <-ctx.Done():
			return
		}
	}
}
