package repository

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// ────────────────────────────────────────────────────────────
//  Observer-only methods.
//  Этот файл расширяет Repository функциями, которые вызывает observer.
//  Записывать в tx_transactions может ТОЛЬКО observer (защищено append-only triggers).
// ────────────────────────────────────────────────────────────

// TxInsertParams — данные для одной строки tx_transactions.
type TxInsertParams struct {
	TxHash            string
	LogIndex          int32
	Type              string // P2P_TRANSFER, NFT_MINT, и т.п.
	Status            string // CONFIRMED | FAILED
	FromUserID        *uuid.UUID
	ToUserID          *uuid.UUID
	FromAddress       string
	ToAddress         string
	Amount            *big.Int // wei (для NFT — это 1 или 0)
	BlockNumber       int64
	ConfirmedAt       time.Time
	RelatedEntityType *string
	RelatedEntityID   *uuid.UUID
	Note              *string
}

// InsertTransaction — append-only вставка в tx_transactions.
//
// Возвращает inserted=true, если строка действительно добавлена, и
// inserted=false, если запись с этим (tx_hash, log_index) уже была
// (ON CONFLICT DO NOTHING → 0 затронутых строк).
//
// Observer использует этот флаг для идемпотентности: при повторной
// обработке блока (краш до UpsertObserverState, iter.Error в середине
// chunk'а) inserted=false сигнализирует «событие уже обработано» —
// нельзя повторно писать tx_donations и слать WS-уведомления.
func (r *Repository) InsertTransaction(ctx context.Context, p TxInsertParams) (bool, error) {
	const q = `
		INSERT INTO tx_transactions (
			tx_hash, log_index, type, status,
			from_user_id, to_user_id,
			from_address, to_address,
			amount, block_number, confirmed_at,
			related_entity_type, related_entity_id, note
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		ON CONFLICT (tx_hash, log_index) DO NOTHING
	`
	tag, err := r.pg.Exec(ctx, q,
		p.TxHash, p.LogIndex, p.Type, p.Status,
		p.FromUserID, p.ToUserID,
		p.FromAddress, p.ToAddress,
		p.Amount.String(), p.BlockNumber, p.ConfirmedAt,
		p.RelatedEntityType, p.RelatedEntityID, p.Note,
	)
	if err != nil {
		return false, fmt.Errorf("repo: insert transaction: %w", err)
	}
	return tag.RowsAffected() == 1, nil
}

// DeletePending удаляет запись из tx_pending. Вызывается observer'ом
// сразу после успешного InsertTransaction. UI больше не показывает
// «отправляется» индикатор.
func (r *Repository) DeletePending(ctx context.Context, txHash string) error {
	const q = `DELETE FROM tx_pending WHERE tx_hash = $1`
	_, err := r.pg.Exec(ctx, q, txHash)
	if err != nil {
		return fmt.Errorf("repo: delete pending: %w", err)
	}
	return nil
}

// FindUserByAddress — найти user_id по EVM-адресу.
// Возвращает (nil, ErrNotFound) если адреса нет в наших кошельках
// (значит это внешний адрес — контракт или old wallet).
func (r *Repository) FindUserByAddress(ctx context.Context, address string) (*uuid.UUID, error) {
	const q = `SELECT user_id FROM tx_wallets WHERE address = $1`
	var id uuid.UUID
	err := r.pg.QueryRow(ctx, q, address).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: find user by address: %w", err)
	}
	return &id, nil
}

// FindChannelByAddress — найти channel_id по EVM-адресу канала.
func (r *Repository) FindChannelByAddress(ctx context.Context, address string) (*uuid.UUID, error) {
	const q = `SELECT channel_id FROM tx_channel_wallets WHERE address = $1`
	var id uuid.UUID
	err := r.pg.QueryRow(ctx, q, address).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: find channel by address: %w", err)
	}
	return &id, nil
}

// ────────────────────────────────────────────────────────────
//  Observer state (tx_observer_state)
//  Хранит last_processed_block per contract, чтобы пережить рестарт.
// ────────────────────────────────────────────────────────────

// GetObserverState возвращает last_processed_block для контракта.
// Если записи нет — возвращает 0 (стартуем с генезиса). Caller должен
// учитывать что 0 может быть и валидным значением для свежего start'а.
func (r *Repository) GetObserverState(ctx context.Context, contractName string) (uint64, error) {
	const q = `SELECT last_processed_block FROM tx_observer_state WHERE contract_name = $1`
	var n uint64
	err := r.pg.QueryRow(ctx, q, contractName).Scan(&n)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrNotFound
		}
		return 0, fmt.Errorf("repo: get observer state: %w", err)
	}
	return n, nil
}

// UpsertObserverState — upsert last_processed_block для контракта.
func (r *Repository) UpsertObserverState(ctx context.Context, contractName string, blockNumber uint64) error {
	const q = `
		INSERT INTO tx_observer_state (contract_name, last_processed_block, updated_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (contract_name) DO UPDATE
		SET last_processed_block = EXCLUDED.last_processed_block,
		    updated_at = NOW()
	`
	_, err := r.pg.Exec(ctx, q, contractName, blockNumber)
	if err != nil {
		return fmt.Errorf("repo: upsert observer state: %w", err)
	}
	return nil
}

// PendingMeta — расширенная информация о pending транзакции, которую observer
// читает перед удалением tx_pending row.
type PendingMeta struct {
	ChatID          uuid.UUID // uuid.Nil если related_chat_id IS NULL
	ExpectedType    string    // "P2P_TRANSFER" | "DONATION" | "PURCHASE" | ...
	DonationMessage string    // пусто для не-DONATION expected_type
}

// GetPendingMeta возвращает chat_id + expected_type + donation_message
// для pending транзакции. Используется observer'ом ПЕРЕД удалением pending:
//
//   - если ChatID != uuid.Nil — это transfer-in-chat или donation в чате,
//     observer шлёт NotifyUserEvent с target_chat_id для создания system message;
//   - если ExpectedType == "DONATION" — observer дополнительно вставляет
//     строку в tx_donations и использует DonationMessage в system message.
//
// ErrNotFound если pending записи нет (внешний адрес или уже обработан).
func (r *Repository) GetPendingMeta(ctx context.Context, txHash string) (*PendingMeta, error) {
	const q = `
		SELECT COALESCE(related_chat_id, '00000000-0000-0000-0000-000000000000'::uuid),
		       expected_type,
		       COALESCE(donation_message, '')
		FROM tx_pending
		WHERE tx_hash = $1
	`
	var m PendingMeta
	err := r.pg.QueryRow(ctx, q, txHash).Scan(&m.ChatID, &m.ExpectedType, &m.DonationMessage)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get pending meta: %w", err)
	}
	return &m, nil
}
