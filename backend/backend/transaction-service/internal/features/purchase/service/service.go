// Package service — бизнес-логика фичи purchase (имитация покупки SUDA).
//
// Реальной оплаты нет. Flow:
//
//  1. POST /purchase/initiate          — создаёт PENDING row, возвращает purchase_id и сумму.
//  2. POST /purchase/:id/confirm       — переводит PENDING -> PROCESSING, делает sleep 2-3 сек
//     (имитация процессинга), затем treasury шлёт SUDA пользователю Transfer'ом, статус COMPLETED.
//  3. Observer индексирует Transfer event (suda_token handler), видит match по tx_hash в
//     tx_suda_purchases, шлёт WS-event PURCHASE_COMPLETED.
//
// Rate-limit на confirm — Redis SET NX TTL=5s, защита от случайного двойного нажатия.
package service

import (
	"context"
	"errors"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet"
	walletrepo "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/blockchain"
)

// Доменные ошибки.
var (
	ErrInvalidInput   = errors.New("invalid input")
	ErrPackageUnknown = errors.New("unknown package code")
	ErrNotFound       = errors.New("purchase not found")
	ErrForbidden      = errors.New("forbidden")
	ErrAlreadyHandled = errors.New("purchase is not in PENDING state")
	ErrWalletNotFound = errors.New("user has no wallet")
	ErrRateLimited    = errors.New("too many purchase confirmations")
)

// ────────────────────────────────────────────────────────────
//  Маленькие интерфейсы для зависимостей — позволяют легко мокать в тестах.
//  Concrete-типы из platform/repository неявно их удовлетворяют (Go duck-typing).
// ────────────────────────────────────────────────────────────

type purchaseRepo interface {
	Create(ctx context.Context, p repository.CreateParams) (uuid.UUID, error)
	GetByID(ctx context.Context, id uuid.UUID) (*purchase.Purchase, error)
	MarkProcessing(ctx context.Context, id uuid.UUID) error
	MarkCompleted(ctx context.Context, id uuid.UUID, txHash string) error
	MarkFailed(ctx context.Context, id uuid.UUID, reason string) error
	ListForUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]purchase.Purchase, error)
}

type walletReader interface {
	GetUserWallet(ctx context.Context, userID uuid.UUID) (*wallet.Wallet, error)
	WriteAudit(ctx context.Context, e walletrepo.AuditEntry) error
	InsertPending(ctx context.Context, txHash string, fromUserID uuid.UUID, expectedType string, chatID uuid.UUID) error
}

// tokenSender — broadcast Token.Transfer одним вызовом (см. blockchain.TokenSender).
type tokenSender interface {
	SendTokenTransfer(ctx context.Context, signer blockchain.Signer, to common.Address, amount *big.Int) (string, error)
}

// treasurySigner — нам нужен только Address; SignTx/TransactOpts вызываются внутри tokenSender.
type treasurySigner interface {
	blockchain.Signer
}

// rateLimiter — обёртка над Redis SETNX. В продакшене — redis.Client; в тестах — miniredis или fake.
type rateLimiter interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
}

// Service — бизнес-логика purchase'ов.
type Service struct {
	repo       purchaseRepo
	walletRepo walletReader
	redis      rateLimiter
	sender     tokenSender
	treasury   treasurySigner

	// simulateProcessing — функция-имитация платёжного процессинга. В продакшене
	// делает sleep 2-3 сек, в тестах заменяется на no-op для скорости.
	simulateProcessing func(ctx context.Context) error
}

// Deps — все зависимости конструктора.
type Deps struct {
	Postgres    *pgxpool.Pool
	Redis       *redis.Client
	Contracts   *blockchain.Contracts
	Broadcaster *blockchain.Broadcaster
	Treasury    *blockchain.TreasurySigner
}

func New(d Deps) *Service {
	return &Service{
		repo:               repository.New(d.Postgres),
		walletRepo:         walletrepo.New(d.Postgres),
		redis:              d.Redis,
		sender:             blockchain.NewTokenSender(d.Broadcaster, d.Contracts),
		treasury:           d.Treasury,
		simulateProcessing: defaultSimulateProcessing,
	}
}

// NewWithDeps — расширенный конструктор для тестов: позволяет инжектить
// мок-репозитории, мок tokenSender и no-op simulateProcessing напрямую.
func NewWithDeps(
	repo purchaseRepo,
	walletRepo walletReader,
	rl rateLimiter,
	sender tokenSender,
	treasury treasurySigner,
	simulateFn func(ctx context.Context) error,
) *Service {
	if simulateFn == nil {
		simulateFn = defaultSimulateProcessing
	}
	return &Service{
		repo:               repo,
		walletRepo:         walletRepo,
		redis:              rl,
		sender:             sender,
		treasury:           treasury,
		simulateProcessing: simulateFn,
	}
}

// auditSubjectUser — re-export для красоты при использовании в этом пакете.
const auditSubjectUser = wallet.SubjectUser
