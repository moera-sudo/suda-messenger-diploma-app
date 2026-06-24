package app

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"

	_ "github.com/moera-sudo/backend/backend/docs/swagger/media"
	"github.com/moera-sudo/backend/backend/media-service/config"
	mediagrpc "github.com/moera-sudo/backend/backend/media-service/internal/delivery/grpc"
	handler "github.com/moera-sudo/backend/backend/media-service/internal/delivery/http"
	gw "github.com/moera-sudo/backend/backend/media-service/internal/pkg/middlewares"
	"github.com/moera-sudo/backend/backend/media-service/internal/platform/clients/messenger"
	pgclient "github.com/moera-sudo/backend/backend/media-service/internal/platform/postgres"
	s3storage "github.com/moera-sudo/backend/backend/media-service/internal/platform/storage/s3"
	repo "github.com/moera-sudo/backend/backend/media-service/internal/repository"
	"github.com/moera-sudo/backend/backend/media-service/internal/service"
	mediapb "github.com/moera-sudo/backend/backend/proto/media"
)

type App struct {
	echo *echo.Echo
	cfg  *config.Config
}

func NewApp(cfg *config.Config) *App {
	return &App{
		echo: echo.New(),
		cfg:  cfg,
	}
}

func (a *App) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgPool, err := pgclient.NewClient(ctx, a.cfg.MediaDatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Postgres")
	}
	defer pgPool.Close()
	log.Info().Msg("Postgres connected")

	fileStorage, err := s3storage.NewS3Storage(
		a.cfg.S3Endpoint,
		a.cfg.S3AccessKey,
		a.cfg.S3SecretKey,
		a.cfg.S3UseSSL,
		a.cfg.S3PublicEndpoint,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to S3 storage")
	}
	log.Info().Msg("S3 Storage initialized")

	var accessChecker messenger.AccessChecker

	accessChecker, err = messenger.NewAccessChecker(a.cfg.MessengerGRPCAddr)
	if err != nil {
		log.Warn().Err(err).
			Str("addr", a.cfg.MessengerGRPCAddr).
			Msg("Failed to connect to messenger-service. Using NoopAccessChecker (all access checks will pass)")
		accessChecker = messenger.NewNoopAccessChecker()
	} else {
		log.Info().
			Str("addr", a.cfg.MessengerGRPCAddr).
			Msg("Connected to messenger-service for access checks")
	}
	defer accessChecker.Close()

	a.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
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
					Msg("REQUEST_ERROR")
				return nil
			}

			log.Info().
				Str("method", v.Method).
				Str("uri", v.URI).
				Int("status", v.Status).
				Dur("latency", v.Latency).
				Str("remote_ip", v.RemoteIP).
				Msg("REQUEST")

			return nil
		},
	}))
	a.echo.Use(middleware.Recover())
	a.echo.Use(middleware.CORS())
	// Inject the gateway signature secret so utils.GetUserId can verify
	// X-Gateway-Signature on every proxied request (mirrors messenger-service).
	a.echo.Use(gw.GatewaySecret(a.cfg.GatewaySignatureSecret))

	a.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK", "service": "media"})
	})

	// Swagger UI (direct and proxied via gateway)
	a.echo.GET("/swagger/*", echoSwagger.EchoWrapHandler())
	a.echo.GET("/api/media/swagger/*", echoSwagger.EchoWrapHandler())

	mediaRepo := repo.NewMediaRepository(pgPool)

	mediaService := service.NewMediaService(mediaRepo, fileStorage, accessChecker, a.cfg)

	mediaHandler := handler.NewMediaHandler(mediaService)
	mediaHandler.RegisterRoutes(a.echo)

	// gRPC server — messenger-service calls LinkMediaToEntity here when a message
	// with an attachment is sent. Without it, media never gets an entity link and
	// the ACL denies every non-owner (403 "private media without entity links").
	go func() {
		lis, err := net.Listen("tcp", a.cfg.MediaGRPCPort)
		if err != nil {
			log.Fatal().Err(err).Str("port", a.cfg.MediaGRPCPort).Msg("Failed to listen gRPC")
		}
		grpcSrv := grpc.NewServer()
		mediapb.RegisterMediaServiceServer(grpcSrv, mediagrpc.NewServer(mediaService))
		log.Info().Str("port", a.cfg.MediaGRPCPort).Msg("Media gRPC server started")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	go func() {
		if err := a.echo.Start(a.cfg.MediaPort); err != nil {
			log.Info().Err(err).Msg("Server shutting down")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info().Msg("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := a.echo.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}