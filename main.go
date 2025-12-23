package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/ruhulfbr/mini-k8s/loadbalancer"
	"github.com/ruhulfbr/mini-k8s/orchestrator"
	"github.com/ruhulfbr/mini-k8s/worker"
)

func main() {

	// --------- Load ENV -----------
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// ----------------- Initialize Datastore -----------------
	ds := orchestrator.NewDatastore(os.Getenv("BADGER_DATA_SOURCE"))
	defer ds.Close()

	// ----------------- Load cluster configuration -----------------
	cfg := orchestrator.LoadConfig("cluster.json")

	// ----------------- Redis / Asynq setup -----------------
	redisAddr := os.Getenv("REDIS_HOST")

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
		log.Println("[LoadBalancer] listening on :", os.Getenv("PORT"))

		lb.Serve(os.Getenv("PORT"))
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
