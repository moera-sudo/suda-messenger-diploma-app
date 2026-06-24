// Package handlers — конкретные event handler'ы observer'а.
// Каждый файл соответствует одному контракту.
package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/transaction-service/contracts/sudatoken"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/donation"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase"
	purchaserepo "github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/blockchain"
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"
)

// Event-types, которые мы шлём через NotifyUserEvent.
// Должны быть согласованы с messenger-service.
const (
	EventSudaReceived       = "SUDA_RECEIVED"
	EventSudaSent           = "SUDA_SENT"
	EventPurchaseCompleted  = "PURCHASE_COMPLETED"
)

// Типы транзакций для tx_transactions.type.
const (
	TxTypeP2PTransfer = "P2P_TRANSFER"
	TxTypeMint        = "MINT"
)

// Адрес 0x0...0 — индикатор mint в ERC-20.
var zeroAddress = common.Address{}

// ────────────────────────────────────────────────────────────
//  Interfaces for observer dependencies (mockable in tests).
// ────────────────────────────────────────────────────────────

type observerWalletRepo interface {
	GetPendingMeta(ctx context.Context, txHash string) (*repository.PendingMeta, error)
	InsertTransaction(ctx context.Context, p repository.TxInsertParams) (bool, error)
	DeletePending(ctx context.Context, txHash string) error
	InsertDonation(ctx context.Context, p repository.DonationInsertParams) error
	FindUserByAddress(ctx context.Context, address string) (*uuid.UUID, error)
	FindChannelByAddress(ctx context.Context, address string) (*uuid.UUID, error)
}

type observerPurchaseRepo interface {
	FindByTxHash(ctx context.Context, txHash string) (*purchaserepo.TxHashMatch, error)
}

type observerMessenger interface {
	NotifyUserEvent(ctx context.Context, req platgrpc.NotifyEventRequest) (*platgrpc.NotifyEventResult, error)
	GetUsersByIDs(ctx context.Context, userIDs []string) (map[string]platgrpc.UserBrief, error)
}

// SudaTokenHandler — обрабатывает Transfer events контракта SudaToken.
type SudaTokenHandler struct {
	repo         observerWalletRepo
	purchaseRepo observerPurchaseRepo   // nil → purchase-lookup skipped
	contracts    *blockchain.Contracts  // for FilterTransfer + contract-address filter
	messenger    observerMessenger
	isContract   func(addr common.Address) bool                                   // injectable for tests
	blockTimeFn  func(ctx context.Context, blockNumber uint64) (time.Time, error) // injectable for tests
}

// SudaTokenDeps — зависимости handler'а.
type SudaTokenDeps struct {
	Repo            *repository.Repository
	PurchaseRepo    *purchaserepo.Repository
	Chain           *blockchain.Client
	Contracts       *blockchain.Contracts
	MessengerClient *platgrpc.MessengerClient
}

func NewSudaTokenHandler(d SudaTokenDeps) *SudaTokenHandler {
	h := &SudaTokenHandler{
		repo:         d.Repo,
		purchaseRepo: d.PurchaseRepo,
		contracts:    d.Contracts,
		messenger:    d.MessengerClient,
	}
	h.isContract = func(addr common.Address) bool {
		return addr == d.Contracts.MarketplaceAddr ||
			addr == d.Contracts.EscrowAddr ||
			addr == d.Contracts.FundraisingAddr
	}
	chain := d.Chain
	h.blockTimeFn = func(ctx context.Context, blockNumber uint64) (time.Time, error) {
		header, err := chain.Eth().HeaderByNumber(ctx, new(big.Int).SetUint64(blockNumber))
		if err != nil {
			return time.Time{}, fmt.Errorf("header by number %d: %w", blockNumber, err)
		}
		return time.Unix(int64(header.Time), 0).UTC(), nil
	}
	return h
}

