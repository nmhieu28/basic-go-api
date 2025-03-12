package server

import (
	"backend/pkg/http"
	"backend/pkg/logger"

	"github.com/labstack/echo/v4"
	"go.uber.org/fx"
)

func MapHandlers(e *echo.Echo, logger logger.Logger, controllers []http.Controller) {
	apiV1 := e.Group("api")

	for _, ctrl := range controllers {
		ctrl.RegisterRoute(apiV1)
	}
}

var Module = fx.Module("handlers",
	fx.Invoke(
		fx.Annotate(
			MapHandlers,
			fx.ParamTags("", "", `group:"controllers"`),
		),
	),
)
