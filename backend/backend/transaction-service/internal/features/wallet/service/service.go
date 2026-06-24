// Package service — бизнес-логика wallet-фичи.
//
// Service имплементирует grpc.TransactionServiceImpl (см. internal/platform/grpc/server.go)
// и предоставляет HTTP-методы для wallet handler'а.
//
// Custodial-флоу: все операции пользователя инициирует бэк, который держит
// зашифрованные приватники в tx_wallets. При каждой подписи:
//   1. Дешифровка приватника
//   2. Запись в tx_signing_audit (ДО broadcast'а — полный journal)
//   3. Подпись + broadcast
//   4. Запись в tx_pending (UI-индикатор)
//   5. Возврат tx_hash юзеру
//   6. Async: observer ловит подтверждение → tx_transactions → WS event
package service

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/blockchain"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/crypto"
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"
)

// ────────────────────────────────────────────────────────────
//  Доменные ошибки.
//  Handler маппит их в HTTP коды, gRPC использует напрямую.
// ────────────────────────────────────────────────────────────

var (
	ErrWalletNotFound       = errors.New("wallet not found")
	ErrInvalidInput         = errors.New("invalid input")
	ErrSelfTransfer         = errors.New("cannot transfer to self")
	ErrRecipientNotFound    = errors.New("recipient not found")
	ErrRecipientHasNoWallet = errors.New("recipient has no wallet")
	ErrInsufficientBalance  = errors.New("insufficient balance")
	ErrForbidden            = errors.New("forbidden")
	ErrNotPaidChannel       = errors.New("channel is not a paid channel")
)

// ────────────────────────────────────────────────────────────
//  Service
// ────────────────────────────────────────────────────────────

type Service struct {
	repo      *repository.Repository
	redis     *redis.Client
	encryptor *crypto.Encryptor

	chain     *blockchain.Client
	contracts *blockchain.Contracts
	reader    *blockchain.Reader
	broadcast *blockchain.Broadcaster
	treasury  *blockchain.TreasurySigner

	messenger *platgrpc.MessengerClient

	ts *TreasuryService // treasury stats sub-service (interface-typed deps)

	welcomeBonusSudaWei string // 100 SUDA в wei (строка)
}

// Deps — все зависимости для конструктора Service.
type Deps struct {
	Postgres            *pgxpool.Pool
	Redis               *redis.Client
	Encryptor           *crypto.Encryptor
	Chain               *blockchain.Client
	Contracts           *blockchain.Contracts
	Reader              *blockchain.Reader
	Broadcaster         *blockchain.Broadcaster
	Treasury            *blockchain.TreasurySigner
	MessengerClient     *platgrpc.MessengerClient
	WelcomeBonusSudaWei string
}

func New(d Deps) *Service {
	repo := repository.New(d.Postgres)
	return &Service{
		repo:                repo,
		redis:               d.Redis,
		encryptor:           d.Encryptor,
		chain:               d.Chain,
		contracts:           d.Contracts,
		reader:              d.Reader,
		broadcast:           d.Broadcaster,
		treasury:            d.Treasury,
		messenger:           d.MessengerClient,
		ts:                  NewTreasuryService(repo, d.Reader, d.MessengerClient),
		welcomeBonusSudaWei: d.WelcomeBonusSudaWei,
	}
}