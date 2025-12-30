package routes

import (
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/database"
	"github.com/ruhulfbr/mini-k8s/internal/http/handlers"
	"github.com/ruhulfbr/mini-k8s/internal/http/middleware"
	"github.com/ruhulfbr/mini-k8s/internal/loadbalancer"
)

func ConfigureRoutes(
	engine *echo.Echo,
	cfg *config.Config,
	ds *database.Database,
	asynqClient *asynq.Client,
	lb *loadbalancer.LoadBalancer,
) {
	appHandlers := handlers.InitHandlers(cfg, ds, asynqClient, lb)

	engine.HTTPErrorHandler = middleware.EchoHTTPErrorHandler

	api := engine.Group("/api")

	app := api.Group("/applications")
	app.GET("", appHandlers.ApplicationHandler.List)
	app.POST("", appHandlers.ApplicationHandler.Create)
	app.PUT("/:id", appHandlers.ApplicationHandler.Update)
	app.DELETE("/:id", appHandlers.ApplicationHandler.Delete)

	svc := api.Group("/applications/:appId/services")
	svc.GET("", appHandlers.ServiceHandler.ListByApplication)
	svc.POST("", appHandlers.ServiceHandler.Create)
	svc.DELETE("/:id", appHandlers.ServiceHandler.Delete)

	pods := api.Group("/services/:serviceId/pods")
	pods.GET("", appHandlers.PodHandler.ListByService)
	pods.POST("", appHandlers.PodHandler.Create)
	pods.DELETE("/:id", appHandlers.PodHandler.Delete)

	api.POST("/deploy", appHandlers.NodeHandler.HandleDeploy)
	api.POST("/scale", appHandlers.NodeHandler.HandleScale)
}
