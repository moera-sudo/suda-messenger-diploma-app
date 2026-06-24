package service_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	walletmodels "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	walletsvc "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/service"
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"
)

// ── Inline mocks ───────────────────────────────────────────

type fakeTreasuryRepo struct {
	wallet    *walletmodels.ChannelWallet
	walletErr error

	stats    func() (int64, *big.Int, error)
	topDonors []repository.TopDonor
	donations []repository.DonationRecord
}

func (r *fakeTreasuryRepo) GetChannelWallet(_ context.Context, _ uuid.UUID) (*walletmodels.ChannelWallet, error) {
	if r.walletErr != nil {
		return nil, r.walletErr
	}
	if r.wallet != nil {
		return r.wallet, nil
	}
	return nil, repository.ErrNotFound
}

func (r *fakeTreasuryRepo) GetChannelDonationStats(ctx context.Context, _ uuid.UUID) (int64, *big.Int, error) {
	if r.stats != nil {
		return r.stats()
	}
	return 0, big.NewInt(0), nil
}

func (r *fakeTreasuryRepo) GetChannelTopDonors(_ context.Context, _ uuid.UUID, _ int) ([]repository.TopDonor, error) {
	return r.topDonors, nil
}

func (r *fakeTreasuryRepo) GetChannelDonations(_ context.Context, _ uuid.UUID, _, _ int) ([]repository.DonationRecord, error) {
	return r.donations, nil
}

type fakeTreasuryReader struct {
	balance *big.Int
	err     error
}

func (r *fakeTreasuryReader) SudaBalanceOf(_ context.Context, _ common.Address) (*big.Int, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.balance != nil {
		return r.balance, nil
	}
	return big.NewInt(0), nil
}

type fakeTreasuryMessenger struct {
	permGranted bool
	permErr     error
	users       map[string]platgrpc.UserBrief
}

func (m *fakeTreasuryMessenger) CheckChannelPermission(_ context.Context, _, _, _ string) (*platgrpc.ChannelPermissionResult, error) {
	if m.permErr != nil {
		return nil, m.permErr
	}
	return &platgrpc.ChannelPermissionResult{Granted: m.permGranted, Reason: "ok"}, nil
}

func (m *fakeTreasuryMessenger) GetUsersByIDs(_ context.Context, _ []string) (map[string]platgrpc.UserBrief, error) {
	if m.users != nil {
		return m.users, nil
	}
	return map[string]platgrpc.UserBrief{}, nil
}

// ── Setup helper ───────────────────────────────────────────

type treasuryEnv struct {
	svc       *walletsvc.TreasuryService
	repo      *fakeTreasuryRepo
	reader    *fakeTreasuryReader
	messenger *fakeTreasuryMessenger
}

func newTreasuryEnv(permGranted bool) *treasuryEnv {
	repo := &fakeTreasuryRepo{}
	reader := &fakeTreasuryReader{balance: big.NewInt(0)}
	mc := &fakeTreasuryMessenger{permGranted: permGranted}
	svc := walletsvc.NewTreasuryForTest(repo, reader, mc)
	return &treasuryEnv{svc: svc, repo: repo, reader: reader, messenger: mc}
}

// ── GetChannelTreasury ─────────────────────────────────────

func TestGetChannelTreasury_OwnerOk(t *testing.T) {
	e := newTreasuryEnv(true)
	channelID := uuid.New()

	e.repo.wallet = &walletmodels.ChannelWallet{
		ChannelID: channelID,
		Address:   "0x000000000000000000000000000000000000beef",
	}
	e.reader.balance = new(big.Int).Mul(big.NewInt(500), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	donorID := uuid.New()
	donorWei, _ := new(big.Int).SetString("50000000000000000000", 10)
	e.repo.topDonors = []repository.TopDonor{
		{FromUserID: donorID, TotalWei: donorWei, DonationCount: 3},
	}
	e.repo.stats = func() (int64, *big.Int, error) {
		totalWei, _ := new(big.Int).SetString("50000000000000000000", 10)
		return 3, totalWei, nil
	}
	e.messenger.users = map[string]platgrpc.UserBrief{
		donorID.String(): {UserID: donorID.String(), Username: "alice", DisplayName: "Alice"},
	}

	stats, err := e.svc.GetChannelTreasury(context.Background(), uuid.New(), channelID)
	require.NoError(t, err)
	assert.Equal(t, channelID.String(), stats.ChannelID)
	assert.Equal(t, "0x000000000000000000000000000000000000beef", stats.Address)
	assert.Equal(t, int64(3), stats.TotalDonationsCount)
	require.Len(t, stats.TopDonors, 1)
	assert.Equal(t, "alice", stats.TopDonors[0].Username)
	assert.Equal(t, int64(3), stats.TopDonors[0].DonationCount)
}

func TestGetChannelTreasury_NotOwner(t *testing.T) {
	e := newTreasuryEnv(false)

	_, err := e.svc.GetChannelTreasury(context.Background(), uuid.New(), uuid.New())
	assert.ErrorIs(t, err, walletsvc.ErrForbidden)
}

func TestGetChannelTreasury_NoChannelWallet(t *testing.T) {
	e := newTreasuryEnv(true)
	// repo.wallet == nil, returns ErrNotFound

	_, err := e.svc.GetChannelTreasury(context.Background(), uuid.New(), uuid.New())
	assert.ErrorIs(t, err, walletsvc.ErrWalletNotFound)
}

func TestGetChannelTreasury_NoDonations(t *testing.T) {
	e := newTreasuryEnv(true)
	e.repo.wallet = &walletmodels.ChannelWallet{ChannelID: uuid.New(), Address: "0xdead"}
	// stats returns 0,0 by default; topDonors empty

	stats, err := e.svc.GetChannelTreasury(context.Background(), uuid.New(), uuid.New())
	require.NoError(t, err)
	assert.Equal(t, int64(0), stats.TotalDonationsCount)
	assert.Equal(t, "0", stats.TotalDonationsWei)
	assert.Empty(t, stats.TopDonors)
}

// ── GetChannelDonations ────────────────────────────────────

func TestGetChannelDonations_OK(t *testing.T) {
	e := newTreasuryEnv(true)
	channelID := uuid.New()
	donorID := uuid.New()
	e.repo.stats = func() (int64, *big.Int, error) {
		return 2, big.NewInt(2000), nil
	}
	e.repo.donations = []repository.DonationRecord{
		{
			ID:          uuid.New(),
			FromUserID:  donorID,
			FromAddress: "0xcafe",
			Amount:      big.NewInt(1000),
			Message:     "nice work",
			TxHash:      "0xabc",
			CreatedAt:   time.Now(),
		},
	}
	e.messenger.users = map[string]platgrpc.UserBrief{
		donorID.String(): {UserID: donorID.String(), Username: "bob", DisplayName: "Bob"},
	}

	list, err := e.svc.GetChannelDonations(context.Background(), uuid.New(), channelID, 10, 0)
	require.NoError(t, err)
	assert.Equal(t, int64(2), list.Total)
	require.Len(t, list.Items, 1)
	assert.Equal(t, "bob", list.Items[0].FromUsername)
	assert.Equal(t, "1000", list.Items[0].AmountWei)
	assert.Equal(t, "nice work", list.Items[0].Message)
}

func TestGetChannelDonations_NotOwner(t *testing.T) {
	e := newTreasuryEnv(false)

	_, err := e.svc.GetChannelDonations(context.Background(), uuid.New(), uuid.New(), 10, 0)
	assert.ErrorIs(t, err, walletsvc.ErrForbidden)
}
