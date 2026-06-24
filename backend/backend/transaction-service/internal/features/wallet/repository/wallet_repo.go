// Package repository — доступ к таблицам wallet-фичи (tx_wallets,
// tx_channel_wallets, tx_pending, tx_signing_audit, tx_transactions для чтения).
//
// Запись в tx_transactions делает ТОЛЬКО observer. Repository здесь
// только читает её для GetHistory.
package repository

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet"
)

// ErrNotFound — кошелёк (или сущность) не найден.
var ErrNotFound = errors.New("not found")

// ErrAlreadyExists — попытка создать кошелёк, который уже есть.
// Возникает из-за UNIQUE constraint на user_id/channel_id.
var ErrAlreadyExists = errors.New("already exists")

type Repository struct {
	pg *pgxpool.Pool
}

func New(pg *pgxpool.Pool) *Repository {
	return &Repository{pg: pg}
}

// ────────────────────────────────────────────────────────────
//  USER WALLETS
// ────────────────────────────────────────────────────────────

// CreateUserWallet вставляет запись о новом кошельке юзера.
// Если UNIQUE(user_id) сработал — возвращает ErrAlreadyExists.
//
// Caller дальше должен сделать GetUserWallet и вернуть existed=true.
// Этот pattern (insert-or-select) реализует идемпотентность
// CreateWalletForUser, как мы договаривались.
func (r *Repository) CreateUserWallet(
	ctx context.Context,
	userID uuid.UUID,
	address string,
	encryptedPK string,
	keyVersion uint8,
) error {
	const q = `
		INSERT INTO tx_wallets (user_id, address, encrypted_private_key, key_version)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.pg.Exec(ctx, q, userID, address, encryptedPK, keyVersion)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("repo: create user wallet: %w", err)
	}
	return nil
}

func (r *Repository) GetUserWallet(ctx context.Context, userID uuid.UUID) (*wallet.Wallet, error) {
	const q = `
		SELECT user_id, address, encrypted_private_key, key_version,
		       COALESCE(suda_balance_cache, 0), balance_updated_at, created_at
		FROM tx_wallets
		WHERE user_id = $1
	`
	var w wallet.Wallet
	var balCache string
	err := r.pg.QueryRow(ctx, q, userID).Scan(
		&w.UserID, &w.Address, &w.EncryptedPrivateKey, &w.KeyVersion,
		&balCache, &w.BalanceUpdatedAt, &w.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get user wallet: %w", err)
	}
	w.SudaBalanceCache, _ = new(big.Int).SetString(balCache, 10)
	return &w, nil
}

func (r *Repository) GetUserWalletByAddress(ctx context.Context, address string) (*wallet.Wallet, error) {
	const q = `
		SELECT user_id, address, encrypted_private_key, key_version,
		       COALESCE(suda_balance_cache, 0), balance_updated_at, created_at
		FROM tx_wallets
		WHERE address = $1
	`
	var w wallet.Wallet
	var balCache string
	err := r.pg.QueryRow(ctx, q, address).Scan(
		&w.UserID, &w.Address, &w.EncryptedPrivateKey, &w.KeyVersion,
		&balCache, &w.BalanceUpdatedAt, &w.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get user wallet by address: %w", err)
	}
	w.SudaBalanceCache, _ = new(big.Int).SetString(balCache, 10)
	return &w, nil
}

// UpdateUserBalanceCache обновляет suda_balance_cache (используется как кеш).
// НЕ источник правды — реальный баланс читается из SudaToken.balanceOf().
func (r *Repository) UpdateUserBalanceCache(
	ctx context.Context,
	userID uuid.UUID,
	balanceWei *big.Int,
) error {
	const q = `
		UPDATE tx_wallets
		SET suda_balance_cache = $1, balance_updated_at = NOW()
		WHERE user_id = $2
	`
	_, err := r.pg.Exec(ctx, q, balanceWei.String(), userID)
	if err != nil {
		return fmt.Errorf("repo: update balance cache: %w", err)
	}
	return nil
}

// ────────────────────────────────────────────────────────────
//  CHANNEL WALLETS
// ────────────────────────────────────────────────────────────

func (r *Repository) CreateChannelWallet(
	ctx context.Context,
	channelID uuid.UUID,
	address string,
	encryptedPK string,
	keyVersion uint8,
) error {
	const q = `
		INSERT INTO tx_channel_wallets (channel_id, address, encrypted_private_key, key_version)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.pg.Exec(ctx, q, channelID, address, encryptedPK, keyVersion)
	if err != nil {
		if isUniqueViolation(err) {
			return ErrAlreadyExists
		}
		return fmt.Errorf("repo: create channel wallet: %w", err)
	}
	return nil
}

