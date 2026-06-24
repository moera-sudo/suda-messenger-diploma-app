package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase/repository"
	walletrepo "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
)

// ConfirmInput — параметры POST /purchase/:id/confirm.
//
// CardNumber/CVV намеренно НЕ валидируются и НЕ сохраняются — это pet-project
// имитация, реальной оплаты нет. Поля принимаются только чтобы UI мог
// "отрендерить форму как настоящий чекаут".
type ConfirmInput struct {
	PurchaseID uuid.UUID
	UserID     uuid.UUID

	RequestIP string
	UserAgent string
}

// ConfirmResult — что отдаём клиенту после успешного confirm.
type ConfirmResult struct {
	PurchaseID    uuid.UUID
	Status        string // COMPLETED | FAILED
	TxHash        string // пусто при FAILED
	SudaAmountWei string
}

// rateLimitWindow — минимальный интервал между двумя confirm'ами одного юзера.
const rateLimitWindow = 5 * time.Second

// Confirm — имитация процессинга:
//
//  1. Redis SET NX rate-limit (5 сек / user).
//  2. Загрузить purchase row, проверить ownership и status=PENDING.
//  3. MarkProcessing (CAS) — защита от двойного confirm.
//  4. Sleep 2-3 сек (имитация работы платёжного провайдера).
//  5. Найти кошелёк юзера → если нет, MarkFailed("no_wallet").
//  6. Audit ДО broadcast (subject = USER (получатель бонуса), op=PURCHASE).
//  7. Treasury Transfer на адрес юзера через abigen-биндинги.
//  8. MarkCompleted с tx_hash.
//  9. Записать pending для UI-индикатора (observer его удалит).
//
// Подтверждение PURCHASE_COMPLETED придёт через observer (см. suda_token.go),
// сразу после индексации блока с этим Transfer.
func (s *Service) Confirm(ctx context.Context, in ConfirmInput) (*ConfirmResult, error) {
	// 1. Rate limit на confirm.
	if err := s.acquireConfirmLock(ctx, in.UserID); err != nil {
		return nil, err
	}

	// 2. Загрузка и проверки.
	p, err := s.repo.GetByID(ctx, in.PurchaseID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("confirm: get: %w", err)
	}
	if p.UserID != in.UserID {
		return nil, ErrForbidden
	}
	if p.Status != purchase.StatusPending {
		return nil, ErrAlreadyHandled
	}

	// 3. CAS-переход PENDING -> PROCESSING.
	if err := s.repo.MarkProcessing(ctx, p.ID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, ErrAlreadyHandled
		}
		return nil, fmt.Errorf("confirm: mark processing: %w", err)
	}

	// 4. Имитация платёжного процессинга (2-3 секунды; в тестах no-op).
	if err := s.simulateProcessing(ctx); err != nil {
		// Контекст отменён — позволим админу/CLI завершить.
		_ = s.repo.MarkFailed(ctx, p.ID, "context_cancelled")
		return nil, err
	}

	// 5. Кошелёк юзера.
	w, err := s.walletRepo.GetUserWallet(ctx, p.UserID)
	if err != nil {
		if errors.Is(err, walletrepo.ErrNotFound) {
			_ = s.repo.MarkFailed(ctx, p.ID, "no_wallet")
			return nil, ErrWalletNotFound
		}
		_ = s.repo.MarkFailed(ctx, p.ID, "wallet_lookup_failed")
		return nil, fmt.Errorf("confirm: wallet lookup: %w", err)
	}

	// 6. Audit запись (до broadcast'а; subject = USER-получатель бонуса).
	_ = s.walletRepo.WriteAudit(ctx, walletrepo.AuditEntry{
		SubjectType: auditSubjectUser,
		SubjectID:   p.UserID,
		Operation:   purchase.OpPurchase,
		TxHash:      "",
		RequestIP:   in.RequestIP,
		UserAgent:   in.UserAgent,
	})

	// 7. Treasury -> user Transfer (один высокоуровневый вызов; маскируется в тестах).
	toAddr := common.HexToAddress(w.Address)
	txHash, err := s.sender.SendTokenTransfer(ctx, s.treasury, toAddr, p.SudaAmountWei)
	if err != nil {
		_ = s.repo.MarkFailed(ctx, p.ID, fmt.Sprintf("broadcast_failed: %v", err))
		return nil, fmt.Errorf("confirm: broadcast: %w", err)
	}

	// 8. Финальный статус.
	if err := s.repo.MarkCompleted(ctx, p.ID, txHash); err != nil {
		// TX уже улетела, статус неконсистентен. Лог + игнор: observer всё равно
		// проиндексирует Transfer, юзер получит SUDA.
		log.Error().Err(err).Str("purchase_id", p.ID.String()).Str("tx_hash", txHash).
			Msg("purchase: failed to mark completed; tx already broadcast")
	}

	// 9. Pending row — UI увидит "transaction confirming" пока observer не закроет.
	_ = s.walletRepo.InsertPending(ctx, txHash, p.UserID, purchase.TxTypePurchase, uuid.Nil)

	log.Info().
		Str("user_id", p.UserID.String()).
		Str("purchase_id", p.ID.String()).
		Str("package", p.PackageCode).
		Str("suda_wei", p.SudaAmountWei.String()).
		Str("tx_hash", txHash).
		Msg("purchase confirmed and broadcast")

	return &ConfirmResult{
		PurchaseID:    p.ID,
		Status:        purchase.StatusCompleted,
		TxHash:        txHash,
		SudaAmountWei: p.SudaAmountWei.String(),
	}, nil
}

// acquireConfirmLock пытается захватить Redis-лок на confirm для конкретного
// user_id. Возвращает ErrRateLimited если лок занят.
//
// Если Redis недоступен — пропускаем (fail-open). Pet-project, риск двойного
// confirm'а в редком кейсе приемлемый, а блокировать всех при сетевой
// проблеме — хуже.
func (s *Service) acquireConfirmLock(ctx context.Context, userID uuid.UUID) error {
	key := fmt.Sprintf("purchase:lock:%s", userID.String())
	ok, err := s.redis.SetNX(ctx, key, "1", rateLimitWindow).Result()
	if err != nil {
		log.Warn().Err(err).Str("user_id", userID.String()).Msg("redis lock failed (fail-open)")
		return nil
	}
	if !ok {
		return ErrRateLimited
	}
	return nil
}

// defaultSimulateProcessing блокирует goroutine на 2-3 секунды или до отмены ctx.
// Не зависит от внешних ресурсов; нужен только чтобы юзер увидел спиннер
// и поверил, что "оплата проходит". В тестах заменяется на no-op.
func defaultSimulateProcessing(ctx context.Context) error {
	jitter := time.Duration(2000+rand.Intn(1000)) * time.Millisecond
	select {
	case <-time.After(jitter):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
