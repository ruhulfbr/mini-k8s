package worker

import (
	"log"
	"time"

	"github.com/hibiken/asynq"
	"github.com/jmoiron/sqlx"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/tasks"
	"github.com/ruhulfbr/mini-k8s/internal/worker/workerHandlers"
	"github.com/ruhulfbr/mini-k8s/internal/worker/workerServices"
)

type Worker struct {
	handlers    workerHandlers.Handlers
	redisConfig config.RedisConfig
}

func NewWorker(DB *sqlx.DB, asynqClient *asynq.Client) *Worker {
	services := workerServices.InitWorkerServices(DB, asynqClient)

	return &Worker{
		handlers:    *workerHandlers.InitWorkerHandlers(services),
		redisConfig: config.GetRedisConfig(),
	}
}

func (w *Worker) StartWorker() {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: w.redisConfig.Host},
		asynq.Config{
			Concurrency: 10,
			RetryDelayFunc: func(n int, e error, t *asynq.Task) time.Duration {
				return time.Duration(n*2) * time.Second // exponential backoff
			},
		},
	)

	mux := asynq.NewServeMux()
	mux.Handle(tasks.BuildDockerImage, w.handlers.ClusterHandler.HandleBuildDockerImage())
	mux.Handle(tasks.PullDockerImage, w.handlers.ClusterHandler.HandlePullDockerImage())
	mux.Handle(tasks.DeployCluster, w.handlers.ClusterHandler.HandleDeployCluster())
	mux.Handle(tasks.RollingDeployCluster, w.handlers.ClusterHandler.HandleRollingDeployCluster())
	mux.Handle(tasks.ScaleCluster, w.handlers.ClusterHandler.HandleScaleCluster())
	mux.Handle(tasks.DeleteCluster, w.handlers.ClusterHandler.HandleDeleteCluster())

	log.Println("[Worker] starting...")

	if err := server.Run(mux); err != nil {
		log.Fatalf("[Worker] failed: %v", err)
	}
}
