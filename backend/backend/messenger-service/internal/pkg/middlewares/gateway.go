package middlewares

import (
	"github.com/labstack/echo/v4"
)

// Мидлвейр ложит секретный ключ из конфига в контекст чтобы getuserid мог расшифровать его
// Это необходимо для безопасности. Чтобы валидировать что запрос пришел от Гейтвея, который подписал запрос с тем же ключом


func GatewaySecret(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("gateway_secret", secret)
			return next(c)
		}
	}
}