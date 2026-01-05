package server

import (
	"context"
	"fmt"
	"log"
	"os/signal"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/database"
	"github.com/ruhulfbr/mini-k8s/internal/logger/slog"
)

func Run() error {
	// Load application configuration once
	cfg := config.Load()

	loggerCleanup, err := initLogger(cfg)
	if err != nil {
		return err
	}
	defer loggerCleanup()

	// Root context cancelled on SIGINT / SIGTERM
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer stop()

	// Initialize database (shared by all components)
	ds := database.NewDatastore(cfg)
	defer ds.Close()

	// Initialize Asynq client (used by API to enqueue jobs)
	asynqClient := asynq.NewClient(asynq.RedisClientOpt{
		Addr: cfg.Redis.Host,
	})
	defer asynqClient.Close()

	// Initialize and start load balancer (is a long-running component)
	lb := InitLoadBalancer(ds)

	// Start background worker (Asynq consumer). Runs independently of the HTTP server
	// go startWorker(cfg, ds)

	// Initialize HTTP API server
	app, err := InitServer(cfg, ds, asynqClient, lb)
	if err != nil {
		return err
	}

	// Start HTTP server in background
	go func() {
		log.Println("[API] started")
		if err := StartServer(app, cfg); err != nil {
			log.Println("[API] stopped:", err)
			stop()
		}
	}()

	log.Println("[System] started")

	// Block until shutdown signal is received
	<-ctx.Done()

	log.Println("[System] shutting down...")

	// NOTE:
	// - Database and Asynq client are closed via defer
	// - Worker and load balancer should stop when process exits
	// - For graceful shutdown, pass ctx to those components

	return nil
}

func initLogger(cfg *config.Config) (func(), error) {
	if err := slog.Init(cfg.Logger); err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}
	return func() {}, nil
}
