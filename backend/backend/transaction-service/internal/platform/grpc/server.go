package grpc

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/moera-sudo/backend/backend/proto/transaction"
	
)

// TransactionServiceImpl — интерфейс, который имплементирует feature-слой
// (wallet/gating). gRPC-слой не знает о бизнес-логике; он только
// принимает proto-запросы, делегирует и возвращает proto-ответы.
//
// Зачем интерфейс, а не прямые ссылки на репозитории: чтобы gRPC server
// мог конструироваться в app.go от любого набора зависимостей, и чтобы
// можно было подменять реализации в тестах.
type TransactionServiceImpl interface {
	CreateWalletForUser(ctx context.Context, userID string) (address string, existed bool, err error)
	CreateWalletForChannel(ctx context.Context, channelID string) (address string, existed bool, err error)
	GetWallet(ctx context.Context, userID string) (address string, err error)
	GetChannelWallet(ctx context.Context, channelID string) (address string, err error)
	GetBalance(ctx context.Context, userID string) (balanceWei string, address string, err error)
	CheckTokenGating(ctx context.Context, chatID, userID string) (passed, required bool, minBalanceWei, userBalanceWei, priceWei, reason string, err error)
	ChargeChannelSubscription(ctx context.Context, userID, channelID string) (txHash, priceWei string, err error)
}

// ErrWalletNotFound — типизированная ошибка для feature-слоя.
// gRPC-handler преобразует её в подходящий gRPC status code.
var (
	ErrWalletNotFound = errors.New("wallet not found")
	ErrUserNotFound   = errors.New("user not found")
)

// Server инкапсулирует grpc.Server и адрес для listen'а.
type Server struct {
	grpcServer *grpc.Server
	listener   net.Listener
	addr       string
}

// NewServer создаёт gRPC-сервер с зарегистрированным TransactionService.
// Не запускает listen — это делает Start.
func NewServer(addr string, impl TransactionServiceImpl) (*Server, error) {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, fmt.Errorf("grpc: listen %s: %w", addr, err)
	}

	srv := grpc.NewServer(
		grpc.UnaryInterceptor(loggingInterceptor),
	)

	transaction.RegisterTransactionServiceServer(srv, &transactionHandler{impl: impl})
	reflection.Register(srv)

	return &Server{
		grpcServer: srv,
		listener:   lis,
		addr:       addr,
	}, nil
}

// Start запускает gRPC-сервер в текущей горутине (БЛОКИРУЮЩИЙ вызов).
// В app.go вызывается через `go server.Start()`.
func (s *Server) Start() error {
	log.Info().Str("addr", s.addr).Msg("gRPC server listening")
	return s.grpcServer.Serve(s.listener)
}

// Stop делает graceful shutdown — ждёт окончания in-flight запросов.
func (s *Server) Stop() {
	log.Info().Msg("gRPC server stopping")
	s.grpcServer.GracefulStop()
}

// ────────────────────────────────────────────────────────────
//  Handler — proto-обёртка над интерфейсом TransactionServiceImpl
// ────────────────────────────────────────────────────────────

type transactionHandler struct {
	transaction.UnimplementedTransactionServiceServer
	impl TransactionServiceImpl
}

func (h *transactionHandler) CreateWalletForUser(
	ctx context.Context, req *transaction.CreateUserWalletRequest,
) (*transaction.WalletResponse, error) {
	address, existed, err := h.impl.CreateWalletForUser(ctx, req.GetUserId())
	if err != nil {
		log.Error().Err(err).Str("user_id", req.GetUserId()).Msg("CreateWalletForUser failed")
		return nil, err
	}
	return &transaction.WalletResponse{Address: address, Existed: existed}, nil
}

func (h *transactionHandler) CreateWalletForChannel(
	ctx context.Context, req *transaction.CreateChannelWalletRequest,
) (*transaction.WalletResponse, error) {
	address, existed, err := h.impl.CreateWalletForChannel(ctx, req.GetChannelId())
	if err != nil {
		log.Error().Err(err).Str("channel_id", req.GetChannelId()).Msg("CreateWalletForChannel failed")
		return nil, err
	}
	return &transaction.WalletResponse{Address: address, Existed: existed}, nil
}

func (h *transactionHandler) GetWallet(
	ctx context.Context, req *transaction.UserRequest,
) (*transaction.WalletResponse, error) {
	address, err := h.impl.GetWallet(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, ErrWalletNotFound) {
			return nil, err
		}
		log.Error().Err(err).Str("user_id", req.GetUserId()).Msg("GetWallet failed")
		return nil, err
	}
	return &transaction.WalletResponse{Address: address, Existed: true}, nil
}

func (h *transactionHandler) GetChannelWallet(
	ctx context.Context, req *transaction.ChannelRequest,
) (*transaction.WalletResponse, error) {
	address, err := h.impl.GetChannelWallet(ctx, req.GetChannelId())
	if err != nil {
		if errors.Is(err, ErrWalletNotFound) {
			return nil, err
		}
		log.Error().Err(err).Str("channel_id", req.GetChannelId()).Msg("GetChannelWallet failed")
		return nil, err
	}
	return &transaction.WalletResponse{Address: address, Existed: true}, nil
}

func (h *transactionHandler) GetBalance(
	ctx context.Context, req *transaction.UserRequest,
) (*transaction.BalanceResponse, error) {
	balanceWei, address, err := h.impl.GetBalance(ctx, req.GetUserId())
	if err != nil {
		log.Error().Err(err).Str("user_id", req.GetUserId()).Msg("GetBalance failed")
		return nil, err
	}
	return &transaction.BalanceResponse{
		SudaBalanceWei: balanceWei,
		Decimals:       18,
		Address:        address,
	}, nil
}

func (h *transactionHandler) CheckTokenGating(
	ctx context.Context, req *transaction.GatingRequest,
) (*transaction.GatingResponse, error) {
	passed, required, minBal, userBal, priceWei, reason, err := h.impl.CheckTokenGating(ctx, req.GetChatId(), req.GetUserId())
	if err != nil {
		log.Error().Err(err).Str("chat_id", req.GetChatId()).Str("user_id", req.GetUserId()).Msg("CheckTokenGating failed")
		return nil, err
	}
	return &transaction.GatingResponse{
		Required:       required,
		Passed:         passed,
		MinBalanceWei:  minBal,
		UserBalanceWei: userBal,
		PriceWei:       priceWei,
		Reason:         reason,
	}, nil
}

func (h *transactionHandler) ChargeChannelSubscription(
	ctx context.Context, req *transaction.SubscriptionRequest,
) (*transaction.SubscriptionResponse, error) {
	txHash, priceWei, err := h.impl.ChargeChannelSubscription(ctx, req.GetUserId(), req.GetChannelId())
	if err != nil {
		log.Error().Err(err).Str("user_id", req.GetUserId()).Str("channel_id", req.GetChannelId()).Msg("ChargeChannelSubscription failed")
		return nil, err
	}
	return &transaction.SubscriptionResponse{TxHash: txHash, PriceWei: priceWei}, nil
}