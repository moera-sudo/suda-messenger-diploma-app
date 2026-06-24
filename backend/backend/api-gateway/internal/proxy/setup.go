package proxy

import (
	"net/url"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/moera-sudo/backend/backend/api-gateway/config"
)

func SetupRoutes(e *echo.Echo, cfg *config.Config) {

	messengerUrl, _ := url.Parse(cfg.MessengerHost)
	e.Group("/api/v1/messenger", middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{URL: messengerUrl},
		}),
		Rewrite: map[string]string{
			"/api/v1/messenger/*": "/$1",
		},
	}))

	txUrl, _ := url.Parse(cfg.TxHost)
	e.Group("/api/v1/tx", middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{URL: txUrl},
		}),
		Rewrite: map[string]string{

			"/api/v1/tx/*": "/api/v1/$1",
		},
	}))

	mediaUrl, _ := url.Parse(cfg.MediaHost)
	e.Group("/api/v1/media", middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{URL: mediaUrl},
		}),
		Rewrite: map[string]string{
			"/api/v1/media/*": "/api/media/$1",
		},
	}))


	e.Group("/ws", middleware.ProxyWithConfig(middleware.ProxyConfig{
		Balancer: middleware.NewRoundRobinBalancer([]*middleware.ProxyTarget{
			{URL: messengerUrl},
		}),
	}))
}
