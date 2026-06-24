package blockchain

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/rs/zerolog/log"
)

// Client — обёртка над ethclient.Client с явным chainID и контекстом.
// Используется всеми остальными компонентами blockchain-слоя (reader,
// broadcaster, observer) как единая точка доступа к Besu.
type Client struct {
	eth     *ethclient.Client
	chainID *big.Int
	rpcURL  string
}

// New подключается к Besu по HTTP-RPC и проверяет, что chainID
// совпадает с ожидаемым (защищает от случайной работы с чужим чейном).
func New(ctx context.Context, rpcURL string, expectedChainID int64) (*Client, error) {
	eth, err := ethclient.DialContext(ctx, rpcURL)
	if err != nil {
		return nil, fmt.Errorf("blockchain: dial %s: %w", rpcURL, err)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	actualChainID, err := eth.ChainID(pingCtx)
	if err != nil {
		eth.Close()
		return nil, fmt.Errorf("blockchain: chainID call failed: %w", err)
	}

	if actualChainID.Int64() != expectedChainID {
		eth.Close()
		return nil, fmt.Errorf(
			"blockchain: chainID mismatch — expected %d, got %d (wrong RPC endpoint?)",
			expectedChainID, actualChainID.Int64(),
		)
	}

	log.Info().
		Str("rpc", rpcURL).
		Int64("chain_id", actualChainID.Int64()).
		Msg("blockchain client connected")

	return &Client{
		eth:     eth,
		chainID: actualChainID,
		rpcURL:  rpcURL,
	}, nil
}

// Eth возвращает нижележащий ethclient.Client для случаев, когда нужны
// специфичные методы (FilterLogs для observer, SubscribeNewHead и т.п.).
// Везде, где возможно — используй высокоуровневые reader/broadcaster.
func (c *Client) Eth() *ethclient.Client {
	return c.eth
}

// ChainID возвращает chain ID сети (1337 для Suda Besu).
func (c *Client) ChainID() *big.Int {
	return new(big.Int).Set(c.chainID)
}

// BlockNumber возвращает номер последнего блока.
func (c *Client) BlockNumber(ctx context.Context) (uint64, error) {
	return c.eth.BlockNumber(ctx)
}

// Close закрывает RPC-соединение. Вызывается при graceful shutdown.
func (c *Client) Close() {
	if c.eth != nil {
		c.eth.Close()
	}
}