// NewSudaTokenHandlerForTest creates a testable handler with interface-typed deps.
func NewSudaTokenHandlerForTest(
	repo observerWalletRepo,
	purchaseRepo observerPurchaseRepo,
	messenger observerMessenger,
	blockTimeFn func(ctx context.Context, blockNumber uint64) (time.Time, error),
) *SudaTokenHandler {
	return &SudaTokenHandler{
		repo:         repo,
		purchaseRepo: purchaseRepo,
		messenger:    messenger,
		isContract:   func(common.Address) bool { return false },
		blockTimeFn:  blockTimeFn,
	}
}

// HandleLogs обрабатывает все Transfer events в окне блоков [fromBlock, toBlock].
//
// Использует FilterTransfer от биндингов — он сам делает eth_getLogs с
// правильным topic и разбирает данные в типизированную структуру.
func (h *SudaTokenHandler) HandleLogs(ctx context.Context, fromBlock, toBlock uint64) error {
	opts := &bind.FilterOpts{
		Start:   fromBlock,
		End:     &toBlock,
		Context: ctx,
	}

	iter, err := h.contracts.Token.FilterTransfer(opts, nil, nil)
	if err != nil {
		return fmt.Errorf("filter transfer: %w", err)
	}
	defer iter.Close()

	count := 0
	for iter.Next() {
		ev := iter.Event
		if err := h.handleTransfer(ctx, ev); err != nil {
			log.Error().
				Err(err).
				Str("tx_hash", ev.Raw.TxHash.Hex()).
				Uint64("block", ev.Raw.BlockNumber).
				Msg("handle SudaToken Transfer failed")
		} else {
			count++
		}
	}
	if err := iter.Error(); err != nil {
		return fmt.Errorf("iter error: %w", err)
	}

	if count > 0 {
		log.Debug().
			Uint64("from_block", fromBlock).
			Uint64("to_block", toBlock).
			Int("events", count).
			Msg("processed SudaToken Transfer events")
	}

	return nil
}

