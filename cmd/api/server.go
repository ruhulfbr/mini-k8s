package api

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/ruhulfbr/mini-k8s/internal/http/handlers"
)

type Server struct {
	ctrl *handlers.NodeHandler
}

func NewServer(ctrl *handlers.NodeHandler) *Server {
	return &Server{ctrl: ctrl}
}

func (s *Server) Start() {
	e := echo.New()

	// Api Routs
	api := e.Group("/api")
	api.POST("/deploy", s.ctrl.HandleDeploy)
	api.POST("/scale", s.ctrl.HandleScale)

	e.Logger.Fatal(e.Start(":" + os.Getenv("API_PORT")))
}
