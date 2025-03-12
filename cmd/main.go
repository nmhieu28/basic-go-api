package main

import (
	"backend/internal/controllers"
	"backend/internal/infrastructures/migrations"
	"backend/internal/infrastructures/repositories"
	"backend/internal/server"
	"backend/internal/services"
	"backend/pkg/cache"
	configs "backend/pkg/config"
	"backend/pkg/database"
	"backend/pkg/http"
	"backend/pkg/jwt_generate"
	"backend/pkg/logger"
	"backend/pkg/mailer"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func NewHttpSever() *echo.Echo {
	return echo.New()
}
func main() {
	fx.New(
		fx.Options(
			fx.Provide(
				configs.InitAppConfig,
				logger.NewLogger,
				http.NewContext,
				NewHttpSever,
				database.NewDatabase,
				mailer.NewSMTPMailer,
				cache.NewRedisClient,
				jwt_generate.NewJwtGenerate,
			),
			repositories.Module,
			services.Module,
			controllers.Module,
			server.Module,
			fx.Invoke(
				server.Run,
				server.ConfigMiddlewares,
				func(dbEngine database.DBEngine) error {
					return migrations.Migrate(dbEngine)
				},
			),
		),
	).Run()
}
