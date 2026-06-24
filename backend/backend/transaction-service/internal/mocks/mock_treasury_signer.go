package mocks

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

// MockTreasurySigner — реализация blockchain.Signer для тестов.
// В тестах не используется для реальной подписи (broadcast мокается через
// MockTokenSender), но интерфейсу удовлетворять должен.
type MockTreasurySigner struct {
	Addr common.Address
}

func NewMockTreasurySigner() *MockTreasurySigner {
	return &MockTreasurySigner{
		Addr: common.HexToAddress("0x0000000000000000000000000000000000000001"),
	}
}

func (s *MockTreasurySigner) Address() common.Address { return s.Addr }

func (s *MockTreasurySigner) SignTx(_ context.Context, tx *types.Transaction, _ *big.Int) (*types.Transaction, error) {
	return tx, nil
}

func (s *MockTreasurySigner) TransactOpts(_ context.Context, _ *big.Int) (*bind.TransactOpts, error) {
	return &bind.TransactOpts{From: s.Addr}, nil
}
