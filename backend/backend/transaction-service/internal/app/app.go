// Package app — DI-контейнер и lifecycle для transaction-service.
package app

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/moera-sudo/backend/backend/transaction-service/config"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/observer"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/blockchain"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/crypto"
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/postgres"
	platredis "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/redis"

	httputil "github.com/moera-sudo/backend/backend/transaction-service/internal/pkg/http"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/pkg/middlewares"

	donationhttp "github.com/moera-sudo/backend/backend/transaction-service/internal/features/donation/delivery/http"
	donationsvc "github.com/moera-sudo/backend/backend/transaction-service/internal/features/donation/service"
	gatinghttp "github.com/moera-sudo/backend/backend/transaction-service/internal/features/gating/delivery/http"
	gatingsvc "github.com/moera-sudo/backend/backend/transaction-service/internal/features/gating/service"
	marketplacehttp "github.com/moera-sudo/backend/backend/transaction-service/internal/features/marketplace/delivery/http"
	nfthttp "github.com/moera-sudo/backend/backend/transaction-service/internal/features/nft/delivery/http"
	purchasehttp "github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase/delivery/http"
	purchasesvc "github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase/service"
	questhttp "github.com/moera-sudo/backend/backend/transaction-service/internal/features/quest/delivery/http"
	wallethttp "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/delivery/http"
	walletsvc "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/service"
)

type App struct {
	cfg *config.Config

	pg        *pgxpool.Pool
	redis     *redis.Client
	encryptor *crypto.Encryptor

	chain     *blockchain.Client
	contracts *blockchain.Contracts
	reader    *blockchain.Reader
	broadcast *blockchain.Broadcaster
	nonces    *blockchain.NonceManager
	treasury  *blockchain.TreasurySigner

	messenger *platgrpc.MessengerClient

	walletService   *walletsvc.Service
	purchaseService *purchasesvc.Service
	donationService *donationsvc.Service
	gatingService   *gatingsvc.Service
	observer        *observer.Observer

	echo       *echo.Echo
	grpcServer *platgrpc.Server
}

// Init создаёт все зависимости в правильном порядке.
func Init(ctx context.Context, cfg *config.Config) (*App, error) {
	a := &App{cfg: cfg}

	// 1. PostgreSQL
	pg, err := postgres.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("init postgres: %w", err)
	}
	a.pg = pg
	log.Info().Msg("postgres connected")

	// 2. Redis
	rdb, err := platredis.New(ctx, cfg.RedisAddr, cfg.RedisPass)
	if err != nil {
		a.closeAll()
		return nil, fmt.Errorf("init redis: %w", err)
	}
	a.redis = rdb
	log.Info().Msg("redis connected")

	// 3. Encryptor
	enc, err := crypto.NewEncryptor(1, cfg.WalletEncryptionKey)
	if err != nil {
		a.closeAll()
		return nil, fmt.Errorf("init encryptor: %w", err)
	}
	a.encryptor = enc
	log.Info().Uint8("current_key_version", enc.CurrentVersion()).Msg("wallet encryptor ready")

	// 4. Blockchain client + sanity check
	chain, err := blockchain.New(ctx, cfg.BesuRPCURL, cfg.ChainID)
	if err != nil {
		a.closeAll()
		return nil, fmt.Errorf("init blockchain client: %w", err)
	}
	a.chain = chain

	// 5. Treasury signer
	treasury, err := blockchain.NewTreasurySigner(cfg.DeployerPrivateKey)
	if err != nil {
		a.closeAll()
		return nil, fmt.Errorf("init treasury signer: %w", err)
	}
	a.treasury = treasury
	log.Info().Str("treasury_addr", treasury.Address().Hex()).Msg("treasury signer ready")

	// 6. Контракты
	contracts, err := blockchain.NewContracts(chain, blockchain.ContractAddresses{
		Token:       cfg.SudaTokenAddress,
		NFT:         cfg.SudaNFTAddress,
		Marketplace: cfg.SudaMarketplaceAddress,
		Escrow:      cfg.SudaEscrowAddress,
		Fundraising: cfg.SudaFundraisingAddress,
	})
	if err != nil {
		a.closeAll()
		return nil, fmt.Errorf("init contracts: %w", err)
	}
	a.contracts = contracts

	// 6.1. Sanity check
	a.reader = blockchain.NewReader(chain, contracts)
	if err := a.verifyContractsAlive(ctx); err != nil {
		a.closeAll()
		return nil, fmt.Errorf("contracts verification: %w", err)
	}
	log.Info().Msg("contracts verified on-chain")

	// 7. Nonce manager + broadcaster
	a.nonces = blockchain.NewNonceManager()
	a.broadcast = blockchain.NewBroadcaster(chain, a.nonces)

	// 8. Messenger gRPC client
	mc, err := platgrpc.NewMessengerClient(cfg.MessengerGRPCAddr)
	if err != nil {
		a.closeAll()
		return nil, fmt.Errorf("init messenger client: %w", err)
	}
	a.messenger = mc

	// 9. Wallet service (имплементирует grpc.TransactionServiceImpl)
	a.walletService = walletsvc.New(walletsvc.Deps{
		Postgres:            pg,
		Redis:               rdb,
		Encryptor:           enc,
		Chain:               chain,
		Contracts:           contracts,
		Reader:              a.reader,
		Broadcaster:         a.broadcast,
		Treasury:            treasury,
		MessengerClient:     mc,
		WelcomeBonusSudaWei: welcomeBonusSudaWei(),
	})

	// 9.1. Purchase service (имитация покупки SUDA)
	a.purchaseService = purchasesvc.New(purchasesvc.Deps{
		Postgres:    pg,
		Redis:       rdb,
		Contracts:   contracts,
		Broadcaster: a.broadcast,
		Treasury:    treasury,
	})

	// 9.2. Donation service (донаты юзерам и каналам)
	a.donationService = donationsvc.New(donationsvc.Deps{
		Postgres:        pg,
		Encryptor:       enc,
		Contracts:       contracts,
		Reader:          a.reader,
		Broadcaster:     a.broadcast,
		MessengerClient: mc,
	})

	// 9.3. Gating service (CRUD token-gating правил каналов)
	a.gatingService = gatingsvc.New(gatingsvc.Deps{
		Postgres:        pg,
		MessengerClient: mc,
	})

	// 10. Observer
	a.observer = observer.New(observer.Deps{
		Postgres:        pg,
		Chain:           chain,
		Contracts:       contracts,
		MessengerClient: mc,
	})

	// 11. Echo HTTP server
	a.echo = a.setupEcho()

	// 12. gRPC server
	grpcSrv, err := platgrpc.NewServer(cfg.GRPCPort, a.walletService)
	if err != nil {
		a.closeAll()
		return nil, fmt.Errorf("init grpc server: %w", err)
	}
	a.grpcServer = grpcSrv

	log.Info().Msg("app initialized successfully")
	return a, nil
}