func (r *Repository) GetChannelWallet(ctx context.Context, channelID uuid.UUID) (*wallet.ChannelWallet, error) {
	const q = `
		SELECT channel_id, address, encrypted_private_key, key_version,
		       COALESCE(suda_balance_cache, 0), balance_updated_at, created_at
		FROM tx_channel_wallets
		WHERE channel_id = $1
	`
	var w wallet.ChannelWallet
	var balCache string
	err := r.pg.QueryRow(ctx, q, channelID).Scan(
		&w.ChannelID, &w.Address, &w.EncryptedPrivateKey, &w.KeyVersion,
		&balCache, &w.BalanceUpdatedAt, &w.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get channel wallet: %w", err)
	}
	w.SudaBalanceCache, _ = new(big.Int).SetString(balCache, 10)
	return &w, nil
}

// ────────────────────────────────────────────────────────────
//  PENDING TRANSACTIONS — UI индикатор «отправляется»
// ────────────────────────────────────────────────────────────

func (r *Repository) InsertPending(
	ctx context.Context,
	txHash string,
	fromUserID uuid.UUID,
	expectedType string,
	chatID uuid.UUID,
) error {
	const q = `
		INSERT INTO tx_pending (tx_hash, from_user_id, expected_type, related_chat_id)
		VALUES ($1, $2, $3, NULLIF($4, '00000000-0000-0000-0000-000000000000'::uuid))
		ON CONFLICT (tx_hash) DO NOTHING
	`
	_, err := r.pg.Exec(ctx, q, txHash, fromUserID, expectedType, chatID)
	if err != nil {
		return fmt.Errorf("repo: insert pending: %w", err)
	}
	return nil
}

// InsertPendingTransfer — как InsertPending для P2P_TRANSFER, но дополнительно
// сохраняет note (комментарий к переводу) в donation_message — observer прочитает
// его и положит в payload system message в чате. note может быть пустым.
func (r *Repository) InsertPendingTransfer(
	ctx context.Context,
	txHash string,
	fromUserID uuid.UUID,
	chatID uuid.UUID,
	note string,
) error {
	const q = `
		INSERT INTO tx_pending (tx_hash, from_user_id, expected_type, related_chat_id, donation_message)
		VALUES ($1, $2, 'P2P_TRANSFER',
		        NULLIF($3, '00000000-0000-0000-0000-000000000000'::uuid),
		        NULLIF($4, ''))
		ON CONFLICT (tx_hash) DO NOTHING
	`
	_, err := r.pg.Exec(ctx, q, txHash, fromUserID, chatID, note)
	if err != nil {
		return fmt.Errorf("repo: insert pending transfer: %w", err)
	}
	return nil
}

// InsertPendingDonation — то же что InsertPending, но дополнительно сохраняет
// donation_message (текст доната), который observer прочитает и положит в
// tx_donations.message + system message в чате при индексации Transfer event'а.
//
// expectedType всегда "DONATION" (передаётся для согласованности с InsertPending).
// chatID = uuid.Nil допустим (донат вне чата), тогда system message не создаётся.
func (r *Repository) InsertPendingDonation(
	ctx context.Context,
	txHash string,
	fromUserID uuid.UUID,
	chatID uuid.UUID,
	donationMessage string,
) error {
	const q = `
		INSERT INTO tx_pending (tx_hash, from_user_id, expected_type, related_chat_id, donation_message)
		VALUES ($1, $2, 'DONATION',
		        NULLIF($3, '00000000-0000-0000-0000-000000000000'::uuid),
		        NULLIF($4, ''))
		ON CONFLICT (tx_hash) DO NOTHING
	`
	_, err := r.pg.Exec(ctx, q, txHash, fromUserID, chatID, donationMessage)
	if err != nil {
		return fmt.Errorf("repo: insert pending donation: %w", err)
	}
	return nil
}
// ────────────────────────────────────────────────────────────
//  SIGNING AUDIT
// ────────────────────────────────────────────────────────────

