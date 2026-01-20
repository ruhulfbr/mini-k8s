package server

import (
	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/worker"
)

func StartWorker(
	DB *sqlx.DB, asynqClient *asynq.Client,
) {
	worker.NewWorker(DB, asynqClient).StartWorker()
}
