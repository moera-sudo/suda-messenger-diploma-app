package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/blockchain"
)

// ChargeChannelSubscription списывает с юзера цену платной подписки на PUBLIC-канал
// (tx_gating_rules.subscription_price_wei) и переводит её в treasury канала.
// Подписывается ключом ПОЛЬЗОВАТЕЛЯ. Возвращает tx_hash и списанную цену.
//
// Ошибки: ErrNotPaidChannel (нет правила или price<=0), ErrWalletNotFound,
// ErrInsufficientBalance. Вызывается из gRPC (Messenger.channel.Subscribe).
func (s *Service) ChargeChannelSubscription(
	ctx context.Context, userIDStr, channelIDStr string,
) (txHash, priceWei string, err error) {
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		return "", "", fmt.Errorf("%w: channel_id must be uuid", ErrInvalidInput)
	}
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return "", "", fmt.Errorf("%w: user_id must be uuid", ErrInvalidInput)
	}

	// 1. Правило + цена подписки.
	rule, err := s.repo.GetGatingRule(ctx, channelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", ErrNotPaidChannel
		}
		return "", "", fmt.Errorf("subscribe: read rule: %w", err)
	}
	price := rule.SubscriptionPrice
	if price == nil || price.Sign() <= 0 {
		return "", "", ErrNotPaidChannel
	}

	// 2. Treasury канала (получатель оплаты).
	treasury, err := s.repo.GetChannelWallet(ctx, channelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", ErrWalletNotFound
		}
		return "", "", fmt.Errorf("subscribe: lookup channel wallet: %w", err)
	}

	// 3. Кошелёк юзера + дешифровка ключа.
	userWallet, err := s.repo.GetUserWallet(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", "", ErrWalletNotFound
		}
		return "", "", fmt.Errorf("subscribe: lookup user wallet: %w", err)
	}
	pkBytes, err := s.encryptor.Decrypt(userWallet.KeyVersion, userWallet.EncryptedPrivateKey)
	if err != nil {
		return "", "", fmt.Errorf("subscribe: decrypt pk: %w", err)
	}
	signer, err := blockchain.NewUserSigner(string(pkBytes))
	if err != nil {
		return "", "", fmt.Errorf("subscribe: build signer: %w", err)
	}

	// 4. On-chain balance check.
	balance, err := s.reader.SudaBalanceOf(ctx, signer.Address())
	if err != nil {
		return "", "", fmt.Errorf("subscribe: read balance: %w", err)
	}
	if balance.Cmp(price) < 0 {
		return "", "", ErrInsufficientBalance
	}

	// 5. Audit ДО broadcast'а (подписка = перевод в казну канала).
	if err := s.repo.WriteAudit(ctx, repository.AuditEntry{
		SubjectType: wallet.SubjectUser,
		SubjectID:   userID,
		Operation:   wallet.OpDonation,
	}); err != nil {
		return "", "", fmt.Errorf("subscribe: audit failed: %w", err)
	}

	// 6. Подпись ключом юзера + broadcast SudaToken.Transfer(treasury, price).
	opts, release, err := s.broadcast.PrepareOpts(ctx, signer)
	if err != nil {
		return "", "", fmt.Errorf("subscribe: prepare opts: %w", err)
	}
	defer release()

	tx, err := s.contracts.Token.Transfer(opts, common.HexToAddress(treasury.Address), price)
	if err != nil {
		return "", "", fmt.Errorf("subscribe: broadcast: %w", err)
	}
	txHash = tx.Hash().Hex()

	// 7. Pending (индикатор/observer); related_chat_id = канал.
	if err := s.repo.InsertPending(ctx, txHash, userID, "CHANNEL_SUBSCRIBE", channelID); err != nil {
		log.Warn().Err(err).Str("tx_hash", txHash).Msg("subscribe: insert pending failed")
	}

	log.Info().
		Str("user_id", userID.String()).
		Str("channel_id", channelID.String()).
		Str("price_wei", price.String()).
		Str("tx_hash", txHash).
		Msg("channel subscription charged")

	return txHash, price.String(), nil
}
