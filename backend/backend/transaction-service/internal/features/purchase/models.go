// Package purchase — имитация покупки SUDA за фиатные деньги.
//
// Это pet-project'овая фича: реальной финансовой транзакции нет, юзер
// вводит фейковые данные карты, сервис делает sleep 2-3 сек ("processing"),
// потом treasury переводит SUDA пользователю обычным on-chain transfer'ом.
//
// Жизненный цикл записи tx_suda_purchases:
//
//	PENDING    — POST /purchase/initiate создал row, ждём confirm
//	PROCESSING — пришёл POST /purchase/:id/confirm, имитируем оплату
//	COMPLETED  — treasury успешно отправил Transfer, tx_hash записан
//	FAILED     — что-то упало (no_wallet / broadcast error); failure_reason заполнен
package purchase

import (
	"math/big"
	"time"

	"github.com/google/uuid"
)

// Статусы tx_suda_purchases.status.
const (
	StatusPending    = "PENDING"
	StatusProcessing = "PROCESSING"
	StatusCompleted  = "COMPLETED"
	StatusFailed     = "FAILED"
)

// Способы оплаты — чисто декоративные, для UI.
const (
	PaymentMethodCard      = "CARD"
	PaymentMethodApplePay  = "APPLE_PAY"
	PaymentMethodGooglePay = "GOOGLE_PAY"
)

// Тип tx_transactions.type для индексации observer'ом.
const TxTypePurchase = "PURCHASE"

// Operation для tx_signing_audit.operation (от лица treasury).
const OpPurchase = "PURCHASE"

// Purchase — одна запись об "имитированной" покупке.
type Purchase struct {
	ID            uuid.UUID
	UserID        uuid.UUID
	PackageCode   string
	SudaAmountWei *big.Int   // wei
	FiatAmount    string     // "9.99" — храним как string чтобы не терять копейки
	FiatCurrency  string     // "USD"
	Status        string     // PENDING | PROCESSING | COMPLETED | FAILED
	PaymentMethod string     // CARD | APPLE_PAY | GOOGLE_PAY
	TxHash        string     // пусто пока broadcast не сделан
	FailureReason string     // заполнен только при Status=FAILED
	CreatedAt     time.Time
	CompletedAt   *time.Time
}