// Run запускает все компоненты до ctx.Done().
func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 3)

	go func() {
		log.Info().Str("addr", a.cfg.AppPort).Msg("HTTP server listening")
		if err := a.echo.Start(a.cfg.AppPort); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- fmt.Errorf("http: %w", err)
		}
	}()

	go func() {
		if err := a.grpcServer.Start(); err != nil {
			errCh <- fmt.Errorf("grpc: %w", err)
		}
	}()

	go a.observer.Run(ctx)

	select {
	case <-ctx.Done():
		log.Info().Msg("shutdown signal received")
	case err := <-errCh:
		log.Error().Err(err).Msg("fatal error in background goroutine")
		return err
	}

	a.shutdown()
	return nil
}

func (a *App) shutdown() {
	log.Info().Msg("starting graceful shutdown")

	httpCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := a.echo.Shutdown(httpCtx); err != nil {
		log.Error().Err(err).Msg("http shutdown")
	}

	if a.grpcServer != nil {
		a.grpcServer.Stop()
	}

	a.closeAll()

	log.Info().Msg("shutdown complete")
}

func (a *App) closeAll() {
	if a.messenger != nil {
		_ = a.messenger.Close()
	}
	if a.chain != nil {
		a.chain.Close()
	}
	if a.redis != nil {
		_ = a.redis.Close()
	}
	if a.pg != nil {
		a.pg.Close()
	}
}

func (a *App) setupEcho() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.Validator = httputil.NewValidator()

	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch},
		AllowHeaders: []string{"*"},
	}))

	// Health check — БЕЗ gateway middleware
	e.GET("/healthz", func(c echo.Context) error {
		return httputil.RespondByMessage(c, http.StatusOK, "ok", "transaction-service is alive")
	})

	// Swagger UI — БЕЗ gateway middleware (доступен напрямую для разработки)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	api := e.Group("/api/v1")
	api.Use(middlewares.GatewaySecret(a.cfg.GatewaySignatureSecret))

	// Wallet — теперь с реальным service
	wallethttp.NewHandler(a.walletService).RegisterRoutes(api.Group("/wallet"))

	// Purchase — имитация покупки SUDA (Stage 4.5)
	purchasehttp.NewHandler(a.purchaseService).RegisterRoutes(api.Group("/purchase"))

	// Donation — донаты юзерам и каналам (Stage 4.1)
	donationhttp.NewHandler(a.donationService).RegisterRoutes(api.Group("/donate"))

	// Gating — CRUD token-gating правил каналов (Stage 4.2)
	gatinghttp.NewHandler(a.gatingService).RegisterRoutes(api.Group("/gating"))

	// Остальные фичи — заглушки (итерация Г)
	nfthttp.NewHandler().RegisterRoutes(api.Group("/nft"))
	marketplacehttp.NewHandler().RegisterRoutes(api.Group("/marketplace"))
	questhttp.NewHandler().RegisterRoutes(api.Group("/quests"))

	return e
}

func (a *App) verifyContractsAlive(ctx context.Context) error {
	checks := []struct {
		name string
		addr string
	}{
		{"SudaToken", a.contracts.TokenAddr.Hex()},
		{"SudaNFT", a.contracts.NFTAddr.Hex()},
		{"SudaMarketplace", a.contracts.MarketplaceAddr.Hex()},
		{"SudaEscrow", a.contracts.EscrowAddr.Hex()},
		{"SudaFundraising", a.contracts.FundraisingAddr.Hex()},
	}

	for _, c := range checks {
		has, err := a.reader.HasContractCode(ctx, a.contracts.TokenAddr)
		if err != nil {
			return fmt.Errorf("check %s code: %w", c.name, err)
		}
		if !has {
			return fmt.Errorf("%s at %s has no code (not deployed?)", c.name, c.addr)
		}
	}
	return nil
}

// welcomeBonusSudaWei — 100 SUDA = 100 × 10^18 в wei как строка.
func welcomeBonusSudaWei() string {
	bonus := new(big.Int).Mul(big.NewInt(100), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	return bonus.String()
}