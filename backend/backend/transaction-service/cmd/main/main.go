// Package main — entry point для transaction-service.
//
// Делает три вещи:
//  1. Загружает config из .env / env переменных
//  2. Инициализирует zerolog
//  3. Передаёт управление в app.Init/app.Run, который держит весь lifecycle
package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/rs/zerolog/log"

	_ "github.com/moera-sudo/backend/backend/docs/swagger/transaction"
	"github.com/moera-sudo/backend/backend/transaction-service/config"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/app"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/pkg/logger"
)

// @title Suda Transaction Service API
// @version 1.0
// @description REST API for Suda Transaction Service: custodial wallets, SUDA transfers, donations, token-gating, purchases, treasury statistics.
// @host localhost:8082
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	logger.Setup()
	log.Info().Msg("transaction-service starting")

	cfg := config.Load()

	// Контекст, который отменяется при SIGINT / SIGTERM.
	// app.Run слушает ctx.Done() и инициирует graceful shutdown.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	application, err := app.Init(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to initialize app")
	}

	if err := application.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("app exited with error")
	}

	log.Info().Msg("transaction-service stopped")
}