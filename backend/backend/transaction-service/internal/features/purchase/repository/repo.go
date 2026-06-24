// Package repository — доступ к tx_suda_purchases.
//
// Реализован отдельно от wallet/repository, потому что таблица не пересекается
// и фича может в перспективе ехать отдельным процессом (биллинг).
package repository

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase"
)

// ErrNotFound — purchase row с таким id не существует.
var ErrNotFound = errors.New("purchase not found")

type Repository struct {
	pg *pgxpool.Pool
}

func New(pg *pgxpool.Pool) *Repository {
	return &Repository{pg: pg}
}

// CreateParams — параметры для INSERT новой PENDING покупки.
type CreateParams struct {
	UserID        uuid.UUID
	PackageCode   string
	SudaAmountWei *big.Int
	FiatAmount    string
	FiatCurrency  string
	PaymentMethod string
}

// Create вставляет purchase row в статусе PENDING и возвращает её id.
func (r *Repository) Create(ctx context.Context, p CreateParams) (uuid.UUID, error) {
	const q = `
		INSERT INTO tx_suda_purchases (
			user_id, package_code, suda_amount_wei,
			fiat_amount, fiat_currency, status, payment_method
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	var id uuid.UUID
	err := r.pg.QueryRow(ctx, q,
		p.UserID, p.PackageCode, p.SudaAmountWei.String(),
		p.FiatAmount, p.FiatCurrency, purchase.StatusPending, p.PaymentMethod,
	).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("repo: create purchase: %w", err)
	}
	return id, nil
}

// GetByID — одна запись по id.
func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*purchase.Purchase, error) {
	const q = `
		SELECT id, user_id, package_code, suda_amount_wei::text,
		       fiat_amount::text, fiat_currency, status,
		       COALESCE(payment_method, ''),
		       COALESCE(tx_hash, ''),
		       COALESCE(failure_reason, ''),
		       created_at, completed_at
		FROM tx_suda_purchases
		WHERE id = $1
	`
	var p purchase.Purchase
	var sudaWei string
	err := r.pg.QueryRow(ctx, q, id).Scan(
		&p.ID, &p.UserID, &p.PackageCode, &sudaWei,
		&p.FiatAmount, &p.FiatCurrency, &p.Status,
		&p.PaymentMethod, &p.TxHash, &p.FailureReason,
		&p.CreatedAt, &p.CompletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: get purchase: %w", err)
	}
	p.SudaAmountWei, _ = new(big.Int).SetString(sudaWei, 10)
	return &p, nil
}

// MarkProcessing — переход PENDING -> PROCESSING. Используется CAS-стилем:
// если status != PENDING, обновление не пройдёт (rows affected = 0) и
// caller получит ErrNotFound — это защищает от двойного confirm параллельно.
func (r *Repository) MarkProcessing(ctx context.Context, id uuid.UUID) error {
	const q = `
		UPDATE tx_suda_purchases
		SET status = $1
		WHERE id = $2 AND status = $3
	`
	tag, err := r.pg.Exec(ctx, q, purchase.StatusProcessing, id, purchase.StatusPending)
	if err != nil {
		return fmt.Errorf("repo: mark processing: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

// MarkCompleted — успешное завершение: пишем tx_hash и completed_at.
func (r *Repository) MarkCompleted(ctx context.Context, id uuid.UUID, txHash string) error {
	const q = `
		UPDATE tx_suda_purchases
		SET status = $1, tx_hash = $2, completed_at = NOW()
		WHERE id = $3
	`
	_, err := r.pg.Exec(ctx, q, purchase.StatusCompleted, txHash, id)
	if err != nil {
		return fmt.Errorf("repo: mark completed: %w", err)
	}
	return nil
}

// MarkFailed — провал: пишем причину, статус FAILED.
func (r *Repository) MarkFailed(ctx context.Context, id uuid.UUID, reason string) error {
	const q = `
		UPDATE tx_suda_purchases
		SET status = $1, failure_reason = $2, completed_at = NOW()
		WHERE id = $3
	`
	_, err := r.pg.Exec(ctx, q, purchase.StatusFailed, reason, id)
	if err != nil {
		return fmt.Errorf("repo: mark failed: %w", err)
	}
	return nil
}

// ListForUser возвращает пагинированную историю покупок юзера.
func (r *Repository) ListForUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]purchase.Purchase, error) {
	const q = `
		SELECT id, user_id, package_code, suda_amount_wei::text,
		       fiat_amount::text, fiat_currency, status,
		       COALESCE(payment_method, ''),
		       COALESCE(tx_hash, ''),
		       COALESCE(failure_reason, ''),
		       created_at, completed_at
		FROM tx_suda_purchases
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.pg.Query(ctx, q, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("repo: list purchases: %w", err)
	}
	defer rows.Close()

	var out []purchase.Purchase
	for rows.Next() {
		var p purchase.Purchase
		var sudaWei string
		if err := rows.Scan(
			&p.ID, &p.UserID, &p.PackageCode, &sudaWei,
			&p.FiatAmount, &p.FiatCurrency, &p.Status,
			&p.PaymentMethod, &p.TxHash, &p.FailureReason,
			&p.CreatedAt, &p.CompletedAt,
		); err != nil {
			return nil, fmt.Errorf("repo: scan purchase: %w", err)
		}
		p.SudaAmountWei, _ = new(big.Int).SetString(sudaWei, 10)
		out = append(out, p)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("repo: list purchases rows: %w", err)
	}
	return out, nil
}

// FindByTxHash — observer'у нужно понять "пришедший Transfer — это purchase?".
// Возвращает purchase id и suda_amount, ErrNotFound если нет такой записи.
type TxHashMatch struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	SudaAmountWei *big.Int
	PackageCode   string
}

// FindByTxHash — observer'у нужно понять, не пришёл ли Transfer от treasury
// именно как purchase confirmation. Возвращает ErrNotFound, если такой
// transaction-hash не привязан ни к одной записи.
func (r *Repository) FindByTxHash(ctx context.Context, txHash string) (*TxHashMatch, error) {
	const q = `
		SELECT id, user_id, suda_amount_wei::text, package_code
		FROM tx_suda_purchases
		WHERE tx_hash = $1
		LIMIT 1
	`
	var m TxHashMatch
	var sudaWei string
	err := r.pg.QueryRow(ctx, q, txHash).Scan(&m.ID, &m.UserID, &sudaWei, &m.PackageCode)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("repo: find by tx_hash: %w", err)
	}
	m.SudaAmountWei, _ = new(big.Int).SetString(sudaWei, 10)
	return &m, nil
}
