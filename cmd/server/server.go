package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/hibiken/asynq"
	"github.com/labstack/echo/v4"
	"github.com/ruhulfbr/mini-k8s/internal/config"
	"github.com/ruhulfbr/mini-k8s/internal/http/routes"
	"github.com/ruhulfbr/mini-k8s/internal/http/validator"
	"github.com/ruhulfbr/mini-k8s/internal/infrastructure/database"
	"github.com/ruhulfbr/mini-k8s/internal/loadBalancer"
)

const shutdownTimeout = 20 * time.Second

type Server struct {
	echo *echo.Echo
}

func NewServer(echo *echo.Echo) *Server {
	return &Server{echo: echo}
}

func (s *Server) Start(addr string) error {
	if err := s.echo.Start(":" + addr); err != nil && errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("start echo: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	if err := s.echo.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown echo: %w", err)
	}

	return nil
}

func InitServer(
	cfg *config.Config,
	ds *database.Database,
	asynqClient *asynq.Client,
	lb *loadBalancer.LoadBalancer,
) (*Server, error) {
	engine := echo.New()
	engine.Validator = validator.New()

	routes.ConfigureRoutes(engine, cfg, ds, asynqClient, lb)

	return NewServer(engine), nil
}

func StartServer(app *Server, cfg *config.Config) error {
	go func() {
		_ = app.Start(strconv.Itoa(cfg.HTTP.Port))
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	return app.Shutdown(ctx)
}