// handleTransfer обрабатывает один Transfer event.
func (h *SudaTokenHandler) handleTransfer(ctx context.Context, ev *sudatoken.SudaTokenTransfer) error {
	from := ev.From
	to := ev.To
	amount := ev.Value
	txHash := ev.Raw.TxHash.Hex()
	logIndex := int32(ev.Raw.Index)
	blockNumber := int64(ev.Raw.BlockNumber)

	// Игнорируем переводы к/от наших contract-сервисов — обработают свои handlers.
	if h.isOurContract(from) || h.isOurContract(to) {
		return nil
	}

	// Получаем timestamp блока.
	blockTime, err := h.getBlockTime(ctx, ev.Raw.BlockNumber)
	if err != nil {
		return fmt.Errorf("get block time: %w", err)
	}

	// Резолвим адреса в user_id'ы (могут быть оба nil — внешний адрес).
	fromUserID, _ := h.lookupUserByAddress(ctx, from)
	toUserID, _ := h.lookupUserByAddress(ctx, to)

	// Определяем тип транзакции.
	var txType string
	switch {
	case from == zeroAddress:
		txType = TxTypeMint
	default:
		txType = TxTypeP2PTransfer
	}

	// Читаем pending-метаданные ДО удаления pending.
	//   - related_chat_id  → transfer-in-chat / donation в чате (нужен system message);
	//   - expected_type    → если "DONATION", делаем доп. обработку (см. ниже);
	//   - donation_message → текст доната для system message.
	var pendingMeta *repository.PendingMeta
	if m, err := h.repo.GetPendingMeta(ctx, txHash); err == nil {
		pendingMeta = m
	}
	// errors.Is(err, repository.ErrNotFound) — норма (не наш перевод или welcome bonus)

	chatID := uuid.Nil
	isDonation := false
	if pendingMeta != nil {
		chatID = pendingMeta.ChatID
		isDonation = pendingMeta.ExpectedType == "DONATION"
	}

	// related_entity для истории в tx_transactions.
	var relatedEntityType *string
	var relatedEntityID *uuid.UUID
	if chatID != uuid.Nil {
		rt := "CHAT_MESSAGE"
		relatedEntityType = &rt
		cid := chatID
		relatedEntityID = &cid
	}

	// DONATION: переопределяем тип транзакции.
	if isDonation {
		txType = donation.TxTypeDonation
	}

	// Проверяем, не является ли этот Transfer покупкой SUDA через /purchase.
	// Если да — переопределяем тип и привязываем к purchase в related_entity.
	purchaseMatch := h.lookupPurchase(ctx, txHash)
	if purchaseMatch != nil {
		txType = purchase.TxTypePurchase
		rt := "PURCHASE"
		relatedEntityType = &rt
		pid := purchaseMatch.ID
		relatedEntityID = &pid
	}

	// Вставляем запись. inserted=false → событие уже обработано в прошлом
	// прогоне observer'а (reprocess блока). В этом случае нельзя повторять
	// InsertDonation и WS-уведомления — иначе дубли в tx_donations и
	// двойные system message в чате.
	inserted, err := h.repo.InsertTransaction(ctx, repository.TxInsertParams{
		TxHash:            txHash,
		LogIndex:          logIndex,
		Type:              txType,
		Status:            "CONFIRMED",
		FromUserID:        fromUserID,
		ToUserID:          toUserID,
		FromAddress:       from.Hex(),
		ToAddress:         to.Hex(),
		Amount:            amount,
		BlockNumber:       blockNumber,
		ConfirmedAt:       blockTime,
		RelatedEntityType: relatedEntityType,
		RelatedEntityID:   relatedEntityID,
	})
	if err != nil {
		return fmt.Errorf("insert tx: %w", err)
	}
	if !inserted {
		// Уже обработан ранее. Чистим pending (идемпотентно) и выходим
		// без InsertDonation и без повторных уведомлений.
		if err := h.repo.DeletePending(ctx, txHash); err != nil {
			log.Warn().Err(err).Str("tx_hash", txHash).Msg("delete pending failed")
		}
		log.Debug().
			Str("tx_hash", txHash).
			Int32("log_index", logIndex).
			Msg("transfer event already processed, skipping donation/notifications")
		return nil
	}

	// DONATION: пишем плоскую строку в tx_donations (история для treasury-статистики).
	if isDonation {
		toChannelID, _ := h.lookupChannelByAddress(ctx, to)
		if err := h.repo.InsertDonation(ctx, repository.DonationInsertParams{
			ChannelID:   toChannelID,
			ToUserID:    toUserID, // nil если получатель — канал
			FromUserID:  derefOrZero(fromUserID),
			FromAddress: from.Hex(),
			ToAddress:   to.Hex(),
			Amount:      amount,
			Message:     pendingMeta.DonationMessage,
			TxHash:      txHash,
		}); err != nil {
			log.Warn().Err(err).Str("tx_hash", txHash).Msg("insert donation row failed")
		}
	}

	// Удаляем pending (best-effort).
	if err := h.repo.DeletePending(ctx, txHash); err != nil {
		log.Warn().Err(err).Str("tx_hash", txHash).Msg("delete pending failed")
	}

	log.Info().
		Str("tx_hash", txHash).
		Str("type", txType).
		Str("from", from.Hex()).
		Str("to", to.Hex()).
		Str("amount_wei", amount.String()).
		Int64("block", blockNumber).
		Str("chat_id", chatID.String()).
		Msg("SudaToken Transfer indexed")

	// WS-уведомления (best-effort). Каждый кейс — свой event_type.
	switch {
	case purchaseMatch != nil:
		h.notifyPurchase(ctx, purchaseMatch, to.Hex(), amount, txHash)
	case isDonation:
		message := ""
		if pendingMeta != nil {
			message = pendingMeta.DonationMessage
		}
		h.notifyDonation(ctx, fromUserID, toUserID, from.Hex(), to.Hex(), amount, txHash, chatID, message)
	default:
		note := ""
		if pendingMeta != nil {
			note = pendingMeta.DonationMessage
		}
		h.notifyTransfer(ctx, fromUserID, toUserID, from.Hex(), to.Hex(), amount, txHash, chatID, note)
	}

	return nil
}

