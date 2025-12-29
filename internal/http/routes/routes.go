package routes

import (
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/datastore"
	"github.com/ruhulfbr/mini-k8s/internal/http/handlers"
	"github.com/ruhulfbr/mini-k8s/internal/http/middleware"
	"github.com/ruhulfbr/mini-k8s/internal/loadbalancer"
)

func ConfigureRoutes(
	engine *echo.Echo,
	cfg *config.Config,
	ds *datastore.Datastore,
	asynqClient *asynq.Client,
	lb *loadbalancer.LoadBalancer,
) error {
	appHandlers := handlers.InitHandlers(cfg, ds, asynqClient, lb)

	engine.HTTPErrorHandler = middleware.EchoHTTPErrorHandler

	api := engine.Group("/api")
	api.POST("/deploy", appHandlers.NodeHandler.HandleDeploy)
	api.POST("/scale", appHandlers.NodeHandler.HandleScale)

	return nil
}
