package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/blockchain"
)

// CreateWalletForUser — gRPC handler для регистрации.
//
// Idempotent:
//  1. SELECT — если кошелёк уже есть, возвращаем его address + existed=true.
//  2. Иначе генерируем новую пару, шифруем приватник, INSERT.
//     Если INSERT упал с UNIQUE constraint (гонка двух регистраций) — снова SELECT.
//  3. Welcome bonus: 100 SUDA от treasury → новому юзеру.
//     Если bonus упал (RPC down) — кошелёк всё равно создан.
func (s *Service) CreateWalletForUser(ctx context.Context, userIDStr string) (address string, existed bool, err error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", false, fmt.Errorf("%w: user_id must be uuid", ErrInvalidInput)
	}

	// 1. Если уже есть — возвращаем.
	if w, errGet := s.repo.GetUserWallet(ctx, userID); errGet == nil {
		return w.Address, true, nil
	} else if !errors.Is(errGet, repository.ErrNotFound) {
		return "", false, fmt.Errorf("create wallet: lookup: %w", errGet)
	}

	// 2. Генерируем новую пару ключей secp256k1.
	hexPK, addr, err := blockchain.GenerateNewWallet()
	if err != nil {
		return "", false, fmt.Errorf("create wallet: gen key: %w", err)
	}

	// 3. Шифруем приватник под текущей версией ключа.
	encryptedPK, err := s.encryptor.Encrypt([]byte(hexPK))
	if err != nil {
		return "", false, fmt.Errorf("create wallet: encrypt: %w", err)
	}

	// 4. Вставляем в БД.
	err = s.repo.CreateUserWallet(ctx, userID, addr.Hex(), encryptedPK, s.encryptor.CurrentVersion())
	if errors.Is(err, repository.ErrAlreadyExists) {
		// Гонка: пока генерировали ключи, другой call успел вставить.
		w, errGet := s.repo.GetUserWallet(ctx, userID)
		if errGet != nil {
			return "", false, fmt.Errorf("create wallet: refetch after race: %w", errGet)
		}
		return w.Address, true, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("create wallet: insert: %w", err)
	}

	log.Info().
		Str("user_id", userID.String()).
		Str("address", addr.Hex()).
		Msg("user wallet created")

	// 5. Welcome bonus — 100 SUDA от treasury (best-effort).
	if err := s.sendWelcomeBonus(ctx, userID, addr.Hex()); err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID.String()).
			Str("address", addr.Hex()).
			Msg("welcome bonus failed (wallet created without bonus)")
	}

	return addr.Hex(), false, nil
}

// sendWelcomeBonus переводит welcomeBonusSudaWei от treasury новому юзеру.
func (s *Service) sendWelcomeBonus(ctx context.Context, userID uuid.UUID, toAddress string) error {
	if s.welcomeBonusSudaWei == "" || s.welcomeBonusSudaWei == "0" {
		return nil
	}

	bonusAmount, ok := new(big.Int).SetString(s.welcomeBonusSudaWei, 10)
	if !ok {
		return fmt.Errorf("welcome bonus: bad amount %q", s.welcomeBonusSudaWei)
	}

	opts, release, err := s.broadcast.PrepareOpts(ctx, s.treasury)
	if err != nil {
		return fmt.Errorf("welcome bonus: prepare opts: %w", err)
	}
	defer release()

	tx, err := s.contracts.Token.Transfer(opts, common.HexToAddress(toAddress), bonusAmount)
	if err != nil {
		return fmt.Errorf("welcome bonus: send transfer: %w", err)
	}

	// Audit (бенефициар — userID получатель).
	_ = s.repo.WriteAudit(ctx, repository.AuditEntry{
		SubjectType: wallet.SubjectUser,
		SubjectID:   userID,
		Operation:   "WELCOME_BONUS",
		TxHash:      tx.Hash().Hex(),
		RequestIP:   "system",
		UserAgent:   "transaction-service/welcome",
	})

	// Pending для UI.
	_ = s.repo.InsertPending(ctx, tx.Hash().Hex(), userID, "P2P_TRANSFER", uuid.Nil)

	log.Info().
		Str("user_id", userID.String()).
		Str("amount_wei", bonusAmount.String()).
		Str("tx_hash", tx.Hash().Hex()).
		Msg("welcome bonus broadcast")

	return nil
}

// CreateWalletForChannel — то же для канала. Без welcome bonus.
func (s *Service) CreateWalletForChannel(ctx context.Context, channelIDStr string) (address string, existed bool, err error) {
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		return "", false, fmt.Errorf("%w: channel_id must be uuid", ErrInvalidInput)
	}

	if w, errGet := s.repo.GetChannelWallet(ctx, channelID); errGet == nil {
		return w.Address, true, nil
	} else if !errors.Is(errGet, repository.ErrNotFound) {
		return "", false, fmt.Errorf("create channel wallet: lookup: %w", errGet)
	}

	hexPK, addr, err := blockchain.GenerateNewWallet()
	if err != nil {
		return "", false, fmt.Errorf("create channel wallet: gen key: %w", err)
	}

	encryptedPK, err := s.encryptor.Encrypt([]byte(hexPK))
	if err != nil {
		return "", false, fmt.Errorf("create channel wallet: encrypt: %w", err)
	}

	err = s.repo.CreateChannelWallet(ctx, channelID, addr.Hex(), encryptedPK, s.encryptor.CurrentVersion())
	if errors.Is(err, repository.ErrAlreadyExists) {
		w, errGet := s.repo.GetChannelWallet(ctx, channelID)
		if errGet != nil {
			return "", false, fmt.Errorf("create channel wallet: refetch: %w", errGet)
		}
		return w.Address, true, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("create channel wallet: insert: %w", err)
	}

	log.Info().
		Str("channel_id", channelID.String()).
		Str("address", addr.Hex()).
		Msg("channel wallet created")

	return addr.Hex(), false, nil
}

// GetWallet — gRPC: вернуть address юзера без RPC к Besu.
func (s *Service) GetWallet(ctx context.Context, userIDStr string) (string, error) {
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", fmt.Errorf("%w: user_id must be uuid", ErrInvalidInput)
	}
	w, err := s.repo.GetUserWallet(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrWalletNotFound
		}
		return "", fmt.Errorf("get wallet: %w", err)
	}
	return w.Address, nil
}

// GetChannelWallet — gRPC: вернуть address канала без RPC к Besu.
// Без permission check (для HTTP версии — GetChannelWalletAndBalance в balance.go).
func (s *Service) GetChannelWallet(ctx context.Context, channelIDStr string) (string, error) {
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		return "", fmt.Errorf("%w: channel_id must be uuid", ErrInvalidInput)
	}
	w, err := s.repo.GetChannelWallet(ctx, channelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrWalletNotFound
		}
		return "", fmt.Errorf("get channel wallet: %w", err)
	}
	return w.Address, nil
}