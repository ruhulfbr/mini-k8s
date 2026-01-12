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
	ds *database.Database,
	asynqClient *asynq.Client,
	lb *loadBalancer.LoadBalancer,
) {
	appHandlers := handlers.InitHandlers(ds, asynqClient, lb)

	engine.HTTPErrorHandler = middleware.NewEchoHTTPErrorHandler()

	if config.GetLoggerConfig().EnableRequestLog {
		tracer := slog.NewTraceStarter(uuid.NewV7)
		engine.Use(middleware.NewRequestLogger(tracer))
	}

	api := engine.Group("/api")

	application := api.Group("/applications")
	application.GET("", appHandlers.ApplicationHandler.List)
	application.POST("", appHandlers.ApplicationHandler.Create)
	application.GET("/:id", appHandlers.ApplicationHandler.Show)
	application.PUT("/:id", appHandlers.ApplicationHandler.Update)
	application.DELETE("/:id", appHandlers.ApplicationHandler.Delete)

	cluster := api.Group("/applications/:appId/clusters")
	cluster.GET("", appHandlers.ClusterHandler.ListByApplication)
	cluster.POST("", appHandlers.ClusterHandler.Create)
	cluster.GET("/:id", appHandlers.ClusterHandler.Show)
	cluster.PUT("/:id", appHandlers.ClusterHandler.Update)
	cluster.DELETE("/:id", appHandlers.ClusterHandler.Delete)
	
	cluster.GET("/:id/builds", appHandlers.ClusterHandler.BuildHistory)
	cluster.POST("/:id/build-image", appHandlers.ClusterHandler.BuildImage)
	cluster.POST("/:id/pull-image", appHandlers.ClusterHandler.PullImage)

	pods := api.Group("/clusters/:id/pods")
	pods.GET("", appHandlers.PodHandler.ListByCluster)
	pods.POST("", appHandlers.PodHandler.Create)
	pods.DELETE("/:id", appHandlers.PodHandler.Delete)

	api.POST("/deploy", appHandlers.NodeHandler.HandleDeploy)
	api.POST("/scale", appHandlers.NodeHandler.HandleScale)
}
