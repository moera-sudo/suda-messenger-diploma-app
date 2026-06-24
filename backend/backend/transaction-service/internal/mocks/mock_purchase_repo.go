// Package mocks — hand-written test doubles для transaction-service.
//
// Паттерн: каждый мок — struct с XxxFn полями (для кастомного поведения)
// и XxxCalls счётчиками (для assertion на количество вызовов).
package mocks

import (
	"context"

	"github.com/google/uuid"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase/repository"
)

type MockPurchaseRepo struct {
	CreateFn         func(ctx context.Context, p repository.CreateParams) (uuid.UUID, error)
	GetByIDFn        func(ctx context.Context, id uuid.UUID) (*purchase.Purchase, error)
	MarkProcessingFn func(ctx context.Context, id uuid.UUID) error
	MarkCompletedFn  func(ctx context.Context, id uuid.UUID, txHash string) error
	MarkFailedFn     func(ctx context.Context, id uuid.UUID, reason string) error
	ListForUserFn    func(ctx context.Context, userID uuid.UUID, limit, offset int) ([]purchase.Purchase, error)
	FindByTxHashFn   func(ctx context.Context, txHash string) (*repository.TxHashMatch, error)

	CreateCalls         int
	GetByIDCalls        int
	MarkProcessingCalls int
	MarkCompletedCalls  int
	MarkFailedCalls     int
	ListForUserCalls    int
	FindByTxHashCalls   int
}

func NewMockPurchaseRepo() *MockPurchaseRepo { return &MockPurchaseRepo{} }

func (m *MockPurchaseRepo) Create(ctx context.Context, p repository.CreateParams) (uuid.UUID, error) {
	m.CreateCalls++
	if m.CreateFn != nil {
		return m.CreateFn(ctx, p)
	}
	return uuid.New(), nil
}

func (m *MockPurchaseRepo) GetByID(ctx context.Context, id uuid.UUID) (*purchase.Purchase, error) {
	m.GetByIDCalls++
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, repository.ErrNotFound
}

func (m *MockPurchaseRepo) MarkProcessing(ctx context.Context, id uuid.UUID) error {
	m.MarkProcessingCalls++
	if m.MarkProcessingFn != nil {
		return m.MarkProcessingFn(ctx, id)
	}
	return nil
}

func (m *MockPurchaseRepo) MarkCompleted(ctx context.Context, id uuid.UUID, txHash string) error {
	m.MarkCompletedCalls++
	if m.MarkCompletedFn != nil {
		return m.MarkCompletedFn(ctx, id, txHash)
	}
	return nil
}

func (m *MockPurchaseRepo) MarkFailed(ctx context.Context, id uuid.UUID, reason string) error {
	m.MarkFailedCalls++
	if m.MarkFailedFn != nil {
		return m.MarkFailedFn(ctx, id, reason)
	}
	return nil
}

func (m *MockPurchaseRepo) ListForUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]purchase.Purchase, error) {
	m.ListForUserCalls++
	if m.ListForUserFn != nil {
		return m.ListForUserFn(ctx, userID, limit, offset)
	}
	return nil, nil
}

func (m *MockPurchaseRepo) FindByTxHash(ctx context.Context, txHash string) (*repository.TxHashMatch, error) {
	m.FindByTxHashCalls++
	if m.FindByTxHashFn != nil {
		return m.FindByTxHashFn(ctx, txHash)
	}
	return nil, repository.ErrNotFound
}
