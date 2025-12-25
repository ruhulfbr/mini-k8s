package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/ruhulfbr/mini-k8s/cmd/api"
	"github.com/ruhulfbr/mini-k8s/internal/datastore"
	"github.com/ruhulfbr/mini-k8s/internal/http/handlers"
	"github.com/ruhulfbr/mini-k8s/internal/loadbalancer"
	"github.com/ruhulfbr/mini-k8s/internal/worker"
)

func main() {
	// ----------------------------------------------------
	// Load Environment Variables
	// ----------------------------------------------------
	loadEnv()

	// ----------------------------------------------------
	// Initialize Core Dependencies
	// ----------------------------------------------------
	ds := initDatastore()
	defer ds.Close()

	asynqClient := initAsynqClient()
	defer asynqClient.Close()

	// ----------------------------------------------------
	// Start Background Services
	// ----------------------------------------------------
	lb := loadbalancer.NewLoadBalancer(ds)
	go lb.Serve()

	startWorker(ds)

	// ----------------------------------------------------
	// Start API Server
	// ----------------------------------------------------
	ctrl := handlers.NewNodeHandler(ds, asynqClient, lb)
	apiServer := api.NewServer(ctrl)

	go apiServer.Start()

	// ----------------------------------------------------
	// Graceful Shutdown
	// ----------------------------------------------------
	waitForShutdown()

	log.Println("[Main] shutting down ", os.Getenv("API_NAME"))
	time.Sleep(2 * time.Second)
}

/* ======================================================
   Helper Functions
====================================================== */

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("[ENV] .env file not found, using system env")
	}
}

func initDatastore() *datastore.Datastore {
	ds := datastore.NewDatastore(os.Getenv("BADGER_DATA_SOURCE"))
	log.Println("[Datastore] initialized")

	return ds
}

func initAsynqClient() *asynq.Client {
	redisAddr := os.Getenv("REDIS_HOST")
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	log.Println("[Asynq] client connected:", redisAddr)
	return client
}

func startWorker(ds *datastore.Datastore) {
	go func() {
		log.Println("[Worker] starting...")
		worker.StartWorker(ds, os.Getenv("REDIS_HOST"))
	}()
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
