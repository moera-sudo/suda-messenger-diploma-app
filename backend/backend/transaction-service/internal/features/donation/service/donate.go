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
	walletrepo "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/blockchain"
)

// DonationInput — параметры POST /donate.
//
// Ровно одно из ToUsername / ToChannelID должно быть задано:
//   - ToUsername  → P2P-донат конкретному юзеру;
//   - ToChannelID → донат на treasury канала.
type DonationInput struct {
	FromUserID  uuid.UUID
	ToUsername  string    // P2P-донат (без префикса @)
	ToChannelID uuid.UUID // донат каналу; uuid.Nil если P2P
	AmountWei   string    // десятичная строка wei
	Message     string    // опционально, <= 200 символов
	ChatID      uuid.UUID // опц.: чат, в котором появится system message (для P2P)
	RequestIP   string
	UserAgent   string
}

// DonationResult — async-ответ: tx_hash сразу, подтверждение через observer/WS.
type DonationResult struct {
	TxHash      string
	FromAddress string
	ToAddress   string
	AmountWei   string
	IsChannel   bool   // true — донат каналу, false — P2P
	ToUserID    string // заполнен только для P2P
	ToChannelID string // заполнен только для donation каналу
}

// Donate выполняет донат SUDA. Flow повторяет wallet.Transfer, но pending
// помечается как DONATION + сохраняется donation_message.
//
// Шаги:
//  1. Парсинг суммы.
//  2. Валидация: ровно один получатель (username XOR channel).
//  3. Резолв адреса получателя.
//  4. Кошелёк отправителя + проверка self-donation.
//  5. (P2P + ChatID) проверка членства обоих в чате.
//  6. Дешифровка приватника отправителя.
//  7. On-chain balance check.
//  8. Audit ДО broadcast'а.
//  9. Подпись + broadcast SudaToken.Transfer.
//  10. tx_pending с expected_type=DONATION + donation_message.
func (s *Service) Donate(ctx context.Context, in DonationInput) (*DonationResult, error) {
	// 1. Сумма.
	amount, ok := new(big.Int).SetString(in.AmountWei, 10)
	if !ok || amount.Sign() <= 0 {
		return nil, fmt.Errorf("%w: amount_wei must be positive integer string", ErrInvalidInput)
	}

	// 2. Ровно один получатель.
	isChannel := in.ToChannelID != uuid.Nil
	hasUsername := in.ToUsername != ""
	if isChannel == hasUsername {
		return nil, fmt.Errorf("%w: exactly one of to_username / to_channel_id must be set", ErrInvalidInput)
	}

	// 3. Резолв получателя.
	var toAddress string
	var toUserID uuid.UUID
	if isChannel {
		w, err := s.repo.GetChannelWallet(ctx, in.ToChannelID)
		if err != nil {
			if errors.Is(err, walletrepo.ErrNotFound) {
				return nil, ErrChannelWalletNotFound
			}
			return nil, fmt.Errorf("donate: lookup channel wallet: %w", err)
		}
		toAddress = w.Address
	} else {
		resolved, err := s.messenger.ResolveUsername(ctx, in.ToUsername)
		if err != nil {
			return nil, fmt.Errorf("donate: resolve recipient: %w", err)
		}
		if !resolved.Found {
			return nil, ErrRecipientNotFound
		}
		if resolved.WalletAddress == "" {
			return nil, ErrRecipientHasNoWallet
		}
		toAddress = resolved.WalletAddress
		toUserID, err = uuid.Parse(resolved.UserID)
		if err != nil {
			return nil, fmt.Errorf("donate: parse to_user_id: %w", err)
		}
	}

	// 4. Кошелёк отправителя.
	senderWallet, err := s.repo.GetUserWallet(ctx, in.FromUserID)
	if err != nil {
		if errors.Is(err, walletrepo.ErrNotFound) {
			return nil, ErrWalletNotFound
		}
		return nil, fmt.Errorf("donate: lookup sender: %w", err)
	}
	if senderWallet.Address == toAddress {
		return nil, ErrSelfDonation
	}

	// 5. P2P-донат в чате — оба должны быть members.
	if !isChannel && in.ChatID != uuid.Nil {
		if err := s.checkChatMembership(ctx, in.ChatID, in.FromUserID, toUserID); err != nil {
			return nil, err
		}
	}

	// 6. Дешифровка приватника.
	pkBytes, err := s.encryptor.Decrypt(senderWallet.KeyVersion, senderWallet.EncryptedPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("donate: decrypt sender pk: %w", err)
	}
	userSigner, err := blockchain.NewUserSigner(string(pkBytes))
	if err != nil {
		return nil, fmt.Errorf("donate: build signer: %w", err)
	}

	// 7. On-chain balance check.
	balance, err := s.reader.SudaBalanceOf(ctx, userSigner.Address())
	if err != nil {
		return nil, fmt.Errorf("donate: read balance: %w", err)
	}
	if balance.Cmp(amount) < 0 {
		log.Warn().
			Str("user_id", in.FromUserID.String()).
			Str("balance", balance.String()).
			Str("requested", amount.String()).
			Msg("donation rejected: insufficient balance")
		return nil, ErrInsufficientBalance
	}

	// 8. Audit ДО broadcast'а.
	if err := s.repo.WriteAudit(ctx, walletrepo.AuditEntry{
		SubjectType: wallet.SubjectUser,
		SubjectID:   in.FromUserID,
		Operation:   wallet.OpDonation,
		TxHash:      "",
		RequestIP:   in.RequestIP,
		UserAgent:   in.UserAgent,
	}); err != nil {
		log.Error().Err(err).Str("user_id", in.FromUserID.String()).Msg("donation audit pre-broadcast failed")
		return nil, fmt.Errorf("donate: audit failed: %w", err)
	}

	// 9. Подпись + broadcast (один высокоуровневый вызов — мокается в тестах).
	txHash, err := s.sender.SendTokenTransfer(ctx, userSigner, common.HexToAddress(toAddress), amount)
	if err != nil {
		return nil, fmt.Errorf("donate: broadcast: %w", err)
	}

	// 10. Pending с DONATION-метаданными.
	// systemMessageChatID — где появится system message:
	//   - канал: сам channelID (канал == чат);
	//   - P2P:   переданный ChatID (может быть uuid.Nil — тогда без system message).
	systemMessageChatID := in.ChatID
	if isChannel {
		systemMessageChatID = in.ToChannelID
	}
	if err := s.repo.InsertPendingDonation(ctx, txHash, in.FromUserID, systemMessageChatID, in.Message); err != nil {
		// TX уже broadcast'нута — операцию не отменяем. Но без pending observer
		// проиндексирует её как обычный P2P_TRANSFER: не будет ни записи в
		// tx_donations, ни system message в чате. Это деградация — логируем Error.
		log.Error().Err(err).Str("tx_hash", txHash).
			Msg("insert pending donation failed: donation will be indexed as a plain transfer")
	}

	log.Info().
		Str("from_user_id", in.FromUserID.String()).
		Str("to_address", toAddress).
		Bool("is_channel", isChannel).
		Str("amount_wei", amount.String()).
		Str("tx_hash", txHash).
		Msg("donation broadcast")

	result := &DonationResult{
		TxHash:      txHash,
		FromAddress: senderWallet.Address,
		ToAddress:   toAddress,
		AmountWei:   amount.String(),
		IsChannel:   isChannel,
	}
	if isChannel {
		result.ToChannelID = in.ToChannelID.String()
	} else {
		result.ToUserID = toUserID.String()
	}
	return result, nil
}

// checkChatMembership проверяет что отправитель и получатель — оба члены чата.
// Делает один batch gRPC к messenger.
func (s *Service) checkChatMembership(ctx context.Context, chatID, fromUserID, toUserID uuid.UUID) error {
	result, err := s.messenger.CheckChatMembership(ctx, chatID.String(), []string{
		fromUserID.String(),
		toUserID.String(),
	})
	if err != nil {
		return fmt.Errorf("donate: check chat membership: %w", err)
	}
	if !result.ChatExists {
		return ErrChatNotFound
	}
	if len(result.Members) < 2 {
		return fmt.Errorf("donate: unexpected messenger response (%d members)", len(result.Members))
	}
	if !result.Members[0].IsMember {
		return ErrSenderNotInChat
	}
	if !result.Members[1].IsMember {
		return ErrRecipientNotInChat
	}
	return nil
}
