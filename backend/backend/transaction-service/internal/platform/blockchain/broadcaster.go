package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
)

// Broadcaster инкапсулирует логику подготовки и отправки транзакции
// через биндинги контрактов.
//
// Типичный flow вызова контракта из feature-слоя:
//
//	opts, release, err := broadcaster.PrepareOpts(ctx, userSigner)
//	if err != nil { return err }
//	defer release()
//
//	tx, err := contracts.Token.Transfer(opts, recipientAddr, amount)
//	if err != nil { return err }
//
//	receipt, err := broadcaster.WaitForReceipt(ctx, tx.Hash())
type Broadcaster struct {
	client *Client
	nm     *NonceManager
}

func NewBroadcaster(client *Client, nm *NonceManager) *Broadcaster {
	return &Broadcaster{client: client, nm: nm}
}

// PrepareOpts собирает *bind.TransactOpts для подписи транзакции:
//
//   - блокирует адрес в NonceManager (чтобы избежать гонок),
//   - читает актуальный nonce,
//   - выставляет gasPrice = 0 (Besu zeroBaseFee),
//   - проставляет sign-функцию от signer'а.
//
// Возвращает opts И release-функцию. Release ОБЯЗАТЕЛЬНО вызвать
// после отправки транзакции (defer release()), иначе адрес останется
// залоченным и следующие операции зависнут.
func (b *Broadcaster) PrepareOpts(ctx context.Context, signer Signer) (*bind.TransactOpts, func(), error) {
	addr := signer.Address()

	// 1. Блокируем все будущие подписи для этого адреса до release()
	release := b.nm.LockAddress(addr)

	// 2. Берём свежий nonce из mempool
	nonce, err := b.nm.NextNonce(ctx, b.client, addr)
	if err != nil {
		release()
		return nil, nil, err
	}

	// 3. Получаем готовые TransactOpts от signer'а
	opts, err := signer.TransactOpts(ctx, b.client.ChainID())
	if err != nil {
		release()
		return nil, nil, err
	}

	// 4. Настраиваем параметры
	opts.Context = ctx
	opts.Nonce = new(big.Int).SetUint64(nonce)
	opts.GasPrice = big.NewInt(0) // Besu private chain, zeroBaseFee
	opts.GasLimit = 0             // 0 = пусть оценит сама ethclient через estimateGas

	return opts, release, nil
}

// SendTokenTransfer — high-level helper: подписать и отправить SudaToken.Transfer
// одним вызовом. Обёртка вокруг PrepareOpts + Token.Transfer.
//
// Прямой вызов: для редких мест где Broadcaster и Contracts инжектируются раздельно.
// Для feature-слоя есть TokenSender (см. ниже) — он обёртывает Broadcaster+Contracts
// и удовлетворяет маленькому интерфейсу, который легко мокается.
func (b *Broadcaster) SendTokenTransfer(
	ctx context.Context,
	contracts *Contracts,
	signer Signer,
	to common.Address,
	amount *big.Int,
) (string, error) {
	opts, release, err := b.PrepareOpts(ctx, signer)
	if err != nil {
		return "", err
	}
	defer release()

	tx, err := contracts.Token.Transfer(opts, to, amount)
	if err != nil {
		return "", err
	}
	return tx.Hash().Hex(), nil
}

// TokenSender — обёртка Broadcaster + Contracts с компактным интерфейсом.
// Внедряется в feature-слой как единое зависимость (вместо двух отдельных),
// что позволяет одним моком накрывать всю broadcast-логику в тестах.
type TokenSender struct {
	broadcaster *Broadcaster
	contracts   *Contracts
}

func NewTokenSender(b *Broadcaster, c *Contracts) *TokenSender {
	return &TokenSender{broadcaster: b, contracts: c}
}

// SendTokenTransfer подписывает и отправляет SudaToken.Transfer от лица signer'а.
// Возвращает hex tx_hash.
func (s *TokenSender) SendTokenTransfer(
	ctx context.Context, signer Signer, to common.Address, amount *big.Int,
) (string, error) {
	return s.broadcaster.SendTokenTransfer(ctx, s.contracts, signer, to, amount)
}

// WaitForReceipt ждёт пока транзакция попадёт в блок и вернёт receipt.
// Опрашивает каждые 500мс, максимум 30 секунд (= 15 блоков при block time 2 сек).
//
// Если transaction reverted (receipt.Status == 0) — возвращает ошибку,
// но receipt не nil — вызывающий может проверить причину через trace.
func (b *Broadcaster) WaitForReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	const (
		pollInterval = 500 * time.Millisecond
		maxAttempts  = 60 // 60 * 500ms = 30 секунд
	)

	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()

	for attempt := 0; attempt < maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
		}

		receipt, err := b.client.Eth().TransactionReceipt(ctx, txHash)
		if err == nil {
			if receipt.Status == types.ReceiptStatusFailed {
				log.Warn().
					Str("tx_hash", txHash.Hex()).
					Uint64("block", receipt.BlockNumber.Uint64()).
					Msg("transaction reverted on-chain")
				return receipt, fmt.Errorf("blockchain: tx %s reverted", txHash.Hex())
			}
			return receipt, nil
		}
		// Receipt ещё не появился — продолжаем опрашивать.
	}

	return nil, fmt.Errorf("blockchain: tx %s not mined within timeout", txHash.Hex())
}