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

// TransferInput — параметры для P2P перевода SUDA.
type TransferInput struct {
	FromUserID uuid.UUID
	ToUsername string    // без префикса @
	AmountWei  string    // десятичная строка
	Note       string    // опционально
	ChatID     uuid.UUID // uuid.Nil = обычный перевод; не Nil = transfer-in-chat
	RequestIP  string
	UserAgent  string
}

// TransferResult — результат, который handler отдаёт клиенту.
// Status подтверждения придёт асинхронно через WS-event от observer'а.
type TransferResult struct {
	TxHash      string
	FromAddress string
	ToAddress   string
	ToUserID    string
	AmountWei   string
	ChatID      string // эхо обратно (пусто если обычный перевод)
}

// Доменные ошибки, специфичные для transfer-in-chat.
var (
	ErrChatNotFound       = errors.New("chat not found")
	ErrSenderNotInChat    = errors.New("sender is not a member of this chat")
	ErrRecipientNotInChat = errors.New("recipient is not a member of this chat")
)

// Transfer выполняет P2P перевод SUDA от FromUserID к юзеру ToUsername.
//
// Если ChatID = uuid.Nil — обычный перевод.
// Если ChatID задан — transfer-in-chat: проверяется членство обоих
// в чате, observer создаст system message в этом чате.
//
// Шаги:
//  1. Парсим amount, проверяем что положительный.
//  2. Resolve ToUsername через messenger gRPC.
//  3. Проверки: получатель найден, у него есть кошелёк, это не self-transfer.
//  4. Если ChatID задан — CheckChatMembership: оба должны быть members.
//  5. Достаём кошелёк отправителя, дешифруем приватник.
//  6. ON-CHAIN balance check — читаем SudaToken.balanceOf().
//  7. Audit log (ДО broadcast'а — полный journal).
//  8. Подпись + broadcast через биндинги SudaToken.Transfer().
//  9. Запись в tx_pending (с chatID, если есть) — UI-индикатор + связь для observer.
//  10. Возврат tx_hash. Подтверждение придёт через observer/WS.
func (s *Service) Transfer(ctx context.Context, in TransferInput) (*TransferResult, error) {
	// 1. Парсинг суммы.
	amount, ok := new(big.Int).SetString(in.AmountWei, 10)
	if !ok || amount.Sign() <= 0 {
		return nil, fmt.Errorf("%w: amount_wei must be positive integer string", ErrInvalidInput)
	}

	// 2. Resolve получателя.
	resolved, err := s.messenger.ResolveUsername(ctx, in.ToUsername)
	if err != nil {
		return nil, fmt.Errorf("transfer: resolve recipient: %w", err)
	}
	if !resolved.Found {
		return nil, ErrRecipientNotFound
	}
	if resolved.WalletAddress == "" {
		return nil, ErrRecipientHasNoWallet
	}

	toUserID, err := uuid.Parse(resolved.UserID)
	if err != nil {
		return nil, fmt.Errorf("transfer: parse to_user_id: %w", err)
	}

	// 3. Кошелёк отправителя.
	senderWallet, err := s.repo.GetUserWallet(ctx, in.FromUserID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, fmt.Errorf("transfer: lookup sender: %w", err)
	}

	// Self-transfer запрещён.
	if senderWallet.Address == resolved.WalletAddress {
		return nil, ErrSelfTransfer
	}

	// 4. Если перевод в чате — проверяем что оба члены чата.
	if in.ChatID != uuid.Nil {
		if err := s.checkChatMembership(ctx, in.ChatID, in.FromUserID, toUserID); err != nil {
			return nil, err
		}
	}

	// 5. Дешифровка приватника.
	pkBytes, err := s.encryptor.Decrypt(senderWallet.KeyVersion, senderWallet.EncryptedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("transfer: decrypt sender pk: %w", err)
	}
	userSigner, err := blockchain.NewUserSigner(string(pkBytes))
	if err != nil {
		return nil, fmt.Errorf("transfer: build signer: %w", err)
	}

	// 6. On-chain balance check (защита от signed-tx-which-will-fail).
	balance, err := s.reader.SudaBalanceOf(ctx, userSigner.Address())
	if err != nil {
		return nil, fmt.Errorf("transfer: read balance: %w", err)
	}
	if balance.Cmp(amount) < 0 {
		log.Warn().
			Str("user_id", in.FromUserID.String()).
			Str("balance", balance.String()).
			Str("requested", amount.String()).
			Msg("transfer rejected: insufficient balance")
		return nil, ErrInsufficientBalance
	}

	// 7. Audit log ДО broadcast'а (по решению: полный journal).
	// Передаём пустой tx_hash — он ещё не известен.
	if err := s.repo.WriteAudit(ctx, repository.AuditEntry{
		SubjectType: wallet.SubjectUser,
		SubjectID:   in.FromUserID,
		Operation:   wallet.OpTransfer,
		TxHash:      "", // tx_hash ещё нет
		RequestIP:   in.RequestIP,
		UserAgent:   in.UserAgent,
	}); err != nil {
		log.Error().Err(err).Str("user_id", in.FromUserID.String()).Msg("audit pre-broadcast failed")
		return nil, fmt.Errorf("transfer: audit failed: %w", err)
	}

	// 8. Подготовка opts + подпись + broadcast.
	opts, release, err := s.broadcast.PrepareOpts(ctx, userSigner)
	if err != nil {
		return nil, fmt.Errorf("transfer: prepare opts: %w", err)
	}
	defer release()

	toAddr := common.HexToAddress(resolved.WalletAddress)
	tx, err := s.contracts.Token.Transfer(opts, toAddr, amount)
	if err != nil {
		return nil, fmt.Errorf("transfer: broadcast: %w", err)
	}

	txHash := tx.Hash().Hex()

	// 9. Pending для UI + связь с чатом (если есть).
	if err := s.repo.InsertPendingTransfer(ctx, txHash, in.FromUserID, in.ChatID, in.Note); err != nil {
		// Не критично — TX уже отправлена, юзер просто не увидит индикатор/system message.
		log.Warn().Err(err).Str("tx_hash", txHash).Msg("insert pending failed")
	}

	resultChatID := ""
	if in.ChatID != uuid.Nil {
		resultChatID = in.ChatID.String()
	}

	log.Info().
		Str("from_user_id", in.FromUserID.String()).
		Str("to_user_id", toUserID.String()).
		Str("amount_wei", amount.String()).
		Str("tx_hash", txHash).
		Str("chat_id", resultChatID).
		Msg("P2P transfer broadcast")

	return &TransferResult{
		TxHash:      txHash,
		FromAddress: senderWallet.Address,
		ToAddress:   resolved.WalletAddress,
		ToUserID:    toUserID.String(),
		AmountWei:   amount.String(),
		ChatID:      resultChatID,
	}, nil
}

// checkChatMembership — обёртка для проверки членства обоих юзеров в чате.
// Делает один batch gRPC к messenger.
//
// Возвращает доменную ошибку с понятным значением:
//   ErrChatNotFound        — чата нет
//   ErrSenderNotInChat     — отправитель не member
//   ErrRecipientNotInChat  — получатель не member
func (s *Service) checkChatMembership(
	ctx context.Context, chatID, fromUserID, toUserID uuid.UUID,
) error {
	result, err := s.messenger.CheckChatMembership(ctx, chatID.String(), []string{
		fromUserID.String(),
		toUserID.String(),
	})
	if err != nil {
		return fmt.Errorf("transfer: check chat membership: %w", err)
	}

	if !result.ChatExists {
		return ErrChatNotFound
	}

	// Members[0] — отправитель, Members[1] — получатель (порядок гарантирован).
	if len(result.Members) < 2 {
		return fmt.Errorf("transfer: unexpected response from messenger (got %d members)", len(result.Members))
	}

	if !result.Members[0].IsMember {
		return ErrSenderNotInChat
	}
	if !result.Members[1].IsMember {
		return ErrRecipientNotInChat
	}

	return nil
}