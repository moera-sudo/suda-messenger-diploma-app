package main

import (
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/media-service/config"
	"github.com/moera-sudo/backend/backend/media-service/internal/pkg/logger"
	"github.com/moera-sudo/backend/backend/media-service/internal/app"
)

// @title Suda Media Service API
// @version 1.0
// @description REST API for file upload, download and management via S3 presigned URLs.
// @host localhost:8084
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger.Setup()
	cfg := config.Load()

	application := app.NewApp(cfg)

	if err := application.Run(); err != nil{
		log.Fatal().Err(err).Msg("Failed to start Media-service")
	}
}