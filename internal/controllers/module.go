package controllers

import (
	app_http "backend/pkg/http"

	"go.uber.org/fx"
)

var Module = fx.Module("controllers",
	fx.Provide(
		fx.Annotate(NewAuthController, fx.As(new(app_http.Controller)), fx.ResultTags(`group:"controllers"`)),
		fx.Annotate(NewUserController, fx.As(new(app_http.Controller)), fx.ResultTags(`group:"controllers"`)),
	),
)