// AuditEntry — параметры для одной строки tx_signing_audit.
// txHash может быть пустым если запись делается ДО broadcast'а.
type AuditEntry struct {
	SubjectType string    // wallet.SubjectUser | wallet.SubjectChannel
	SubjectID   uuid.UUID
	Operation   string    // wallet.OpTransfer и т.п.
	TxHash      string    // опционально
	RequestIP   string
	UserAgent   string
}

func (r *Repository) WriteAudit(ctx context.Context, e AuditEntry) error {
	const q = `
		INSERT INTO tx_signing_audit (subject_type, subject_id, operation, tx_hash, request_ip, user_agent)
		VALUES ($1, $2, $3, NULLIF($4, ''), $5, $6)
	`
	_, err := r.pg.Exec(ctx, q, e.SubjectType, e.SubjectID, e.Operation, e.TxHash, e.RequestIP, e.UserAgent)
	if err != nil {
		return fmt.Errorf("repo: write audit: %w", err)
	}
	return nil
}

// ────────────────────────────────────────────────────────────
//  TRANSACTION HISTORY — чтение (запись делает observer)
// ────────────────────────────────────────────────────────────

// TransactionHistoryItem — одна запись истории для UI.
type TransactionHistoryItem struct {
	ID                uuid.UUID
	TxHash            string
	LogIndex          int32
	FromUserID        *uuid.UUID
	ToUserID          *uuid.UUID
	FromAddress       string
	ToAddress         string
	Amount            *big.Int // wei
	Type              string
	Status            string
	RelatedEntityType *string
	RelatedEntityID   *uuid.UUID
	Note              *string
	BlockNumber       int64
	ConfirmedAt       time.Time
	CreatedAt         time.Time
}

