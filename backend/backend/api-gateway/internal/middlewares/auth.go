package middleware

import (
	"net/http"
	"strings"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Мидлвейр авторизации для автоматической расшифровки токенов всех запросов проходящих через gateway
func AuthMiddleware(secret string, gatewaySecret string) echo.MiddlewareFunc {
	config := echojwt.Config{
		SigningKey: []byte(secret),

		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.MapClaims)
		},


		TokenLookupFuncs: []middleware.ValuesExtractor{
			func(c echo.Context) ([]string, error) {
				auth := c.Request().Header.Get("Authorization")
				if len(auth) > 7 && strings.EqualFold(auth[:7], "Bearer ") {
					return []string{auth[7:]}, nil
				}
				return nil, nil // Не нашли, идем дальше
			},
			func(c echo.Context) ([]string, error) {
				token := c.QueryParam("token")
				if token != "" {
					return []string{token}, nil
				}
				return nil, nil
			},
		},
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path

			if strings.HasPrefix(path, "/api/v1/messenger/auth/") {
				return true
			}


			if path == "/ws" || strings.HasPrefix(path, "/ws/") {
				return true
			}

			if strings.HasPrefix(path, "/health") {
				return true
			}

			if strings.HasPrefix(path, "/health") {
				return true
			}

			// Swagger UI and service specs
			if strings.HasPrefix(path, "/swagger/") {
				return true
			}
			if strings.Contains(path, "/swagger/") {
				return true
			}

			return false
		},
		

		SuccessHandler: func(c echo.Context) {
			token := c.Get("user").(*jwt.Token)
			claims := token.Claims.(*jwt.MapClaims)

			tokenType, ok := (*claims)["type"].(string)
			if !ok || tokenType != "access" {
				c.Error(echo.NewHTTPError(http.StatusUnauthorized, "Invalid token type"))
			}

			var userID string
			if val, ok := (*claims)["user_id"].(string); ok {
				userID = val
			} else if val, ok := (*claims)["sub"].(string); ok {
				userID = val
			}

			if userID == "" {
				c.Error(echo.NewHTTPError(http.StatusUnauthorized, "Missing user ID"))
				return 
			}

			role := "user"
			if val, ok := (*claims)["role"].(string); ok{
				role = val
			}

			timestamp := strconv.FormatInt(time.Now().Unix(), 10)
			path := c.Request().URL.Path

			c.Request().Header.Set("X-User-ID", userID)
			c.Request().Header.Set("X-User-Role", role)
			c.Request().Header.Set("X-Gateway-Timestamp", timestamp)
			c.Request().Header.Set("X-Gateway-Path", path)

			// Forward the session id ("sid") so downstream services can mark the
			// current session in the Active Sessions list. Absent for stateless
			// tokens (e.g. webapp-session) — then the header is simply empty.
			if sid, ok := (*claims)["sid"].(string); ok && sid != "" {
				c.Request().Header.Set("X-Session-ID", sid)
			} else {
				c.Request().Header.Del("X-Session-ID")
			}

			message := fmt.Sprintf("%s:%s:%s:%s", userID, role, timestamp, path)
			mac := hmac.New(sha256.New, []byte(gatewaySecret))
			mac.Write([]byte(message))
			signature := hex.EncodeToString(mac.Sum(nil))

			c.Request().Header.Set("X-Gateway-Signature", signature)
		},
	}
	return echojwt.WithConfig(config)
}
