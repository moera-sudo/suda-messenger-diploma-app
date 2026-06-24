package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"
)

// GetBalance — gRPC: вернуть текущий on-chain баланс юзера.
//
// Возвращает: (balance_wei, address, err).
// Параллельно обновляет suda_balance_cache в БД (best-effort).
func (s *Service) GetBalance(ctx context.Context, userIDStr string) (string, string, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", "", fmt.Errorf("%w: user_id must be uuid", ErrInvalidInput)
	}

	w, err := s.repo.GetUserWallet(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", ErrWalletNotFound
		}
		return "", "", fmt.Errorf("get balance: lookup: %w", err)
	}

	balance, err := s.reader.SudaBalanceOf(ctx, common.HexToAddress(w.Address))
	if err != nil {
		return "", "", fmt.Errorf("get balance: on-chain read: %w", err)
	}

	// Обновляем кеш в БД (best-effort).
	if err := s.repo.UpdateUserBalanceCache(ctx, userID, balance); err != nil {
		log.Warn().Err(err).Str("user_id", userID.String()).Msg("update balance cache failed")
	}

	return balance.String(), w.Address, nil
}

// GetMyWalletAndBalance — HTTP-helper для GET /wallet/me.
// Возвращает: (address, balance_wei).
func (s *Service) GetMyWalletAndBalance(ctx context.Context, userID uuid.UUID) (string, string, error) {
	balanceWei, addr, err := s.GetBalance(ctx, userID.String())
	if err != nil {
		return "", "", err
	}
	return addr, balanceWei, nil
}

// GetChannelWalletAndBalance — HTTP-helper для GET /wallet/channel/:id.
//
// Проверяет, что вызывающий имеет права OWNER или ADMIN канала
// через gRPC к messenger. Без прав — ErrForbidden.
func (s *Service) GetChannelWalletAndBalance(
	ctx context.Context, userID, channelID uuid.UUID,
) (string, string, error) {
	// 1. Permission check через messenger.
	perm, err := s.messenger.CheckChannelPermission(
		ctx, userID.String(), channelID.String(), platgrpc.PermissionOwnerOrAdmin,
	)
	if err != nil {
		log.Error().Err(err).
			Str("user_id", userID.String()).
			Str("channel_id", channelID.String()).
			Msg("permission check failed")
		return "", "", fmt.Errorf("get channel wallet: permission check: %w", err)
	}
	if !perm.Granted {
		return "", "", ErrForbidden
	}

	// 2. Кошелёк канала.
	w, err := s.repo.GetChannelWallet(ctx, channelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", ErrWalletNotFound
		}
		return "", "", fmt.Errorf("get channel wallet: lookup: %w", err)
	}

	// 3. On-chain баланс.
	balance, err := s.reader.SudaBalanceOf(ctx, common.HexToAddress(w.Address))
	if err != nil {
		return "", "", fmt.Errorf("get channel wallet: on-chain read: %w", err)
	}

	return w.Address, balance.String(), nil
}