package controllers

import (
	"backend/internal/models/requests"
	auth_service "backend/internal/services"
	app_errors "backend/pkg/errors"
	"backend/pkg/logger"
	"backend/pkg/response"
	"net/http"

	app_http "backend/pkg/http"

	"github.com/labstack/echo/v4"
)

type AuthController struct {
	auth_service *auth_service.IdentityService
	logger       logger.Logger
}

func NewAuthController(auth_service *auth_service.IdentityService, logger logger.Logger) app_http.Controller {
	return &AuthController{auth_service: auth_service, logger: logger}
}

func (c *AuthController) RegisterRoute(r *echo.Group) {
	r.POST("/accounts/register", c.Register)
	r.POST("/accounts/verify-email", c.VerifyEmail)
	r.POST("/accounts/login", c.Login)
}

func (c *AuthController) Register(ctx echo.Context) error {
	var request requests.CreateUserRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Failure(app_errors.NewGeneralError(app_errors.DataInvalid)))
	}

	result, err := c.auth_service.Register(ctx.Request().Context(), request)

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Failure(err.(app_errors.AppError)))
	}

	return ctx.JSON(http.StatusOK, response.Success(result))
}

func (c *AuthController) Login(ctx echo.Context) error {
	var request requests.LoginRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Failure(app_errors.NewGeneralError(app_errors.DataInvalid)))
	}
	result := c.auth_service.Login(ctx.Request().Context(), request)
	if !result.IsSuccess {
		return ctx.JSON(http.StatusBadRequest, result)
	}
	ctx.SetCookie(&http.Cookie{
		Name:     "refreshToken",
		Value:    result.Data.RefreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})
	return ctx.JSON(http.StatusOK, result)
}

func (c *AuthController) VerifyEmail(ctx echo.Context) error {
	var verifyRequest requests.VerifyEmailRequest
	if err := ctx.Bind(&verifyRequest); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Failure(app_errors.NewGeneralError(app_errors.DataInvalid)))
	}
	result := c.auth_service.VerifyEmail(ctx.Request().Context(), verifyRequest.Token)
	if !result.IsSuccess {
		return ctx.JSON(http.StatusBadRequest, result)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *AuthController) ResetPassword(ctx echo.Context) error {
	var request requests.ResetPasswordRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Failure(app_errors.NewGeneralError(app_errors.DataInvalid)))
	}
	result := c.auth_service.ResetPassword(ctx.Request().Context(), request)
	if !result.IsSuccess {
		return ctx.JSON(http.StatusBadRequest, result)
	}
	return ctx.JSON(http.StatusOK, result)
}
func (c *AuthController) ForgotPassword(ctx echo.Context) error {
	var request requests.ForgotPasswordRequest
	if err := ctx.Bind(&request); err != nil {
		return ctx.JSON(http.StatusBadRequest, response.Failure(app_errors.NewGeneralError(app_errors.DataInvalid)))
	}
	result := c.auth_service.ForgotPassword(ctx.Request().Context(), request)
	if !result.IsSuccess {
		return ctx.JSON(http.StatusBadRequest, result)
	}
	return ctx.JSON(http.StatusOK, result)
}
