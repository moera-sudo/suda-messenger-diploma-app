package middlewares

import (
	"github.com/labstack/echo/v4"
)

// GatewaySecret кладёт секретный ключ из конфига в контекст, чтобы GetUserId
// мог его прочитать и валидировать HMAC-подпись запроса от api-gateway.
//
// Это необходимо для безопасности: gateway подписывает каждый запрос своим
// секретом (HMAC-SHA256 от user_id + role + timestamp + path), и сервис
// проверяет подпись тем же секретом. Если злоумышленник попробует пойти
// напрямую в transaction-service, минуя gateway — без знания секрета он
// не сможет подделать X-Gateway-Signature, и запрос будет отклонён.
func GatewaySecret(secret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("gateway_secret", secret)
			return next(c)
		}
	}
}