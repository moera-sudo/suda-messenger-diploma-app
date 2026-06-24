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
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"
)

// WithdrawResult — результат вывода средств из казны канала (async-ответ).
type WithdrawResult struct {
	TxHash      string
	FromAddress string // адрес казны
	ToAddress   string // кошелёк владельца
	AmountWei   string
}

// WithdrawTreasury выводит SUDA из казны канала на кошелёк его владельца (= автора).
// Доступно только OWNER канала; получатель — кошелёк вызывающего (он же OWNER).
// amountWei — сумма к выводу (частичный или полный вывод).
//
// Подпись делается ключом КАЗНЫ (tx_channel_wallets.encrypted_private_key).
// Возвращает tx_hash сразу; подтверждение придёт асинхронно через observer.
func (s *Service) WithdrawTreasury(
	ctx context.Context, callerUserID, channelID uuid.UUID, amountWei, requestIP, userAgent string,
) (*WithdrawResult, error) {
	// 1. Сумма.
	amount, ok := new(big.Int).SetString(amountWei, 10)
	if !ok || amount.Sign() <= 0 {
		return nil, fmt.Errorf("%w: amount_wei must be positive integer string", ErrInvalidInput)
	}

	// 2. Права: вывести может только OWNER канала.
	perm, err := s.messenger.CheckChannelPermission(
		ctx, callerUserID.String(), channelID.String(), platgrpc.PermissionOwner,
	)
	if err != nil {
		return nil, fmt.Errorf("withdraw: permission check: %w", err)
	}
	if !perm.Granted {
		return nil, ErrForbidden
	}

	// 3. Казна канала (источник средств + ключ для подписи).
	treasury, err := s.repo.GetChannelWallet(ctx, channelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, fmt.Errorf("withdraw: lookup channel wallet: %w", err)
	}

	// 4. Кошелёк владельца (получатель) = кошелёк вызывающего.
	ownerWallet, err := s.repo.GetUserWallet(ctx, callerUserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, fmt.Errorf("withdraw: lookup owner wallet: %w", err)
	}

	// 5. Дешифровка ключа КАЗНЫ + signer.
	pkBytes, err := s.encryptor.Decrypt(treasury.KeyVersion, treasury.EncryptedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("withdraw: decrypt treasury pk: %w", err)
	}
	treasurySigner, err := blockchain.NewUserSigner(string(pkBytes))
	if err != nil {
		return nil, fmt.Errorf("withdraw: build signer: %w", err)
	}

	// 6. On-chain баланс казны.
	balance, err := s.reader.SudaBalanceOf(ctx, treasurySigner.Address())
	if err != nil {
		return nil, fmt.Errorf("withdraw: read treasury balance: %w", err)
	}
	if balance.Cmp(amount) < 0 {
		log.Warn().
			Str("channel_id", channelID.String()).
			Str("balance", balance.String()).
			Str("requested", amount.String()).
			Msg("withdraw rejected: insufficient treasury balance")
		return nil, ErrInsufficientBalance
	}

	// 7. Audit ДО broadcast'а (subject = CHANNEL).
	if err := s.repo.WriteAudit(ctx, repository.AuditEntry{
		SubjectType: wallet.SubjectChannel,
		SubjectID:   channelID,
		Operation:   wallet.OpChannelWithdraw,
		RequestIP:   requestIP,
		UserAgent:   userAgent,
	}); err != nil {
		log.Error().Err(err).Str("channel_id", channelID.String()).Msg("withdraw audit pre-broadcast failed")
		return nil, fmt.Errorf("withdraw: audit failed: %w", err)
	}

	// 8. Подпись ключом казны + broadcast SudaToken.Transfer(owner, amount).
	opts, release, err := s.broadcast.PrepareOpts(ctx, treasurySigner)
	if err != nil {
		return nil, fmt.Errorf("withdraw: prepare opts: %w", err)
	}
	defer release()

	tx, err := s.contracts.Token.Transfer(opts, common.HexToAddress(ownerWallet.Address), amount)
	if err != nil {
		return nil, fmt.Errorf("withdraw: broadcast: %w", err)
	}
	txHash := tx.Hash().Hex()

	// 9. Pending: from_user_id = владелец, related_chat_id = канал.
	if err := s.repo.InsertPending(ctx, txHash, callerUserID, "CHANNEL_WITHDRAW", channelID); err != nil {
		log.Warn().Err(err).Str("tx_hash", txHash).Msg("withdraw: insert pending failed")
	}

	log.Info().
		Str("channel_id", channelID.String()).
		Str("owner_id", callerUserID.String()).
		Str("amount_wei", amount.String()).
		Str("tx_hash", txHash).
		Msg("treasury withdraw broadcast")

	return &WithdrawResult{
		TxHash:      txHash,
		FromAddress: treasury.Address,
		ToAddress:   ownerWallet.Address,
		AmountWei:   amount.String(),
	}, nil
}
