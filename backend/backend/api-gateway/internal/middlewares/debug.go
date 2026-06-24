package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

func GatewayLoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			path := c.Request().URL.Path

			// Публичные маршруты — не логируем auth warnings
			if isPublicRoute(path) {
				return next(c)
			}

			authHeader := c.Request().Header.Get("Authorization")
			logger := log.With().
				Str("method", c.Request().Method).
				Str("path", path).
				Str("remote_ip", c.RealIP()).
				Logger()

			if authHeader == "" {
				logger.Warn().Msg("No Authorization header")
			} else if len(authHeader) > 7 && strings.EqualFold(authHeader[:7], "Bearer ") {
				token := authHeader[7:]
				masked := token[:4] + "..." + token[len(token)-4:]
				logger.Debug().
					Str("token_preview", masked).
					Int("token_length", len(token)).
					Msg("Bearer token received")
			} else {
				logger.Warn().
					Str("auth_header", authHeader[:min(20, len(authHeader))]).
					Msg("Malformed Authorization header")
			}

			return next(c)
		}
	}
}

func isPublicRoute(path string) bool {
	return strings.HasPrefix(path, "/health") ||
		strings.Contains(path, "/swagger/") ||
		strings.HasPrefix(path, "/api/v1/messenger/auth/")
}