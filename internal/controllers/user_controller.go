package controllers

import (
	"net/http"

	"backend/internal/services"
	"backend/pkg/cache"
	configs "backend/pkg/config"
	app_http "backend/pkg/http"
	"backend/pkg/middlewares"

	"github.com/labstack/echo/v4"
)

type UserController struct {
	app_http.BaseController
	userService *services.UserService
	redisCache  cache.Cache
	appConfig   *configs.AppConfig
}

func NewUserController(userService *services.UserService,
	redisCache cache.Cache, appConfig *configs.AppConfig) app_http.Controller {
	return &UserController{userService: userService, redisCache: redisCache, appConfig: appConfig}
}
func (c *UserController) RegisterRoute(r *echo.Group) {
	r.GET("/users/me", c.Me, middlewares.ValidateTokenMiddleware(c.appConfig, c.redisCache))
}
func (c *UserController) Me(ctx echo.Context) error {
	id := c.CurrentUser(ctx).UserId
	result := c.userService.GetUser(id, ctx.Request().Context())
	if !result.IsSuccess {
		return ctx.JSON(http.StatusBadRequest, result)
	}
	return ctx.JSON(http.StatusOK, result)
}
