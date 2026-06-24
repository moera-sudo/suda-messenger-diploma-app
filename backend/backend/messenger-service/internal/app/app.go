package app

// Пакет для инициализации сервисов и хэндлеров, чтобы не делать этого в main.go

import (
	// Стандартная библиотека
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Внешние зависимости
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"google.golang.org/grpc"

	// Конфигурация и инфраструктура сервиса
	_ "github.com/moera-sudo/backend/backend/docs/swagger/messenger"
	"github.com/moera-sudo/backend/backend/messenger-service/config"
	appcrypto "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/crypto"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/notification"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/postgres"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/redis"

	// Middleware и валидация
	gw "github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/middlewares"
	customValidator "github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/validator"

	// gRPC: сервер и клиенты
	grpcServer "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/grpc"
	pb "github.com/moera-sudo/backend/backend/proto/messenger_access"
	pbTxAccess "github.com/moera-sudo/backend/backend/proto/transaction_access"

	// Email
	emailConfig "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/email/config"
	emailSender "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/email/sender"

	// WebSocket
	wsDelivery "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket/delivery"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket/hub"

	// Фича: Auth
	authConfig "github.com/moera-sudo/backend/backend/messenger-service/internal/features/auth/config"
	authHttp "github.com/moera-sudo/backend/backend/messenger-service/internal/features/auth/delivery/http"
	authRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/auth/repository"
	authService "github.com/moera-sudo/backend/backend/messenger-service/internal/features/auth/service"

	// Фича: User
	userHttp "github.com/moera-sudo/backend/backend/messenger-service/internal/features/user/delivery/http"
	userRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/user/repository"
	userService "github.com/moera-sudo/backend/backend/messenger-service/internal/features/user/service"

	// Фича: Chat
	chatHttp "github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/delivery/http"
	chatRepository "github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/repository"
	chatService "github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/service"

	// Фича: Search
	searchHttp "github.com/moera-sudo/backend/backend/messenger-service/internal/features/search/delivery/http"
	searchRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/search/repository"
	searchService "github.com/moera-sudo/backend/backend/messenger-service/internal/features/search/service"

	// Фича: Preferences
	prefsRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences/repository"
	prefsService "github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences/service"
	prefsHttp "github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences/delivery/http"

	// Фича: Channel
	channelRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel/repository"
	channelService "github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel/service"
	channelHttp "github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel/delivery/http"

	// Фича: Pins
	pinsRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/userpins/repository"
	pinsService "github.com/moera-sudo/backend/backend/messenger-service/internal/features/userpins/service"
	pinsHttp "github.com/moera-sudo/backend/backend/messenger-service/internal/features/userpins/delivery/http"

	// Фича: Friends
	friendsHttp "github.com/moera-sudo/backend/backend/messenger-service/internal/features/friends/delivery/http"
	friendsRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/friends/repository"
	friendsService "github.com/moera-sudo/backend/backend/messenger-service/internal/features/friends/service"
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

	// * Infrastructure. Connect to external services and platforms: Pgsql/Redis/Firebase/Email/Grpc
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pgPool, err := postgres.NewClient(ctx, a.cfg.DatabaseURL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Postgres")
	}
	defer pgPool.Close()
	log.Info().Msg("Postgres connected")

	redisClient, err := redis.NewClient(ctx, a.cfg.RedisAddr, a.cfg.RedisPass)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to Redis")
	}
	defer redisClient.Close()
	log.Info().Msg("Redis connected")

	emailCfg := emailConfig.Load()
	emailSenderInstance := emailSender.NewSender(emailCfg)

	fcmKeyPath := os.Getenv("FIREBASE_CREDENTIALS")
	log.Info().Msgf("Loading Firebase key from: %s", fcmKeyPath)
	notifyService, err := notification.NewFCMService(fcmKeyPath, pgPool)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to init FCM service (Push notifications disabled)")
	} else {
		log.Info().Msg("FCM Service initialized. Notifications enabled")
	}

	// * First of all. Connection to media service
	mediaClient, err := grpcServer.NewMediaClient(a.cfg.MediaServiceAddr)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to connect to Media Service (attachments disabled)")
	} else {
		log.Info().Msg("Media Service connected")
		defer mediaClient.Close()
	}

	// * Init chat service and repositories for shared functions, also init websocket's hub for chatservice
	// WS Hub. Pass JWT secret so it can authenticate direct WS connections that supply token via query
	// string (browsers/WebViews can't send custom headers on the WS handshake).
	wsHub := hub.NewHub(redisClient, authConfig.Load().JWTSecret)

	// Message content encryption at rest (AES-256-GCM). Empty key => no-op cipher
	// (messages stored as plaintext), so the service still runs without a key.
	contentCipher, err := appcrypto.NewContentCipher(a.cfg.MessageEncryptionKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Invalid MESSAGE_ENCRYPTION_KEY")
	}
	if appcrypto.Enabled(contentCipher) {
		log.Info().Msg("Message content encryption enabled (AES-256-GCM)")
	} else {
		log.Warn().Msg("Message content encryption DISABLED (no MESSAGE_ENCRYPTION_KEY)")
	}

	// Repositories
	memberRepo := chatRepository.NewMemberRepo(pgPool)
	chatRepo := chatRepository.NewChatRepo(pgPool, contentCipher)
	messageRepo := chatRepository.NewMessageRepo(pgPool, contentCipher)
	moderationRepo := chatRepository.NewModerationRepo(pgPool, contentCipher)
	prefRepo := prefsRepo.NewPreferencesRepo(pgPool)
	hiddenRepo := prefsRepo.NewLastSeenHiddenRepo(pgPool)
	userPinsRepository := pinsRepo.NewRepository(pgPool)
	
	// User Repository
	userRepository := userRepo.NewRepository(pgPool)


	// Transaction service gRPC client (must be created BEFORE chatSrv so it can be injected).
	txClient, err := grpcServer.NewTransactionClient(a.cfg.TransactionServiceAddr)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to create Transaction Service client (wallets will not work)")
		txClient = nil
	} else {
		log.Info().Str("addr", a.cfg.TransactionServiceAddr).Msg("Transaction Service client created")
		defer txClient.Close()
	}

	//Services
	prefsServiceInstance := prefsService.NewService(prefRepo, hiddenRepo, redisClient, wsHub)
	chatSrv := chatService.NewChatService(
		chatRepo, messageRepo, memberRepo, moderationRepo,
		prefsServiceInstance, userPinsRepository, wsHub, redisClient, notifyService, mediaClient,
		txClient,
	)
	// * Start messenger grpc server. Connection to other microservices
	accessServer := grpcServer.NewAccessServer(memberRepo)

	txAccessServer := grpcServer.NewTransactionAccessServer(userRepository, memberRepo, messageRepo, wsHub)

	grpcSrv := grpc.NewServer()
	pb.RegisterMessengerAccessServiceServer(grpcSrv, accessServer)
	pbTxAccess.RegisterTransactionAccessServiceServer(grpcSrv, txAccessServer)



	go func() {
		lis, err := net.Listen("tcp", a.cfg.GRPCPort)
		if err != nil {
			log.Fatal().Err(err).Str("port", a.cfg.GRPCPort).Msg("Failed to listen gRPC")
		}
		log.Info().Str("port", a.cfg.GRPCPort).Msg("gRPC server started")
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatal().Err(err).Msg("gRPC server failed")
		}
	}()

	// * Echo middlewares
	a.echo.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		// Пропускаем /ws — у WebSocket своё управление соединением
		Skipper: func(c echo.Context) bool {
			return c.Request().URL.Path == "/ws"
		},
		LogURI:           true,
		LogStatus:        true,
		LogMethod:        true,
		LogLatency:       true,
		LogRemoteIP:      true,
		LogUserAgent:     true,
		LogError:         true,
		LogHost:          true,
		LogContentLength: true,
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
	a.echo.Use(gw.GatewaySecret(a.cfg.GatewaySignatureSecret))
	a.echo.Validator = customValidator.New()

	// * Feature's repositories

	authRepository := authRepo.NewRepository(pgPool)
	searchRepository := searchRepo.NewRepository(pgPool, contentCipher)
	chRepo := channelRepo.NewChannelRepo(pgPool)
	commentsRepo := channelRepo.NewCommentsRepo(pgPool)
	appsRepo := channelRepo.NewAppsRepo(pgPool)
	joinReqRepo := channelRepo.NewJoinRequestRepo(pgPool)

	// * Feature's services

	authCfg := authConfig.Load()
	authServiceInstance := authService.NewService(authRepository, emailSenderInstance, chatSrv, authCfg, userRepository, txClient, a.cfg.InitDataSignToken)

	userServiceInstance := userService.NewService(userRepository, authRepository, wsHub, redisClient, notifyService, mediaClient, prefsServiceInstance, moderationRepo, a.cfg)

	searchServiceInstance := searchService.NewService(searchRepository, memberRepo)

	channelServiceInstance := channelService.NewService(
		chRepo, commentsRepo, appsRepo, joinReqRepo, chatRepo, memberRepo, userRepository, redisClient, wsHub, txClient,
	)

	pinsServiceInstance := pinsService.NewService(userPinsRepository, memberRepo)

	// Friends feature (Facebook-like requests).
	friendsRepository := friendsRepo.New(pgPool)
	friendsServiceInstance := friendsService.New(friendsRepository, userRepository, moderationRepo, wsHub)

	// * Feature's handlers and endpoints
	// Health Checker Handler
	a.echo.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "OK", "service": "messenger"})
	})

	// Swagger UI
	a.echo.GET("/swagger/*", echoSwagger.EchoWrapHandler())
	// Ws Router
	wsRouter := wsDelivery.NewRouter(chatSrv)
	wsHub.SetRouter(wsRouter.HandleIncomingMessage)

	// Online/Offline callbacks
	wsHub.SetOnConnect(func(userID uuid.UUID) {
		userServiceInstance.OnUserConnected(userID)
	})
	wsHub.SetOnDisconnect(func(userID uuid.UUID) {
		userServiceInstance.OnUserDisconnected(userID)
	})

	// Chat's Handlers
	chatHandler := chatHttp.NewChatHandler(chatSrv)
	membersHandler := chatHttp.NewMemberHandler(chatSrv)
	messageHandler := chatHttp.NewMessageHandler(chatSrv)
	// Other handlers
	authHandler := authHttp.NewAuthHandler(authServiceInstance)
	userHandler := userHttp.NewUserHandler(userServiceInstance)
	searchHandler := searchHttp.NewSearchHandler(searchServiceInstance)
	prefsHandler := prefsHttp.NewHandler(prefsServiceInstance)
	channelHandler := channelHttp.NewHandler(channelServiceInstance)
	pinsHandler := pinsHttp.NewHandler(pinsServiceInstance)
	friendsHandler := friendsHttp.NewHandler(friendsServiceInstance)

	a.echo.GET("/ws", func(c echo.Context) error {
		return wsHub.HandleRequest(c)
	})

	// * Register routes
	chatHandler.RegisterRoutes(a.echo)
	membersHandler.RegisterRoutes(a.echo)
	messageHandler.RegisterRoutes(a.echo)
	authHandler.RegisterAuthRoutes(a.echo)
	userHandler.RegisterRoutes(a.echo)
	searchHandler.RegisterRoutes(a.echo)
	prefsHandler.RegisterRoutes(a.echo)
	channelHandler.RegisterRoutes(a.echo)
	pinsHandler.RegisterRoutes(a.echo)
	friendsHandler.RegisterRoutes(a.echo)

	// * Server launch and graceful shutdown
	go func() {
		if err := a.echo.Start(a.cfg.AppPort); err != nil {
			log.Info().Err(err).Msg("Server shutting down")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info().Msg("Shutting down server...")

	grpcSrv.GracefulStop()
	log.Info().Msg("gRPC server stopped")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := a.echo.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	return nil
}
