// Package observer — фоновый процесс, который читает on-chain event logs
// из Besu и обновляет PostgreSQL read model (tx_transactions, tx_pending).
package observer

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"

	purchaserepo "github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/observer/handlers"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/blockchain"
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"
)

// Контракты которые слушаем. Имя — ключ в tx_observer_state.
const (
	ContractSudaToken       = "SUDA_TOKEN"
	ContractSudaNFT         = "SUDA_NFT"
	ContractSudaMarketplace = "SUDA_MARKETPLACE"
	ContractSudaEscrow      = "SUDA_ESCROW"
	ContractSudaFundraising = "SUDA_FUNDRAISING"
)

// Observer — главный struct.
type Observer struct {
	repo         *repository.Repository
	chain        *blockchain.Client
	contracts    *blockchain.Contracts
	messenger    *platgrpc.MessengerClient
	pollInterval time.Duration

	tokenHandler *handlers.SudaTokenHandler
}

// Deps — зависимости observer'а.
type Deps struct {
	Postgres        *pgxpool.Pool
	Chain           *blockchain.Client
	Contracts       *blockchain.Contracts
	MessengerClient *platgrpc.MessengerClient
}

const maxBlockRange = 1000

func New(d Deps) *Observer {
	repo := repository.New(d.Postgres)
	purchases := purchaserepo.New(d.Postgres)
	return &Observer{
		repo:         repo,
		chain:        d.Chain,
		contracts:    d.Contracts,
		messenger:    d.MessengerClient,
		pollInterval: 2 * time.Second,

		tokenHandler: handlers.NewSudaTokenHandler(handlers.SudaTokenDeps{
			Repo:            repo,
			PurchaseRepo:    purchases,
			Chain:           d.Chain,
			Contracts:       d.Contracts,
			MessengerClient: d.MessengerClient,
		}),
	}
}

// Run блокирующе крутит poll-цикл до отмены ctx.
func (o *Observer) Run(ctx context.Context) {
	log.Info().Dur("poll_interval", o.pollInterval).Msg("observer starting")

	if err := o.initStateIfMissing(ctx); err != nil {
		log.Error().Err(err).Msg("observer init state failed (will retry)")
	}

	ticker := time.NewTicker(o.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Info().Msg("observer stopping")
			return
		case <-ticker.C:
			if err := o.tick(ctx); err != nil {
				log.Error().Err(err).Msg("observer tick failed")
			}
		}
	}
}

// initStateIfMissing — для каждого контракта, если в БД нет записи,
// ставим last_processed_block = текущий блок чейна (старт с "сейчас",
// не пытаемся пройти историю при первом запуске).
func (o *Observer) initStateIfMissing(ctx context.Context) error {
	currentBlock, err := o.chain.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("get current block: %w", err)
	}

	contracts := []string{
		ContractSudaToken,
		ContractSudaNFT,
		ContractSudaMarketplace,
		ContractSudaEscrow,
		ContractSudaFundraising,
	}

	for _, name := range contracts {
		_, err := o.repo.GetObserverState(ctx, name)
		if errors.Is(err, repository.ErrNotFound) {
			if err := o.repo.UpsertObserverState(ctx, name, currentBlock); err != nil {
				return fmt.Errorf("init state for %s: %w", name, err)
			}
			log.Info().
				Str("contract", name).
				Uint64("start_block", currentBlock).
				Msg("observer state initialized")
		} else if err != nil {
			return fmt.Errorf("check state for %s: %w", name, err)
		}
	}
	return nil
}

// tick — одна итерация poll-цикла.
func (o *Observer) tick(ctx context.Context) error {
	currentBlock, err := o.chain.BlockNumber(ctx)
	if err != nil {
		return fmt.Errorf("get current block: %w", err)
	}

	// SudaToken — основной контракт итерации В.
	if err := o.processContract(ctx, ContractSudaToken, currentBlock, o.tokenHandler.HandleLogs); err != nil {
		log.Error().Err(err).Str("contract", ContractSudaToken).Msg("process logs failed")
	}

	// Остальные контракты — placeholders на итерацию Г.
	// Просто двигаем cursor вперёд (нечего обрабатывать).
	for _, name := range []string{
		ContractSudaNFT, ContractSudaMarketplace, ContractSudaEscrow, ContractSudaFundraising,
	} {
		if err := o.repo.UpsertObserverState(ctx, name, currentBlock); err != nil {
			log.Warn().Err(err).Str("contract", name).Msg("update cursor failed")
		}
	}

	return nil
}

// processContract — обработка одного контракта.
//
// handler — функция которая обрабатывает logs пакетом. Если она вернёт ошибку —
// НЕ продвигаем cursor (попробуем снова на следующем тике).
func (o *Observer) processContract(
	ctx context.Context,
	contractName string,
	currentBlock uint64,
	handler func(ctx context.Context, fromBlock, toBlock uint64) error,
) error {
	last, err := o.repo.GetObserverState(ctx, contractName)
	if err != nil {
		return fmt.Errorf("get state: %w", err)
	}

	if last >= currentBlock {
		return nil
	}

	fromBlock := last + 1
	toBlock := currentBlock

	// Ограничиваем размер окна — если разрыв большой, доедем в следующих тиках.
	if toBlock-fromBlock+1 > maxBlockRange {
		toBlock = fromBlock + maxBlockRange - 1
		log.Debug().
			Str("contract", contractName).
			Uint64("from_block", fromBlock).
			Uint64("to_block", toBlock).
			Uint64("current_block", currentBlock).
			Uint64("blocks_behind", currentBlock-toBlock).
			Msg("processing chunk (will catch up over next ticks)")
	}

	if err := handler(ctx, fromBlock, toBlock); err != nil {
		return fmt.Errorf("handler %s [%d-%d]: %w", contractName, fromBlock, toBlock, err)
	}

	if err := o.repo.UpsertObserverState(ctx, contractName, toBlock); err != nil {
		return fmt.Errorf("upsert state: %w", err)
	}

	return nil
}