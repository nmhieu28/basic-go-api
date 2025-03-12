package http

import (
	"backend/pkg/middlewares"

	"github.com/labstack/echo/v4"
)

type Controller interface {
	RegisterRoute(*echo.Group)
}
type BaseController struct{}

func NewBaseController() *BaseController {
	return &BaseController{}
}
func (c *BaseController) CurrentUser(ctx echo.Context) middlewares.CurrentUser {
	return ctx.Get("currentUser").(middlewares.CurrentUser)
}
