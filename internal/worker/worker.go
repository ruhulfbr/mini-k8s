package worker

import (
	"log"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/datastore"
)

func StartWorker(ds *datastore.Datastore, redisAddr string) {
	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: redisAddr},
		asynq.Config{Concurrency: 10},
	)

	mux := asynq.NewServeMux()
	mux.Handle("deploy", NewDeployHandler(ds))
	mux.Handle("terminate", NewTerminateHandler(ds))

	log.Println("[Worker] starting...")
	if err := server.Run(mux); err != nil {
		log.Fatalf("[Worker] failed: %v", err)
	}
}