// lookupChannelByAddress — обёртка с graceful nil-возвратом для адресов,
// которых нет в tx_channel_wallets (значит получатель — обычный юзер).
func (h *SudaTokenHandler) lookupChannelByAddress(ctx context.Context, addr common.Address) (*uuid.UUID, error) {
	if id, err := h.repo.FindChannelByAddress(ctx, addr.Hex()); err == nil {
		return id, nil
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	return nil, nil
}

// derefOrZero разыменовывает *uuid.UUID или возвращает uuid.Nil.
func derefOrZero(id *uuid.UUID) uuid.UUID {
	if id == nil {
		return uuid.Nil
	}
	return *id
}

// lookupPurchase возвращает match если этот tx_hash был записан как purchase'овый
// transfer от treasury, иначе nil.
func (h *SudaTokenHandler) lookupPurchase(ctx context.Context, txHash string) *purchaserepo.TxHashMatch {
	if h.purchaseRepo == nil {
		return nil
	}
	m, err := h.purchaseRepo.FindByTxHash(ctx, txHash)
	if err != nil {
		if !errors.Is(err, purchaserepo.ErrNotFound) {
			log.Warn().Err(err).Str("tx_hash", txHash).Msg("purchase lookup failed")
		}
		return nil
	}
	return m
}

// notifyPurchase шлёт PURCHASE_COMPLETED только покупателю.
// SUDA_RECEIVED для этого Transfer'а намеренно не шлём, чтобы UI не показал
// двойное уведомление "вам пришли SUDA" + "покупка завершена".
func (h *SudaTokenHandler) notifyPurchase(
	ctx context.Context,
	match *purchaserepo.TxHashMatch,
	toAddr string,
	amount *big.Int,
	txHash string,
) {
	payload := map[string]any{
		"tx_hash":         txHash,
		"purchase_id":     match.ID.String(),
		"package_code":    match.PackageCode,
		"amount_wei":      amount.String(),
		"to_address":      toAddr,
		"user_id":         match.UserID.String(),
	}
	payloadJSON, _ := json.Marshal(payload)

	_, err := h.messenger.NotifyUserEvent(ctx, platgrpc.NotifyEventRequest{
		UserID:      match.UserID.String(),
		EventType:   EventPurchaseCompleted,
		PayloadJSON: string(payloadJSON),
	})
	if err != nil {
		log.Warn().Err(err).
			Str("tx_hash", txHash).
			Str("user_id", match.UserID.String()).
			Msg("notify PURCHASE_COMPLETED failed")
	}
}

// notifyDonation шлёт WS-уведомления о донате.
//
//   - DONATION_SENT отправителю. Несёт target_chat_id (если есть), чтобы
//     messenger создал system message типа DONATION в чате/канале — его
//     увидят все участники (для канала это все подписчики).
//   - DONATION_RECEIVED получателю-юзеру (только P2P-донат, toUserID != nil).
//     Без target_chat_id, чтобы не дублировать system message.
//
// Для доната каналу получатель — treasury (не юзер), поэтому отдельного
// DONATION_RECEIVED нет; system message в канале служит уведомлением для всех.
func (h *SudaTokenHandler) notifyDonation(
	ctx context.Context,
	fromUserID, toUserID *uuid.UUID,
	fromAddr, toAddr string,
	amount *big.Int,
	txHash string,
	chatID uuid.UUID,
	message string,
) {
	chatIDStr := ""
	systemMessageType := ""
	if chatID != uuid.Nil {
		chatIDStr = chatID.String()
		systemMessageType = donation.SystemMessageDonation
	}

	// Payload одинаков для обоих уведомлений — собираем один раз.
	payload := map[string]any{
		"tx_hash":      txHash,
		"from_address": fromAddr,
		"to_address":   toAddr,
		"amount_wei":   amount.String(),
	}
	if message != "" {
		payload["message"] = message
	}
	if fromUserID != nil {
		payload["from_user_id"] = fromUserID.String()
	}
	if toUserID != nil {
		payload["to_user_id"] = toUserID.String()
	}
	if chatID != uuid.Nil {
		payload["chat_id"] = chatID.String()
	}
	payloadJSON, _ := json.Marshal(payload)

	// DONATION_SENT отправителю — несёт target_chat_id (создаёт system message).
	if fromUserID != nil {
		_, err := h.messenger.NotifyUserEvent(ctx, platgrpc.NotifyEventRequest{
			UserID:            fromUserID.String(),
			EventType:         donation.EventDonationSent,
			PayloadJSON:       string(payloadJSON),
			TargetChatID:      chatIDStr,
			SystemMessageType: systemMessageType,
		})
		if err != nil {
			log.Warn().Err(err).
				Str("tx_hash", txHash).
				Str("from_user_id", fromUserID.String()).
				Msg("notify DONATION_SENT failed")
		}
	}

	// DONATION_RECEIVED получателю-юзеру (P2P) — без target_chat_id.
	if toUserID != nil {
		_, err := h.messenger.NotifyUserEvent(ctx, platgrpc.NotifyEventRequest{
			UserID:      toUserID.String(),
			EventType:   donation.EventDonationReceived,
			PayloadJSON: string(payloadJSON),
		})
		if err != nil {
			log.Warn().Err(err).
				Str("tx_hash", txHash).
				Str("to_user_id", toUserID.String()).
				Msg("notify DONATION_RECEIVED failed")
		}
	}
}

// notifyTransfer шлёт WS-event получателю (SUDA_RECEIVED) и отправителю (SUDA_SENT).
//
// Если chatID != uuid.Nil — это transfer-in-chat, и system message создаётся
// ровно ОДИН раз (через первое уведомление, у которого chat target указан).
// Логика: указываем target_chat_id только в SUDA_RECEIVED, чтобы система не
// создала два system message от двух событий.
func (h *SudaTokenHandler) notifyTransfer(
	ctx context.Context,
	fromUserID, toUserID *uuid.UUID,
	fromAddr, toAddr string,
	amount *big.Int,
	txHash string,
	chatID uuid.UUID,
	note string,
) {
	chatIDStr := ""
	systemMessageType := ""
	if chatID != uuid.Nil {
		chatIDStr = chatID.String()
		systemMessageType = "SUDA_TRANSFER"
	}

	// Резолвим display name'ы обоих участников одним batch-вызовом, чтобы
	// UI мог показать «Имя → Имя», а не сырые UUID/адреса. Сумму UI форматирует
	// сам из amount_wei — бэкенд не занимается presentation.
	fromName, toName := h.resolveTransferNames(ctx, fromUserID, toUserID)

	// SUDA_RECEIVED — получателю. Несёт target_chat_id (если есть),
	// чтобы messenger создал system message в чате.
	if toUserID != nil {
		payload := map[string]any{
			"tx_hash":      txHash,
			"from_address": fromAddr,
			"to_address":   toAddr,
			"amount_wei":   amount.String(),
		}
		if fromUserID != nil {
			payload["from_user_id"] = fromUserID.String()
		}
		payload["to_user_id"] = toUserID.String()
		if fromName != "" {
			payload["from_display_name"] = fromName
		}
		if toName != "" {
			payload["to_display_name"] = toName
		}
		if note != "" {
			payload["comment"] = note
		}
		if chatID != uuid.Nil {
			payload["chat_id"] = chatID.String()
		}
		payloadJSON, _ := json.Marshal(payload)

		_, err := h.messenger.NotifyUserEvent(ctx, platgrpc.NotifyEventRequest{
			UserID:            toUserID.String(),
			EventType:         EventSudaReceived,
			PayloadJSON:       string(payloadJSON),
			TargetChatID:      chatIDStr,        // создаст system message если задан
			SystemMessageType: systemMessageType,
		})
		if err != nil {
			log.Warn().Err(err).
				Str("tx_hash", txHash).
				Str("to_user_id", toUserID.String()).
				Msg("notify SUDA_RECEIVED failed")
		}
	}

	// SUDA_SENT — отправителю. БЕЗ target_chat_id, чтобы не создавать
	// system message повторно (он уже создан вместе с SUDA_RECEIVED).
	if fromUserID != nil {
		payload := map[string]any{
			"tx_hash":    txHash,
			"to_address": toAddr,
			"amount_wei": amount.String(),
		}
		if toUserID != nil {
			payload["to_user_id"] = toUserID.String()
		}
		payload["from_user_id"] = fromUserID.String()
		if fromName != "" {
			payload["from_display_name"] = fromName
		}
		if toName != "" {
			payload["to_display_name"] = toName
		}
		if note != "" {
			payload["comment"] = note
		}
		if chatID != uuid.Nil {
			payload["chat_id"] = chatID.String()
		}
		payloadJSON, _ := json.Marshal(payload)

		_, err := h.messenger.NotifyUserEvent(ctx, platgrpc.NotifyEventRequest{
			UserID:      fromUserID.String(),
			EventType:   EventSudaSent,
			PayloadJSON: string(payloadJSON),
			// TargetChatID и SystemMessageType — пусто.
		})
		if err != nil {
			log.Warn().Err(err).
				Str("tx_hash", txHash).
				Str("from_user_id", fromUserID.String()).
				Msg("notify SUDA_SENT failed")
		}
	}
}

// resolveTransferNames batch-резолвит display name'ы отправителя и получателя.
// Любой из id может быть nil (внешний адрес) — тогда соответствующее имя пустое.
// Ошибки резолва не фатальны: возвращаем то, что удалось получить.
func (h *SudaTokenHandler) resolveTransferNames(
	ctx context.Context, fromUserID, toUserID *uuid.UUID,
) (fromName, toName string) {
	ids := make([]string, 0, 2)
	if fromUserID != nil {
		ids = append(ids, fromUserID.String())
	}
	if toUserID != nil {
		ids = append(ids, toUserID.String())
	}
	if len(ids) == 0 {
		return "", ""
	}

	users, err := h.messenger.GetUsersByIDs(ctx, ids)
	if err != nil {
		log.Warn().Err(err).Msg("notifyTransfer: resolve display names failed")
		return "", ""
	}

	pick := func(id *uuid.UUID) string {
		if id == nil {
			return ""
		}
		if u, ok := users[id.String()]; ok {
			if u.DisplayName != "" {
				return u.DisplayName
			}
			return u.Username
		}
		return ""
	}
	return pick(fromUserID), pick(toUserID)
}

// lookupUserByAddress — обёртка с graceful nil-возвратом для адресов,
// которых нет в наших кошельках.
func (h *SudaTokenHandler) lookupUserByAddress(ctx context.Context, addr common.Address) (*uuid.UUID, error) {
	if id, err := h.repo.FindUserByAddress(ctx, addr.Hex()); err == nil {
		return id, nil
	} else if !errors.Is(err, repository.ErrNotFound) {
		return nil, err
	}
	return nil, nil
}

func (h *SudaTokenHandler) isOurContract(addr common.Address) bool {
	return h.isContract(addr)
}

func (h *SudaTokenHandler) getBlockTime(ctx context.Context, blockNumber uint64) (time.Time, error) {
	return h.blockTimeFn(ctx, blockNumber)
}