// GetHistoryForUser — последние транзакции, в которых userID был отправителем
// или получателем. Сортировка по дате подтверждения (новые сверху).
func (r *Repository) GetHistoryForUser(
	ctx context.Context,
	userID uuid.UUID,
	limit, offset int,
) ([]TransactionHistoryItem, error) {
	const q = `
		SELECT id, tx_hash, log_index, from_user_id, to_user_id,
		       from_address, to_address, amount::text, type, status,
		       related_entity_type, related_entity_id, note,
		       block_number, confirmed_at, created_at
		FROM tx_transactions
		WHERE from_user_id = $1 OR to_user_id = $1
		ORDER BY confirmed_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pg.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repo: get history: %w", err)
	}
	defer rows.Close()

	var out []TransactionHistoryItem
	for rows.Next() {
		var it TransactionHistoryItem
		var amount string
		if err := rows.Scan(
			&it.ID, &it.TxHash, &it.LogIndex, &it.FromUserID, &it.ToUserID,
			&it.FromAddress, &it.ToAddress, &amount, &it.Type, &it.Status,
			&it.RelatedEntityType, &it.RelatedEntityID, &it.Note,
			&it.BlockNumber, &it.ConfirmedAt, &it.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("repo: scan history row: %w", err)
		}
		it.Amount, _ = new(big.Int).SetString(amount, 10)
		out = append(out, it)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo: history rows: %w", err)
	}
	return out, nil
}

// ────────────────────────────────────────────────────────────
//  GATING — чтение правила, читать NFT-ownership
// ────────────────────────────────────────────────────────────

// GatingRule — одно правило token-gating для чата/канала.
type GatingRule struct {
	ChatID                  uuid.UUID
	MinSudaBalance          *big.Int   // wei (legacy hold-порог; для платной подписки не используется)
	RequiredNFTCollectionID *uuid.UUID // опционально (legacy)
	SubscriptionPrice       *big.Int   // wei — цена платной подписки на PUBLIC-канал (0 = не платный)
}

// GetGatingRule возвращает правило по chat_id, или ErrNotFound если правила
// для этого чата нет (= чат open для всех).
func (r *Repository) GetGatingRule(ctx context.Context, chatID uuid.UUID) (*GatingRule, error) {
	const q = `
		SELECT chat_id, min_suda_balance::text, required_nft_collection_id, subscription_price_wei::text
		FROM tx_gating_rules
		WHERE chat_id = $1
	`
	var rule GatingRule
	var minBal, price string
	err := r.pg.QueryRow(ctx, q, chatID).Scan(&rule.ChatID, &minBal, &rule.RequiredNFTCollectionID, &price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get gating rule: %w", err)
	}
	rule.MinSudaBalance, _ = new(big.Int).SetString(minBal, 10)
	rule.SubscriptionPrice, _ = new(big.Int).SetString(price, 10)
	return &rule, nil
}

// UpsertGatingRule создаёт или обновляет правило token-gating для чата/канала.
// При повторном вызове (ON CONFLICT chat_id) перезаписывает min_suda_balance
// и required_nft_collection_id, created_by НЕ трогаем (остаётся первый автор).
func (r *Repository) UpsertGatingRule(
	ctx context.Context,
	chatID uuid.UUID,
	minSudaBalance *big.Int,
	nftCollectionID *uuid.UUID,
	subscriptionPrice *big.Int,
	createdBy uuid.UUID,
) error {
	if subscriptionPrice == nil {
		subscriptionPrice = big.NewInt(0)
	}
	const q = `
		INSERT INTO tx_gating_rules (chat_id, min_suda_balance, required_nft_collection_id, subscription_price_wei, created_by)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (chat_id) DO UPDATE
		SET min_suda_balance = EXCLUDED.min_suda_balance,
		    required_nft_collection_id = EXCLUDED.required_nft_collection_id,
		    subscription_price_wei = EXCLUDED.subscription_price_wei,
		    updated_at = NOW()
	`
	_, err := r.pg.Exec(ctx, q, chatID, minSudaBalance.String(), nftCollectionID, subscriptionPrice.String(), createdBy)
	if err != nil {
		return fmt.Errorf("repo: upsert gating rule: %w", err)
	}
	return nil
}

// DeleteGatingRule снимает правило gating для чата/канала.
// Идемпотентно: если правила не было — это не ошибка.
func (r *Repository) DeleteGatingRule(ctx context.Context, chatID uuid.UUID) error {
	const q = `DELETE FROM tx_gating_rules WHERE chat_id = $1`
	_, err := r.pg.Exec(ctx, q, chatID)
	if err != nil {
		return fmt.Errorf("repo: delete gating rule: %w", err)
	}
	return nil
}

// ────────────────────────────────────────────────────────────
//  DONATIONS — запись делает observer после индексации Transfer event'а.
// ────────────────────────────────────────────────────────────

// DonationInsertParams — параметры строки tx_donations.
//
// ChannelID и ToUserID — взаимоисключающие: для donation каналу заполняем
// ChannelID, для P2P доната заполняем ToUserID, второе остаётся nil.
type DonationInsertParams struct {
	ChannelID   *uuid.UUID
	ToUserID    *uuid.UUID
	FromUserID  uuid.UUID
	FromAddress string
	ToAddress   string
	Amount      *big.Int
	Message     string
	TxHash      string
}

// InsertDonation — append-only вставка в tx_donations.
// Дубли пишем не пытаемся отлавливать: UNIQUE-индекса на (tx_hash) нет, но
// observer вызывается ровно один раз за event благодаря (tx_hash, log_index)
// constraint'у в tx_transactions.
func (r *Repository) InsertDonation(ctx context.Context, p DonationInsertParams) error {
	const q = `
		INSERT INTO tx_donations (
			channel_id, to_user_id, from_user_id,
			from_address, to_address, amount, message, tx_hash
		)
		VALUES ($1, $2, $3, $4, $5, $6, NULLIF($7, ''), $8)
	`
	_, err := r.pg.Exec(ctx, q,
		p.ChannelID, p.ToUserID, p.FromUserID,
		p.FromAddress, p.ToAddress, p.Amount.String(), p.Message, p.TxHash,
	)
	if err != nil {
		return fmt.Errorf("repo: insert donation: %w", err)
	}
	return nil
}

// ────────────────────────────────────────────────────────────
//  TREASURY STATS — агрегаты и список донатов канала (чтение tx_donations).
// ────────────────────────────────────────────────────────────

// GetChannelDonationStats возвращает количество донатов и суммарный объём
// (wei) для канала. Если донатов нет — (0, big.NewInt(0), nil).
func (r *Repository) GetChannelDonationStats(
	ctx context.Context, channelID uuid.UUID,
) (count int64, totalWei *big.Int, err error) {
	const q = `
		SELECT COUNT(*), COALESCE(SUM(amount), 0)::text
		FROM tx_donations
		WHERE channel_id = $1
	`
	var total string
	if err := r.pg.QueryRow(ctx, q, channelID).Scan(&count, &total); err != nil {
		return 0, nil, fmt.Errorf("repo: channel donation stats: %w", err)
	}
	totalWei, _ = new(big.Int).SetString(total, 10)
	if totalWei == nil {
		totalWei = big.NewInt(0)
	}
	return count, totalWei, nil
}

// TopDonor — агрегат одного донатера канала.
type TopDonor struct {
	FromUserID    uuid.UUID
	TotalWei      *big.Int
	DonationCount int64
}

// GetChannelTopDonors возвращает топ-N донатеров канала по суммарному объёму.
func (r *Repository) GetChannelTopDonors(
	ctx context.Context, channelID uuid.UUID, limit int,
) ([]TopDonor, error) {
	const q = `
		SELECT from_user_id, SUM(amount)::text, COUNT(*)
		FROM tx_donations
		WHERE channel_id = $1
		GROUP BY from_user_id
		ORDER BY SUM(amount) DESC
		LIMIT $2
	`
	rows, err := r.pg.Query(ctx, q, channelID, limit)
	if err != nil {
		return nil, fmt.Errorf("repo: channel top donors: %w", err)
	}
	defer rows.Close()

	var out []TopDonor
	for rows.Next() {
		var d TopDonor
		var total string
		if err := rows.Scan(&d.FromUserID, &total, &d.DonationCount); err != nil {
			return nil, fmt.Errorf("repo: scan top donor: %w", err)
		}
		d.TotalWei, _ = new(big.Int).SetString(total, 10)
		out = append(out, d)
	}
	return out, rows.Err()
}

// DonationRecord — одна строка списка донатов канала.
type DonationRecord struct {
	ID          uuid.UUID
	FromUserID  uuid.UUID
	FromAddress string
	Amount      *big.Int
	Message     string
	TxHash      string
	CreatedAt   time.Time
}

// GetChannelDonations возвращает пагинированный список донатов канала,
// новые сверху.
func (r *Repository) GetChannelDonations(
	ctx context.Context, channelID uuid.UUID, limit, offset int,
) ([]DonationRecord, error) {
	const q = `
		SELECT id, from_user_id, from_address, amount::text,
		       COALESCE(message, ''), tx_hash, created_at
		FROM tx_donations
		WHERE channel_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pg.Query(ctx, q, channelID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repo: channel donations: %w", err)
	}
	defer rows.Close()

	var out []DonationRecord
	for rows.Next() {
		var d DonationRecord
		var amount string
		if err := rows.Scan(
			&d.ID, &d.FromUserID, &d.FromAddress, &amount,
			&d.Message, &d.TxHash, &d.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("repo: scan donation: %w", err)
		}
		d.Amount, _ = new(big.Int).SetString(amount, 10)
		out = append(out, d)
	}
	return out, rows.Err()
}

