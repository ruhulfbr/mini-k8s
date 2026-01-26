package routes

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/http/handlers"
	"github.com/ruhulfbr/mini-k8s/internal/http/middleware"
	"github.com/ruhulfbr/mini-k8s/internal/http/web"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/database"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/logger/slog"
)

func ConfigureRoutes(
	engine *echo.Echo,
	ds *database.Database,
	asynqClient *asynq.Client,
) {
	engine.Renderer = web.NewRenderer()

	appHandlers := handlers.InitHandlers(ds, asynqClient)

	engine.HTTPErrorHandler = middleware.NewEchoHTTPErrorHandler()

	if config.GetLoggerConfig().EnableRequestLog {
		tracer := slog.NewTraceStarter(uuid.NewV7)
		engine.Use(middleware.NewRequestLogger(tracer))
	}

	engine.GET("/", func(c echo.Context) error {
		return c.Render(http.StatusOK, "index.html", map[string]any{
			"title":    "Applications",
			"pageName": "dashboard",
		})
	})

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

	cluster.POST("/:id/deploy", appHandlers.ClusterHandler.Deploy)
	cluster.POST("/:id/rolling-deploy", appHandlers.ClusterHandler.RollingDeploy)
	cluster.POST("/:id/scale", appHandlers.ClusterHandler.HandleScale)

	//pods := api.Group("/clusters/:id/pods")
	//pods.GET("", appHandlers.PodHandler.ListByCluster)
	//pods.POST("", appHandlers.PodHandler.Create)
	//pods.DELETE("/:id", appHandlers.PodHandler.Delete)
}
