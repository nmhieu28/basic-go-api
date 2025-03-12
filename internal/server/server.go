package server

import (
	configs "backend/pkg/config"
	"context"
	"errors"
	"net/http"
	"time"

	"backend/pkg/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

const (
	certFile       = "ssl/Server.crt"
	keyFile        = "ssl/Server.pem"
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

func BootStrapServer(ctx context.Context, e *echo.Echo, config *configs.AppConfig, logger logger.Logger) error {
	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Infof("shutting down Http PORT: {%s}", config.Server.Port)
				err := e.Shutdown(ctx)
				if err != nil {
					logger.Errorf("(Shutdown) err: {%v}", err)
					return
				}
				logger.Info("server exited properly")
				return
			}
		}
	}()
	err := e.Start(config.Server.Port)
	return err
}
func Run(lc fx.Lifecycle, ctx context.Context, e *echo.Echo, config *configs.AppConfig, logger logger.Logger) error {

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			e.Server.ReadTimeout = config.Server.ReadTimeout * time.Second
			e.Server.WriteTimeout = config.Server.WriteTimeout * time.Second
			e.Server.MaxHeaderBytes = maxHeaderBytes
			go func() {
				if err := BootStrapServer(ctx, e, config, logger); !errors.Is(err, http.ErrServerClosed) {
					logger.Fatalf("error running http server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			logger.Infof("all servers shutdown gracefully...")
			return nil
		},
	})

	return nil
}
