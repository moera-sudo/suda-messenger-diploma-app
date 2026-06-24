package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/api-gateway/config"
	gm "github.com/moera-sudo/backend/backend/api-gateway/internal/middlewares"
    gatewaySwagger "github.com/moera-sudo/backend/backend/api-gateway/internal/swagger"

	"github.com/moera-sudo/backend/backend/api-gateway/internal/logger"
	"github.com/moera-sudo/backend/backend/api-gateway/internal/proxy"
)

func main() {
	cfg := config.Load()
	e := echo.New()
	logger.Setup()

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		Skipper: func(c echo.Context) bool {
			return false
		},
		LogURI:      true,
		LogStatus:   true,
		LogMethod:   true,
		LogLatency:  true,
		LogRemoteIP: true,
		LogError:    true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error != nil {
				log.Error().
					Err(v.Error).
					Str("method", v.Method).
					Str("uri", v.URI).
					Int("status", v.Status).
					Dur("latency", v.Latency).
					Str("remote_ip", v.RemoteIP).
					Msg("GATEWAY_ERROR")
				return nil
			}
			log.Info().
				Str("method", v.Method).
				Str("uri", v.URI).
				Int("status", v.Status).
				Dur("latency", v.Latency).
				Str("remote_ip", v.RemoteIP).
				Msg("GATEWAY")
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PATCH, echo.POST, echo.DELETE, echo.PUT, echo.OPTIONS, echo.CONNECT},
	}))

	e.Use(gm.GatewayLoggerMiddleware())

	// Forward the client-facing host so downstream services (media) can sign
	// presigned URLs for the same origin the client used to reach the API.
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get("X-Forwarded-Host") == "" {
				c.Request().Header.Set("X-Forwarded-Host", c.Request().Host)
			}
			return next(c)
		}
	})

	e.Use(gm.AuthMiddleware(cfg.JWTSecret, cfg.GatewaySignatureSecret))

	swaggerAgg := gatewaySwagger.NewAggregator([]gatewaySwagger.ServiceSpec{
		{
			Name:   "messenger",
			URL:    cfg.MessengerHost + "/swagger/doc.json",
			Prefix: "/api/v1/messenger",
		},
		{
			Name:   "media",
			URL:    cfg.MediaHost + "/swagger/doc.json",
			Prefix: "/api/v1/media",
		},
		// Добавишь позже:
		// {Name: "transactions", URL: txHost + "/swagger/doc.json", Prefix: "/api/v1/transactions"},
		// {Name: "applications", URL: appHost + "/swagger/doc.json", Prefix: "/api/v1/apps"},
	})
	swaggerAgg.RegisterRoutes(e)
	// Реализация хэндлера для /health
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	proxy.SetupRoutes(e, cfg)

	log.Info().Msg(" Gateway started on %s\n" + cfg.Port)

	if err := e.Start(cfg.Port); err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("failed to start server")
	}
}
