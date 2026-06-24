// Package service — бизнес-логика donation-фичи.
//
// Donate переиспользует ту же инфраструктуру что и wallet.Transfer:
// custodial-подпись (дешифровка приватника → audit → broadcast), но
// помечает pending как DONATION и сохраняет donation_message, чтобы
// observer создал плоскую запись в tx_donations + system message в чате.
package service

import (
	"context"
	"errors"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet"
	walletrepo "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/blockchain"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/crypto"
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"
)

// Доменные ошибки. Handler маппит их в HTTP-коды.
var (
	ErrInvalidInput          = errors.New("invalid input")
	ErrSelfDonation          = errors.New("cannot donate to yourself")
	ErrRecipientNotFound     = errors.New("recipient not found")
	ErrRecipientHasNoWallet  = errors.New("recipient has no wallet")
	ErrChannelWalletNotFound = errors.New("channel treasury wallet not found")
	ErrWalletNotFound        = errors.New("wallet not found")
	ErrInsufficientBalance   = errors.New("insufficient balance")
	ErrChatNotFound          = errors.New("chat not found")
	ErrSenderNotInChat       = errors.New("sender is not a member of this chat")
	ErrRecipientNotInChat    = errors.New("recipient is not a member of this chat")
)

// ────────────────────────────────────────────────────────────
//  Маленькие интерфейсы для зависимостей (мокаются в тестах).
// ────────────────────────────────────────────────────────────

type walletReader interface {
	GetUserWallet(ctx context.Context, userID uuid.UUID) (*wallet.Wallet, error)
	GetChannelWallet(ctx context.Context, channelID uuid.UUID) (*wallet.ChannelWallet, error)
	WriteAudit(ctx context.Context, e walletrepo.AuditEntry) error
	InsertPendingDonation(ctx context.Context, txHash string, fromUserID uuid.UUID, chatID uuid.UUID, donationMessage string) error
}

type chainReader interface {
	SudaBalanceOf(ctx context.Context, addr common.Address) (*big.Int, error)
}

type encryptor interface {
	Decrypt(keyVersion uint8, encoded string) ([]byte, error)
}

type tokenSender interface {
	SendTokenTransfer(ctx context.Context, signer blockchain.Signer, to common.Address, amount *big.Int) (string, error)
}

type messengerClient interface {
	ResolveUsername(ctx context.Context, username string) (*platgrpc.ResolvedUser, error)
	CheckChatMembership(ctx context.Context, chatID string, userIDs []string) (*platgrpc.ChatMembershipResult, error)
}

// Service — бизнес-логика донатов.
type Service struct {
	repo      walletReader
	encryptor encryptor
	reader    chainReader
	sender    tokenSender
	messenger messengerClient
}

// Deps — зависимости конструктора.
type Deps struct {
	Postgres        *pgxpool.Pool
	Encryptor       *crypto.Encryptor
	Contracts       *blockchain.Contracts
	Reader          *blockchain.Reader
	Broadcaster     *blockchain.Broadcaster
	MessengerClient *platgrpc.MessengerClient
}

func New(d Deps) *Service {
	return &Service{
		repo:      walletrepo.New(d.Postgres),
		encryptor: d.Encryptor,
		reader:    d.Reader,
		sender:    blockchain.NewTokenSender(d.Broadcaster, d.Contracts),
		messenger: d.MessengerClient,
	}
}

// NewWithDeps — для тестов: позволяет инжектить моки напрямую.
func NewWithDeps(repo walletReader, enc encryptor, reader chainReader, sender tokenSender, mc messengerClient) *Service {
	return &Service{
		repo:      repo,
		encryptor: enc,
		reader:    reader,
		sender:    sender,
		messenger: mc,
	}
}
