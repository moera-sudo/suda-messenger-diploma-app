package mocks

import (
	"context"
	"math/big"

	"github.com/google/uuid"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
)

// MockWalletRepo покрывает все методы wallet/repository, используемые сервисами
// purchase, donation, gating, treasury и observer. Если в продуктовом repo
// добавится новый метод — его нужно продублировать сюда.
type MockWalletRepo struct {
	// Wallets
	GetUserWalletFn          func(ctx context.Context, userID uuid.UUID) (*wallet.Wallet, error)
	GetUserWalletByAddressFn func(ctx context.Context, address string) (*wallet.Wallet, error)
	GetChannelWalletFn       func(ctx context.Context, channelID uuid.UUID) (*wallet.ChannelWallet, error)

	// Pending + audit
	InsertPendingFn         func(ctx context.Context, txHash string, fromUserID uuid.UUID, expectedType string, chatID uuid.UUID) error
	InsertPendingDonationFn func(ctx context.Context, txHash string, fromUserID uuid.UUID, chatID uuid.UUID, donationMessage string) error
	WriteAuditFn            func(ctx context.Context, e repository.AuditEntry) error

	// Gating
	GetGatingRuleFn    func(ctx context.Context, chatID uuid.UUID) (*repository.GatingRule, error)
	UpsertGatingRuleFn func(ctx context.Context, chatID uuid.UUID, minSudaBalance *big.Int, nftCollectionID *uuid.UUID, subscriptionPrice *big.Int, createdBy uuid.UUID) error
	DeleteGatingRuleFn func(ctx context.Context, chatID uuid.UUID) error
	UserOwnsNFTFn      func(ctx context.Context, userID, collectionID uuid.UUID) (bool, error)

	// Treasury reads
	GetChannelDonationStatsFn func(ctx context.Context, channelID uuid.UUID) (int64, *big.Int, error)
	GetChannelTopDonorsFn     func(ctx context.Context, channelID uuid.UUID, limit int) ([]repository.TopDonor, error)
	GetChannelDonationsFn     func(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]repository.DonationRecord, error)

	// Observer-side
	InsertTransactionFn  func(ctx context.Context, p repository.TxInsertParams) (bool, error)
	InsertDonationFn     func(ctx context.Context, p repository.DonationInsertParams) error
	DeletePendingFn      func(ctx context.Context, txHash string) error
	GetPendingMetaFn     func(ctx context.Context, txHash string) (*repository.PendingMeta, error)
	FindUserByAddressFn  func(ctx context.Context, address string) (*uuid.UUID, error)
	FindChannelByAddressFn func(ctx context.Context, address string) (*uuid.UUID, error)

	// Counters (по мере надобности)
	InsertPendingCalls         int
	InsertPendingDonationCalls int
	WriteAuditCalls            int
	UpsertGatingRuleCalls      int
	DeleteGatingRuleCalls      int
	InsertTransactionCalls     int
	InsertDonationCalls        int
	DeletePendingCalls         int
}

func NewMockWalletRepo() *MockWalletRepo { return &MockWalletRepo{} }

// ── Wallets ───────────────────────────────────────────────────

func (m *MockWalletRepo) GetUserWallet(ctx context.Context, userID uuid.UUID) (*wallet.Wallet, error) {
	if m.GetUserWalletFn != nil {
		return m.GetUserWalletFn(ctx, userID)
	}
	return nil, repository.ErrNotFound
}

func (m *MockWalletRepo) GetUserWalletByAddress(ctx context.Context, address string) (*wallet.Wallet, error) {
	if m.GetUserWalletByAddressFn != nil {
		return m.GetUserWalletByAddressFn(ctx, address)
	}
	return nil, repository.ErrNotFound
}

func (m *MockWalletRepo) GetChannelWallet(ctx context.Context, channelID uuid.UUID) (*wallet.ChannelWallet, error) {
	if m.GetChannelWalletFn != nil {
		return m.GetChannelWalletFn(ctx, channelID)
	}
	return nil, repository.ErrNotFound
}

// ── Pending + audit ───────────────────────────────────────────

func (m *MockWalletRepo) InsertPending(ctx context.Context, txHash string, fromUserID uuid.UUID, expectedType string, chatID uuid.UUID) error {
	m.InsertPendingCalls++
	if m.InsertPendingFn != nil {
		return m.InsertPendingFn(ctx, txHash, fromUserID, expectedType, chatID)
	}
	return nil
}

