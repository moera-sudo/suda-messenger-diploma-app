package main

import (
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/messenger-service/config"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/app"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/logger"
)

// @title Suda Messenger Service API
// @version 1.0
// @description REST API for Suda Messenger: authentication, chats, messaging, users, search.
// @host localhost:8081
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger.Setup()
	
	cfg := config.Load()

	application := app.NewApp(cfg)

	if err := application.Run(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start Messenger-service")
	}

}