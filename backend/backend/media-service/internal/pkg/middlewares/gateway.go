package middlewares

import (
	"github.com/labstack/echo/v4"
)

// GatewaySecret puts the gateway signature secret into the request context so that
// utils.GetUserId can validate the X-Gateway-Signature header. The api-gateway signs
// every proxied request with the same key; without this middleware GetUserId fails
// at the "secret not configured" check and every authorized endpoint returns 401.
func GatewaySecret(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("gateway_secret", secret)
			return next(c)
		}
	}
}
