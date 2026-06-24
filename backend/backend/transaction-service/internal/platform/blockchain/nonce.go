package blockchain

import (
	"context"
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

// NonceManager сериализует получение nonce и broadcast для каждого
// адреса отдельно. Между Lock() и Release() другие операции с этим же
// адресом ждут в очереди.
//
// Зачем это нужно: если юзер за секунду делает два перевода, оба
// handler'а параллельно запросят eth_getTransactionCount(latest), оба
// получат одинаковый nonce, обе подписанные транзакции пойдут в Besu,
// вторая упадёт ("nonce too low"). NonceManager превращает гонку
// в очередь: первая операция возьмёт nonce N, broadcast'нет, отпустит
// lock, вторая — N+1.
//
// Ограничения v1:
//   - In-memory: один инстанс transaction-service. При горизонтальном
//     масштабировании заменить на Redis INCR.
//   - sync.Mutex per address: первая операция блокирует всю очередь
//     для этого адреса. Если broadcast завис на 10 сек — следующие
//     запросы юзера тоже зависнут. Это правильное поведение, но юзер
//     должен понимать.
type NonceManager struct {
	mu    sync.Mutex
	locks map[common.Address]*sync.Mutex
}

func NewNonceManager() *NonceManager {
	return &NonceManager{
		locks: make(map[common.Address]*sync.Mutex),
	}
}

// LockAddress блокирует все подписи для адреса до Release.
// Используется так:
//
//	release := nm.LockAddress(addr)
//	defer release()
//	nonce, _ := client.Eth().PendingNonceAt(ctx, addr)
//	// build tx, sign, broadcast...
func (nm *NonceManager) LockAddress(addr common.Address) (release func()) {
	nm.mu.Lock()
	m, ok := nm.locks[addr]
	if !ok {
		m = &sync.Mutex{}
		nm.locks[addr] = m
	}
	nm.mu.Unlock()

	m.Lock()
	return m.Unlock
}

// NextNonce берёт текущий nonce из mempool. Должен вызываться ВНУТРИ
// LockAddress, иначе двое параллельных запросов получат одинаковое значение.
//
// PendingNonceAt включает транзакции, которые уже в mempool но ещё не
// в блоке — это то что нам нужно, чтобы две быстрых операции получили
// разные nonce.
func (nm *NonceManager) NextNonce(ctx context.Context, client *Client, addr common.Address) (uint64, error) {
	nonce, err := client.Eth().PendingNonceAt(ctx, addr)
	if err != nil {
		return 0, fmt.Errorf("blockchain: pending nonce for %s: %w", addr.Hex(), err)
	}
	return nonce, nil
}