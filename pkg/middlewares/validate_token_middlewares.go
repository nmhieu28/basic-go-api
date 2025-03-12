package middlewares

import (
	"net/http"
	"strings"

	"backend/pkg/cache"
	configs "backend/pkg/config"
	"backend/pkg/jwt_generate"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CurrentUser struct {
	UserId uuid.UUID
	Email  string
}

func ValidateTokenMiddleware(appConfig *configs.AppConfig, redisCache cache.Cache) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			auth, err := BearerAuth(c.Request())
			if err != nil {
				return err
			}
			jwtGen := jwt_generate.NewJwtGenerate(c.Request().Context(), appConfig, redisCache)
			token, err := jwtGen.VerifyToken(auth, appConfig.Jwt.SecretKey)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
			}
			currentUser := CurrentUser{
				UserId: token.UserId,
				Email:  token.Email,
			}
			c.Set("currentUser", currentUser)
			return next(c)
		}
	}
}
func BearerAuth(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Authorization header is required")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", echo.NewHTTPError(http.StatusUnauthorized, "Authorization header format must be Bearer {token}")
	}
	return parts[1], nil
}
