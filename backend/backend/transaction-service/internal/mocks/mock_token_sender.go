package mocks

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/blockchain"
)

// MockTokenSender — мок tokenSender интерфейса. По умолчанию возвращает
// фейковый tx_hash, чтобы test'ы не падали на broadcast-шаге.
type MockTokenSender struct {
	SendFn func(ctx context.Context, signer blockchain.Signer, to common.Address, amount *big.Int) (string, error)

	SendCalls int

	// Capture последних аргументов вызова — для проверки в тестах.
	LastSigner blockchain.Signer
	LastTo     common.Address
	LastAmount *big.Int
}

func NewMockTokenSender() *MockTokenSender { return &MockTokenSender{} }

const defaultMockTxHash = "0x1111111111111111111111111111111111111111111111111111111111111111"

func (m *MockTokenSender) SendTokenTransfer(
	ctx context.Context, signer blockchain.Signer, to common.Address, amount *big.Int,
) (string, error) {
	m.SendCalls++
	m.LastSigner = signer
	m.LastTo = to
	m.LastAmount = amount
	if m.SendFn != nil {
		return m.SendFn(ctx, signer, to, amount)
	}
	return defaultMockTxHash, nil
}
