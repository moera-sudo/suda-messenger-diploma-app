package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

// Reader — обёртка над view-вызовами контрактов с консистентным API.
// Используется фичами для on-demand чтения on-chain состояния (балансы,
// владельцы NFT, статусы квестов и т.п.).
//
// Все методы Reader НЕ обращаются к БД — только к чейну. Если нужен кеш,
// он живёт в repository-слое выше.
type Reader struct {
	client    *Client
	contracts *Contracts
}

func NewReader(client *Client, contracts *Contracts) *Reader {
	return &Reader{client: client, contracts: contracts}
}

// ────────────────────────────────────────────────────────────
//  SudaToken (ERC-20)
// ────────────────────────────────────────────────────────────

// SudaBalanceOf — текущий баланс SUDA по адресу.
// Возвращает значение в wei (uint256 как *big.Int).
func (r *Reader) SudaBalanceOf(ctx context.Context, addr common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	bal, err := r.contracts.Token.BalanceOf(opts, addr)
	if err != nil {
		return nil, fmt.Errorf("blockchain: balanceOf %s: %w", addr.Hex(), err)
	}
	return bal, nil
}

// SudaTotalSupply — суммарная эмиссия SUDA на чейне.
// Используется для метрик / админ-панели.
func (r *Reader) SudaTotalSupply(ctx context.Context) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return r.contracts.Token.TotalSupply(opts)
}

// SudaAllowance — сколько SUDA `owner` разрешил тратить `spender`.
// Нужно для marketplace/escrow/fundraising — они тянут SUDA через
// transferFrom, который требует предварительного approve.
func (r *Reader) SudaAllowance(ctx context.Context, owner, spender common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return r.contracts.Token.Allowance(opts, owner, spender)
}

// ────────────────────────────────────────────────────────────
//  SudaNFT (ERC-721)
// ────────────────────────────────────────────────────────────

// NFTOwnerOf — текущий владелец NFT с данным tokenId.
// Если токен не существует — контракт reverts (ошибка возвращается).
func (r *Reader) NFTOwnerOf(ctx context.Context, tokenID *big.Int) (common.Address, error) {
	opts := &bind.CallOpts{Context: ctx}
	owner, err := r.contracts.NFT.OwnerOf(opts, tokenID)
	if err != nil {
		return common.Address{}, fmt.Errorf("blockchain: ownerOf token %s: %w", tokenID.String(), err)
	}
	return owner, nil
}

// NFTTokenURI — metadata URI токена (HTTP-ссылка к нашему transaction API).
func (r *Reader) NFTTokenURI(ctx context.Context, tokenID *big.Int) (string, error) {
	opts := &bind.CallOpts{Context: ctx}
	return r.contracts.NFT.TokenURI(opts, tokenID)
}

// NFTBalanceOf — сколько NFT у адреса.
func (r *Reader) NFTBalanceOf(ctx context.Context, addr common.Address) (*big.Int, error) {
	opts := &bind.CallOpts{Context: ctx}
	return r.contracts.NFT.BalanceOf(opts, addr)
}

// ────────────────────────────────────────────────────────────
//  Generic helpers
// ────────────────────────────────────────────────────────────

// CodeAt — есть ли байткод по адресу. Используется например при
// валидации что transaction-service подключился к чейну, где
// контракты действительно задеплоены.
func (r *Reader) CodeAt(ctx context.Context, addr common.Address) ([]byte, error) {
	return r.client.Eth().CodeAt(ctx, addr, nil)
}

// HasContractCode — true если по адресу есть байткод.
func (r *Reader) HasContractCode(ctx context.Context, addr common.Address) (bool, error) {
	code, err := r.CodeAt(ctx, addr)
	if err != nil {
		return false, err
	}
	return len(code) > 0, nil
}