func (m *MockWalletRepo) InsertPendingDonation(ctx context.Context, txHash string, fromUserID uuid.UUID, chatID uuid.UUID, donationMessage string) error {
	m.InsertPendingDonationCalls++
	if m.InsertPendingDonationFn != nil {
		return m.InsertPendingDonationFn(ctx, txHash, fromUserID, chatID, donationMessage)
	}
	return nil
}

func (m *MockWalletRepo) WriteAudit(ctx context.Context, e repository.AuditEntry) error {
	m.WriteAuditCalls++
	if m.WriteAuditFn != nil {
		return m.WriteAuditFn(ctx, e)
	}
	return nil
}

// ── Gating ────────────────────────────────────────────────────

func (m *MockWalletRepo) GetGatingRule(ctx context.Context, chatID uuid.UUID) (*repository.GatingRule, error) {
	if m.GetGatingRuleFn != nil {
		return m.GetGatingRuleFn(ctx, chatID)
	}
	return nil, repository.ErrNotFound
}

func (m *MockWalletRepo) UpsertGatingRule(ctx context.Context, chatID uuid.UUID, minSudaBalance *big.Int, nftCollectionID *uuid.UUID, subscriptionPrice *big.Int, createdBy uuid.UUID) error {
	m.UpsertGatingRuleCalls++
	if m.UpsertGatingRuleFn != nil {
		return m.UpsertGatingRuleFn(ctx, chatID, minSudaBalance, nftCollectionID, subscriptionPrice, createdBy)
	}
	return nil
}

func (m *MockWalletRepo) DeleteGatingRule(ctx context.Context, chatID uuid.UUID) error {
	m.DeleteGatingRuleCalls++
	if m.DeleteGatingRuleFn != nil {
		return m.DeleteGatingRuleFn(ctx, chatID)
	}
	return nil
}

func (m *MockWalletRepo) UserOwnsNFTFromCollection(ctx context.Context, userID, collectionID uuid.UUID) (bool, error) {
	if m.UserOwnsNFTFn != nil {
		return m.UserOwnsNFTFn(ctx, userID, collectionID)
	}
	return false, nil
}

// ── Treasury reads ────────────────────────────────────────────

func (m *MockWalletRepo) GetChannelDonationStats(ctx context.Context, channelID uuid.UUID) (int64, *big.Int, error) {
	if m.GetChannelDonationStatsFn != nil {
		return m.GetChannelDonationStatsFn(ctx, channelID)
	}
	return 0, big.NewInt(0), nil
}

func (m *MockWalletRepo) GetChannelTopDonors(ctx context.Context, channelID uuid.UUID, limit int) ([]repository.TopDonor, error) {
	if m.GetChannelTopDonorsFn != nil {
		return m.GetChannelTopDonorsFn(ctx, channelID, limit)
	}
	return nil, nil
}

func (m *MockWalletRepo) GetChannelDonations(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]repository.DonationRecord, error) {
	if m.GetChannelDonationsFn != nil {
		return m.GetChannelDonationsFn(ctx, channelID, limit, offset)
	}
	return nil, nil
}

// ── Observer-side ─────────────────────────────────────────────

func (m *MockWalletRepo) InsertTransaction(ctx context.Context, p repository.TxInsertParams) (bool, error) {
	m.InsertTransactionCalls++
	if m.InsertTransactionFn != nil {
		return m.InsertTransactionFn(ctx, p)
	}
	return true, nil
}

func (m *MockWalletRepo) InsertDonation(ctx context.Context, p repository.DonationInsertParams) error {
	m.InsertDonationCalls++
	if m.InsertDonationFn != nil {
		return m.InsertDonationFn(ctx, p)
	}
	return nil
}

func (m *MockWalletRepo) DeletePending(ctx context.Context, txHash string) error {
	m.DeletePendingCalls++
	if m.DeletePendingFn != nil {
		return m.DeletePendingFn(ctx, txHash)
	}
	return nil
}

func (m *MockWalletRepo) GetPendingMeta(ctx context.Context, txHash string) (*repository.PendingMeta, error) {
	if m.GetPendingMetaFn != nil {
		return m.GetPendingMetaFn(ctx, txHash)
	}
	return nil, repository.ErrNotFound
}

func (m *MockWalletRepo) FindUserByAddress(ctx context.Context, address string) (*uuid.UUID, error) {
	if m.FindUserByAddressFn != nil {
		return m.FindUserByAddressFn(ctx, address)
	}
	return nil, repository.ErrNotFound
}

func (m *MockWalletRepo) FindChannelByAddress(ctx context.Context, address string) (*uuid.UUID, error) {
	if m.FindChannelByAddressFn != nil {
		return m.FindChannelByAddressFn(ctx, address)
	}
	return nil, repository.ErrNotFound
}
