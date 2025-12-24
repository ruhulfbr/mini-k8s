package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hibiken/asynq"
	"github.com/joho/godotenv"

	"github.com/ruhulfbr/mini-k8s/api"
	"github.com/ruhulfbr/mini-k8s/loadbalancer"
	"github.com/ruhulfbr/mini-k8s/orchestrator"
	"github.com/ruhulfbr/mini-k8s/worker"
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
	// Load Cluster Configuration
	// ----------------------------------------------------
	orchestrator.LoadConfig("cluster.json")

	// ----------------------------------------------------
	// Start Background Services
	// ----------------------------------------------------
	startWorker(ds)
	startLoadBalancer(ds)

	// ----------------------------------------------------
	// Start API Server
	// ----------------------------------------------------
	lbMgr := loadbalancer.NewManager(ds)
	ctrl := orchestrator.NewController(ds, asynqClient, lbMgr)
	apiServer := api.NewServer(ctrl)

	go func() {
		apiPort := os.Getenv("API_PORT")

		log.Println("[API] listening on :", apiPort)
		if err := apiServer.Start(":" + apiPort); err != nil {
			log.Fatal(err)
		}
	}()

	// ----------------------------------------------------
	// Graceful Shutdown
	// ----------------------------------------------------
	waitForShutdown()

	log.Println("[Main] shutting down MiniKube-Go...")
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

func initDatastore() *orchestrator.Datastore {
	ds := orchestrator.NewDatastore(os.Getenv("BADGER_DATA_SOURCE"))
	log.Println("[Datastore] initialized")

	return ds
}

func initAsynqClient() *asynq.Client {
	redisAddr := os.Getenv("REDIS_HOST")
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: redisAddr})
	log.Println("[Asynq] client connected:", redisAddr)
	return client
}

func startWorker(ds *orchestrator.Datastore) {
	go func() {
		log.Println("[Worker] starting...")
		worker.StartWorker(ds, os.Getenv("REDIS_HOST"))
	}()
}

func startLoadBalancer(ds *orchestrator.Datastore) {
	go func() {
		lb := loadbalancer.NewLoadBalancer(ds)
		lb.Serve()
	}()
}

func waitForShutdown() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
}
