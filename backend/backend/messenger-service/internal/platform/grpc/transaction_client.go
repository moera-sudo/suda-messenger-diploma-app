package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/moera-sudo/backend/backend/proto/transaction"
)

// TransactionClient — gRPC-клиент к transaction-service.
//
// Используется в messenger-service:
//   - auth/service.Verify  → CreateWalletForUser (создание кошелька + welcome bonus)
//   - channel/service.Create → CreateWalletForChannel (создание treasury)
//   - user/service.GetMe   → GetBalance (отображение баланса в профиле)
//   - channel/service.Subscribe → CheckTokenGating (проверка платных каналов)
//
// Все методы возвращают доменные типы (string, bool), а не pb-структуры,
// чтобы вызывающий код не зависел от proto-deps.
type TransactionClient struct {
	client pb.TransactionServiceClient
	conn   *grpc.ClientConn
}

// NewTransactionClient устанавливает gRPC-соединение с transaction-service.
// Соединение ленивое: NewClient не блокирует, реальный коннект случится
// на первом запросе. Если transaction-service в этот момент недоступен —
// первый вызов вернёт ошибку, но сам конструктор не упадёт.
func NewTransactionClient(addr string) (*TransactionClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("transaction client: dial %s: %w", addr, err)
	}

	return &TransactionClient{
		client: pb.NewTransactionServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close закрывает gRPC-соединение. Вызывать при graceful shutdown.
func (c *TransactionClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ────────────────────────────────────────────────────────────
//  Wallets
// ────────────────────────────────────────────────────────────

// CreatedWallet — результат CreateWalletFor{User,Channel}.
type CreatedWallet struct {
	Address string // 0x... EVM address (42 chars)
	Existed bool   // true если кошелёк уже был, false — создан в этом вызове
}

// CreateWalletForUser создаёт кастодиальный кошелёк юзера через transaction-service.
// Идемпотентно: если уже есть — вернётся существующий address и Existed=true.
//
// Вызывается из auth/service.Verify ПОСЛЕ успешной email-верификации.
// Welcome bonus 100 SUDA уходит автоматически в момент создания нового кошелька.
func (c *TransactionClient) CreateWalletForUser(ctx context.Context, userID string) (*CreatedWallet, error) {
	resp, err := c.client.CreateWalletForUser(ctx, &pb.CreateUserWalletRequest{
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("transaction: CreateWalletForUser: %w", err)
	}
	return &CreatedWallet{
		Address: resp.GetAddress(),
		Existed: resp.GetExisted(),
	}, nil
}

// CreateWalletForChannel создаёт treasury-кошелёк канала.
// Идемпотентно, как и CreateWalletForUser. Без welcome bonus.
func (c *TransactionClient) CreateWalletForChannel(ctx context.Context, channelID string) (*CreatedWallet, error) {
	resp, err := c.client.CreateWalletForChannel(ctx, &pb.CreateChannelWalletRequest{
		ChannelId: channelID,
	})
	if err != nil {
		return nil, fmt.Errorf("transaction: CreateWalletForChannel: %w", err)
	}
	return &CreatedWallet{
		Address: resp.GetAddress(),
		Existed: resp.GetExisted(),
	}, nil
}

// GetWallet возвращает address юзера БЕЗ обращения к Besu (только БД).
// Быстрый метод для UI, когда баланс не нужен (например, в профиле).
func (c *TransactionClient) GetWallet(ctx context.Context, userID string) (string, error) {
	resp, err := c.client.GetWallet(ctx, &pb.UserRequest{UserId: userID})
	if err != nil {
		return "", fmt.Errorf("transaction: GetWallet: %w", err)
	}
	return resp.GetAddress(), nil
}

// GetChannelWallet возвращает address канала. Без обращения к Besu.
func (c *TransactionClient) GetChannelWallet(ctx context.Context, channelID string) (string, error) {
	resp, err := c.client.GetChannelWallet(ctx, &pb.ChannelRequest{ChannelId: channelID})
	if err != nil {
		return "", fmt.Errorf("transaction: GetChannelWallet: %w", err)
	}
	return resp.GetAddress(), nil
}

// ────────────────────────────────────────────────────────────
//  Balance
// ────────────────────────────────────────────────────────────

// Balance — текущий on-chain баланс юзера.
type Balance struct {
	SudaBalanceWei string // uint256 как строка
	Decimals       int32  // всегда 18
	Address        string // адрес для дебага
}

// GetBalance делает SudaToken.balanceOf() на Besu через transaction-service
// и возвращает реальный on-chain баланс юзера.
//
// Стоимость: 1 RPC-call к Besu (~10ms) + 1 UPDATE кеша. Не использовать
// в hot-path'ах (например, не звать на каждое сообщение в чате).
func (c *TransactionClient) GetBalance(ctx context.Context, userID string) (*Balance, error) {
	resp, err := c.client.GetBalance(ctx, &pb.UserRequest{UserId: userID})
	if err != nil {
		return nil, fmt.Errorf("transaction: GetBalance: %w", err)
	}
	return &Balance{
		SudaBalanceWei: resp.GetSudaBalanceWei(),
		Decimals:       resp.GetDecimals(),
		Address:        resp.GetAddress(),
	}, nil
}

// ────────────────────────────────────────────────────────────
//  Token Gating
// ────────────────────────────────────────────────────────────

// GatingResult — результат проверки правила доступа к чату/каналу.
type GatingResult struct {
	Required       bool
	Passed         bool
	MinBalanceWei  string
	UserBalanceWei string
	PriceWei       string // цена платной подписки в wei ("0" = не платный)
	Reason         string // "ok" | "no_wallet" | "insufficient_balance" | "nft_required"
}

// CheckTokenGating проверяет, может ли юзер войти в чат/канал по правилу
// в tx_gating_rules. Если правила нет — Required=false, чат open для всех.
//
// Вызывается из channel/service.Subscribe ПЕРЕД добавлением юзера в подписчики.
// При Passed=false возвращаем юзеру 403 с понятным сообщением.
func (c *TransactionClient) CheckTokenGating(ctx context.Context, chatID, userID string) (*GatingResult, error) {
	resp, err := c.client.CheckTokenGating(ctx, &pb.GatingRequest{
		ChatId: chatID,
		UserId: userID,
	})
	if err != nil {
		return nil, fmt.Errorf("transaction: CheckTokenGating: %w", err)
	}
	return &GatingResult{
		Required:       resp.GetRequired(),
		Passed:         resp.GetPassed(),
		MinBalanceWei:  resp.GetMinBalanceWei(),
		UserBalanceWei: resp.GetUserBalanceWei(),
		PriceWei:       resp.GetPriceWei(),
		Reason:         resp.GetReason(),
	}, nil
}

// SubscriptionCharge — результат списания за платную подписку.
type SubscriptionCharge struct {
	TxHash   string
	PriceWei string
}

// ChargeChannelSubscription списывает цену платной подписки с юзера в treasury
// канала. Вызывается из channel/service.Subscribe для платного PUBLIC-канала
// ПЕРЕД добавлением в подписчики. Ошибка → подписку не оформляем.
func (c *TransactionClient) ChargeChannelSubscription(ctx context.Context, userID, channelID string) (*SubscriptionCharge, error) {
	resp, err := c.client.ChargeChannelSubscription(ctx, &pb.SubscriptionRequest{
		UserId:    userID,
		ChannelId: channelID,
	})
	if err != nil {
		return nil, fmt.Errorf("transaction: ChargeChannelSubscription: %w", err)
	}
	return &SubscriptionCharge{TxHash: resp.GetTxHash(), PriceWei: resp.GetPriceWei()}, nil
}