package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"
	"github.com/ruhulfbr/mini-k8s/cmd/api"
	"github.com/ruhulfbr/mini-k8s/internal/datastore"
	"github.com/ruhulfbr/mini-k8s/internal/http/handlers"
	"github.com/ruhulfbr/mini-k8s/internal/loadbalancer"
	"github.com/ruhulfbr/mini-k8s/internal/services"
	"github.com/ruhulfbr/mini-k8s/internal/worker"
)

func main() {
	loadEnv()

	ds := initDatastore()
	defer ds.Close()

	asynqClient := initAsynqClient()
	defer asynqClient.Close()

	lb := loadbalancer.NewLoadBalancer(ds)
	go lb.Start()

	startWorker(ds)

	startAPIServer(ds, asynqClient, lb)

	waitForShutdown()
	log.Println("[Main] shutting down", os.Getenv("API_NAME"))
}

/* ================= Helpers ================= */

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("[ENV] using system environment variables")
	}
}

func initDatastore() *datastore.Datastore {
	ds := datastore.NewDatastore(os.Getenv("BADGER_DATA_SOURCE"))
	log.Println("[Datastore] initialized")
	return ds
}

func initAsynqClient() *asynq.Client {
	addr := os.Getenv("REDIS_HOST")
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: addr})
	log.Println("[Asynq] connected:", addr)
	return client
}

func startWorker(ds *datastore.Datastore) {
	go func() {
		log.Println("[Worker] started")
		worker.StartWorker(ds, os.Getenv("REDIS_HOST"))
	}()
}

func startAPIServer(
	ds *datastore.Datastore,
	queue *asynq.Client,
	lb *loadbalancer.LoadBalancer,
) {
	nodeService := services.NewNodeService(ds, queue, lb)
	handler := handlers.NewNodeHandler(nodeService)

	server := api.NewServer(handler)

	go func() {
		log.Println("[API] server started")
		server.Start()
	}()
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
