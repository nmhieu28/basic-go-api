package middlewares

import (
	"context"
	"encoding/base64"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CtxKey string

const correlationIdHeader = echo.HeaderXCorrelationID

func generateShortUUID() string {
	u := uuid.New()
	uuidBytes, err := u.MarshalBinary()
	if err != nil {
		panic(err)
	}
	base64Uuid := base64.RawURLEncoding.EncodeToString(uuidBytes)
	return base64Uuid
}

func CorrelationIdMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		correlationId := c.Request().Header.Get(correlationIdHeader)
		if correlationId == "" {
			correlationId = generateShortUUID()
		}
		key := CtxKey(correlationIdHeader)
		newCtx := context.WithValue(c.Request().Context(), key, correlationId)
		c.SetRequest(c.Request().WithContext(newCtx))
		return next(c)
	}
}
