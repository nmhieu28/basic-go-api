package server

import (
	configs "backend/pkg/config"
	"backend/pkg/constants"
	app_middlewares "backend/pkg/middlewares"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ConfigMiddlewares(e *echo.Echo, appConfig *configs.AppConfig) {

	if appConfig.Cors.Enable {
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: appConfig.Cors.Allows,
			AllowHeaders: []string{
				echo.HeaderOrigin,
				echo.HeaderContentType,
				echo.HeaderAccept,
				echo.HeaderXRequestID,
				echo.HeaderXCSRFToken,
				echo.HeaderAuthorization,
			},
			AllowCredentials: true,
		}))
	}

	e.HideBanner = false
	e.Use(middleware.Logger())
	e.Use(app_middlewares.CorrelationIdMiddleware)
	e.Use(middleware.RequestID())
	e.Use(middleware.BodyLimit(constants.BodyLimit))

}
