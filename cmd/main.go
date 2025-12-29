package main

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/cmd/api"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/datastore"
	"github.com/ruhulfbr/mini-k8s/internal/loadbalancer"
	"github.com/ruhulfbr/mini-k8s/internal/worker"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	ds := initDatastore(cfg)
	defer ds.Close()

	asynqClient := initAsynqClient(cfg)
	defer asynqClient.Close()

	var wg sync.WaitGroup

	lb := loadbalancer.NewLoadBalancer(ds)
	wg.Add(1)
	go runLoadBalancer(ctx, &wg, lb)

	wg.Add(1)
	go runWorker(ctx, &wg, ds, cfg)

	app, err := api.InitServer(cfg, ds, asynqClient, lb)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := api.StartServer(app, cfg); err != nil {
			log.Println("[API] stopped:", err)
			stop()
		}
	}()

	log.Println("[System] started")

	<-ctx.Done()
	log.Println("[System] shutting down...")

	wg.Wait()
	log.Println("[System] shutdown complete")
}

/* ================= Setup ================= */

func initDatastore(cfg *config.Config) *datastore.Datastore {
	ds := datastore.NewDatastore(cfg.Badger.DataSource)
	log.Println("[Datastore] initialized")
	return ds
}

func initAsynqClient(cfg *config.Config) *asynq.Client {
	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: cfg.Redis.Host,
	})
	log.Println("[Asynq] connected:", cfg.Redis.Host)
	return client
}

/* ================= Runners ================= */

func runLoadBalancer(ctx context.Context, wg *sync.WaitGroup, lb *loadbalancer.LoadBalancer) {
	defer wg.Done()
	log.Println("[LoadBalancer] started")

	go lb.Start()

	<-ctx.Done()
	log.Println("[LoadBalancer] stopped")
}

func runWorker(ctx context.Context, wg *sync.WaitGroup, ds *datastore.Datastore, cfg *config.Config) {
	defer wg.Done()
	log.Println("[Worker] started")

	redisWorker := worker.NewWorker(cfg, ds)
	go redisWorker.StartWorker()

	<-ctx.Done()
	log.Println("[Worker] stopped")
}
