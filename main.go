package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/loadbalancer"
	"github.com/ruhulfbr/mini-k8s/orchestrator"
	"github.com/ruhulfbr/mini-k8s/worker"
)

func main() {
	// ----------------- Initialize Datastore -----------------
	ds := orchestrator.NewDatastore("badger")
	defer ds.Close()

	// ----------------- Load cluster configuration -----------------
	cfg := orchestrator.LoadConfig("cluster.json")

	// ----------------- Redis / Asynq setup -----------------
	redisAddr := "localhost:6379"

	// Asynq client (used by scheduler/controller to enqueue tasks)
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	defer asynqClient.Close()

	// ----------------- Start Worker -----------------
	go func() {
		log.Println("[Worker] starting...")
		worker.StartWorker(ds, redisAddr)
	}()

	// ----------------- Start Controller / Scheduler -----------------
	ctrl := orchestrator.NewController(cfg, orchestrator.NewScheduler(ds, asynqClient))
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		log.Println("[Controller] starting reconciliation loop...")
		ctrl.Run(ctx)
	}()

	// ----------------- Start Load Balancer -----------------
	lb := loadbalancer.NewLoadBalancer(ds)
	go func() {
		log.Println("[LoadBalancer] listening on :8080")

		lb.Serve("8080")
	}()

	// ----------------- Graceful shutdown -----------------
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Println("[Main] shutting down MiniKube-Go...")
	cancel()

	// Give goroutines a moment to exit cleanly
	time.Sleep(2 * time.Second)
}
