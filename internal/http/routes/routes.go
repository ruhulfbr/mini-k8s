package routes

import (
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/http/handlers"
	"github.com/ruhulfbr/mini-k8s/internal/http/middleware"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/database"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger/slog"
	"github.com/ruhulfbr/mini-k8s/internal/loadBalancer"
)

func ConfigureRoutes(
	engine *echo.Echo,
	cfg *config.Config,
	ds *database.Database,
	asynqClient *asynq.Client,
	lb *loadBalancer.LoadBalancer,
) {
	appHandlers := handlers.InitHandlers(cfg, ds, asynqClient, lb)

	engine.HTTPErrorHandler = middleware.NewEchoHTTPErrorHandler(cfg)

	if cfg.Logger.EnableRequestLog {
		tracer := slog.NewTraceStarter(uuid.NewV7)
		engine.Use(middleware.NewRequestLogger(tracer))
	}

	api := engine.Group("/api")

	app := api.Group("/applications")
	app.GET("", appHandlers.ApplicationHandler.List)
	app.POST("", appHandlers.ApplicationHandler.Create)
	app.GET("/:id", appHandlers.ApplicationHandler.Show)
	app.PUT("/:id", appHandlers.ApplicationHandler.Update)
	app.DELETE("/:id", appHandlers.ApplicationHandler.Delete)

	svc := api.Group("/applications/:appId/services")
	svc.GET("", appHandlers.ServiceHandler.ListByApplication)
	svc.POST("", appHandlers.ServiceHandler.Create)
	svc.GET("/:id", appHandlers.ServiceHandler.Show)
	svc.PUT("/:id", appHandlers.ServiceHandler.Update)
	svc.DELETE("/:id", appHandlers.ServiceHandler.Delete)

	pods := api.Group("/services/:serviceId/pods")
	pods.GET("", appHandlers.PodHandler.ListByService)
	pods.POST("", appHandlers.PodHandler.Create)
	pods.DELETE("/:id", appHandlers.PodHandler.Delete)

	api.POST("/deploy", appHandlers.NodeHandler.HandleDeploy)
	api.POST("/scale", appHandlers.NodeHandler.HandleScale)
}
