package server

import (
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/loadBalancer"
	"github.com/ruhulfbr/mini-k8s/internal/worker"
)

func StartWorker(
	DB *sqlx.DB, asynqClient *asynq.Client, lb *loadBalancer.LoadBalancer,
) {
	worker.NewWorker(DB, asynqClient, lb).StartWorker()
}