// UserOwnsNFTFromCollection — проверяет что у юзера есть хотя бы один NFT
// из указанной коллекции. Дешёвый SELECT по индексу idx_tx_nft_owner.
//
// Доверяемся observer'у — если он отстаёт, временно gating может пропустить
// юзера, у которого NFT только что появился, или (наоборот) не пропустить
// юзера, который только что NFT передал. Лаг < 2 секунды.
func (r *Repository) UserOwnsNFTFromCollection(
	ctx context.Context,
	userID uuid.UUID,
	collectionID uuid.UUID,
) (bool, error) {
	const q = `
		SELECT EXISTS (
			SELECT 1 FROM tx_nft_items
			WHERE owner_user_id = $1 AND collection_id = $2
			LIMIT 1
		)
	`
	var owns bool
	if err := r.pg.QueryRow(ctx, q, userID, collectionID).Scan(&owns); err != nil {
		return false, fmt.Errorf("repo: nft ownership: %w", err)
	}
	return owns, nil
}

// ────────────────────────────────────────────────────────────
//  Helpers
// ────────────────────────────────────────────────────────────

// isUniqueViolation определяет ошибку pgx как нарушение UNIQUE constraint.
// pgx v5 предоставляет PgError со SQLState() — используем его через интерфейс,
// чтобы не тянуть прямую зависимость на pgconn.
func isUniqueViolation(err error) bool {
	if e, ok := err.(interface{ SQLState() string }); ok {
		return e.SQLState() == "23505"
	}
	return false
}