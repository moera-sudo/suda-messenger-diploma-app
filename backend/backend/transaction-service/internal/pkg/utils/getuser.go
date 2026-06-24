package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// maxTimestampDrift — максимально допустимое расхождение между временем
// в запросе и текущим (в секундах). Защищает от replay-атак.
const maxTimestampDrift = 30

// GetUserId извлекает user_id из заголовков и проверяет HMAC-подпись,
// которую проставил api-gateway.
//
// Заголовки которые шлёт gateway:
//
//	X-User-ID            — UUID юзера
//	X-User-Role          — USER / ADMIN / SUPERADMIN
//	X-Gateway-Timestamp  — Unix timestamp когда gateway подписал запрос
//	X-Gateway-Path       — оригинальный путь до прокси (если отличается)
//	X-Gateway-Signature  — HMAC-SHA256(secret, "{user_id}:{role}:{ts}:{path}")
//
// Если подпись невалидна, timestamp устарел или заголовков нет — возвращает
// echo.HTTPError с подходящим кодом.
func GetUserId(c echo.Context) (uuid.UUID, error) {

	if os.Getenv("DEV_SKIP_GATEWAY") == "true" {
		idStr := c.Request().Header.Get("X-User-ID")
		if idStr == "" {
			return uuid.Nil, echo.NewHTTPError(401, "X-User-ID header required (dev mode)")
		}
		uid, err := uuid.Parse(idStr)
		if err != nil {
			return uuid.Nil, echo.NewHTTPError(401, "Invalid User ID format")
		}
		return uid, nil
	}

	// Секрет кладёт middleware GatewaySecret в начале цепочки.
	secret, ok := c.Get("gateway_secret").(string)
	if !ok || secret == "" {
		return uuid.Nil, echo.NewHTTPError(500, "Gateway secret not configured")
	}

	idStr := c.Request().Header.Get("X-User-ID")
	role := c.Request().Header.Get("X-User-Role")
	timestamp := c.Request().Header.Get("X-Gateway-Timestamp")
	signature := c.Request().Header.Get("X-Gateway-Signature")

	if idStr == "" || timestamp == "" || signature == "" {
		return uuid.Nil, echo.NewHTTPError(401, "Missing gateway headers")
	}

	ts, err := strconv.ParseInt(timestamp, 10, 64)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(400, "Invalid Timestamp")
	}
	if math.Abs(float64(time.Now().Unix()-ts)) > maxTimestampDrift {
		return uuid.Nil, echo.NewHTTPError(403, "Request expired")
	}

	// Если gateway проксирует, оригинальный путь может отличаться от
	// того, который видит наш сервис. Gateway может передать оригинал
	// через X-Gateway-Path, иначе используем то что видим у себя.
	path := c.Request().Header.Get("X-Gateway-Path")
	if path == "" {
		path = c.Request().URL.Path
	}

	message := fmt.Sprintf("%s:%s:%s:%s", idStr, role, timestamp, path)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	expected := hex.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return uuid.Nil, echo.NewHTTPError(401, "Invalid gateway signature")
	}

	uid, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, echo.NewHTTPError(401, "Invalid User ID format")
	}

	return uid, nil
}

// GetUserRole возвращает роль юзера, проставленную gateway'ем.
// Не проверяет подпись — для этого используется GetUserId. Зови сначала
// GetUserId (для проверки), потом GetUserRole (для использования).
func GetUserRole(c echo.Context) string {
	return c.Request().Header.Get("X-User-Role")
}

// GetClientIP возвращает IP клиента (для audit-лога).
// Берём X-Forwarded-For (его проставляет gateway/nginx), иначе RemoteAddr.
func GetClientIP(c echo.Context) string {
	if ip := c.Request().Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	return c.Request().RemoteAddr
}

// GetUserAgent возвращает User-Agent (для audit-лога).
func GetUserAgent(c echo.Context) string {
	return c.Request().Header.Get("User-Agent")
}
