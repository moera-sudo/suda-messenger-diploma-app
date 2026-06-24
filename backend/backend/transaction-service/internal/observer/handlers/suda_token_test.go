package handlers

import (
	"context"
	"errors"
	"math/big"
	"sync/atomic"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moera-sudo/backend/backend/transaction-service/contracts/sudatoken"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/donation"
	purchaserepo "github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"
)

// ── Inline mocks ─────────────────────────────────────────────

type mockObserverRepo struct {
	insertResult bool
	insertErr    error

	pendingMeta *repository.PendingMeta
	pendingErr  error

	findUserFn    func(address string) (*uuid.UUID, error)
	findChannelFn func(address string) (*uuid.UUID, error)

	insertTransactionCalls int
	insertDonationCalls    int
	deletePendingCalls     int
}

func (r *mockObserverRepo) GetPendingMeta(_ context.Context, _ string) (*repository.PendingMeta, error) {
	if r.pendingErr != nil {
		return nil, r.pendingErr
	}
	if r.pendingMeta != nil {
		return r.pendingMeta, nil
	}
	return nil, repository.ErrNotFound
}

func (r *mockObserverRepo) InsertTransaction(_ context.Context, _ repository.TxInsertParams) (bool, error) {
	r.insertTransactionCalls++
	return r.insertResult, r.insertErr
}

func (r *mockObserverRepo) DeletePending(_ context.Context, _ string) error {
	r.deletePendingCalls++
	return nil
}

func (r *mockObserverRepo) InsertDonation(_ context.Context, _ repository.DonationInsertParams) error {
	r.insertDonationCalls++
	return nil
}

func (r *mockObserverRepo) FindUserByAddress(_ context.Context, address string) (*uuid.UUID, error) {
	if r.findUserFn != nil {
		return r.findUserFn(address)
	}
	return nil, repository.ErrNotFound
}

func (r *mockObserverRepo) FindChannelByAddress(_ context.Context, address string) (*uuid.UUID, error) {
	if r.findChannelFn != nil {
		return r.findChannelFn(address)
	}
	return nil, repository.ErrNotFound
}

type mockPurchaseRepo struct {
	match *purchaserepo.TxHashMatch
	err   error
}

func (r *mockPurchaseRepo) FindByTxHash(_ context.Context, _ string) (*purchaserepo.TxHashMatch, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.match != nil {
		return r.match, nil
	}
	return nil, purchaserepo.ErrNotFound
}

type mockObserverMessenger struct {
	notifyCalls int32 // atomic for thread-safety
	events      []string
}

func (m *mockObserverMessenger) NotifyUserEvent(_ context.Context, req platgrpc.NotifyEventRequest) (*platgrpc.NotifyEventResult, error) {
	atomic.AddInt32(&m.notifyCalls, 1)
	m.events = append(m.events, req.EventType)
	return &platgrpc.NotifyEventResult{}, nil
}

func (m *mockObserverMessenger) GetUsersByIDs(_ context.Context, ids []string) (map[string]platgrpc.UserBrief, error) {
	out := make(map[string]platgrpc.UserBrief, len(ids))
	for _, id := range ids {
		out[id] = platgrpc.UserBrief{UserID: id, Username: "user", DisplayName: "User"}
	}
	return out, nil
}

func (m *mockObserverMessenger) notifyCallsTotal() int {
	return int(atomic.LoadInt32(&m.notifyCalls))
}

// ── Setup helpers ─────────────────────────────────────────────

var fixedBlockTime = time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)

func noopBlockTime(_ context.Context, _ uint64) (time.Time, error) {
	return fixedBlockTime, nil
}

func makeEvent(from, to common.Address, amount int64, txHash string) *sudatoken.SudaTokenTransfer {
	return &sudatoken.SudaTokenTransfer{
		From:  from,
		To:    to,
		Value: big.NewInt(amount),
		Raw: types.Log{
			TxHash:      common.HexToHash(txHash),
			Index:       0,
			BlockNumber: 100,
		},
	}
}

// ── Tests ─────────────────────────────────────────────────────

// TestHandleTransfer_Idempotency — если InsertTransaction returns inserted=false,
// handler не вызывает InsertDonation и не шлёт уведомления. DeletePending вызван.
func TestHandleTransfer_Idempotency(t *testing.T) {
	repo := &mockObserverRepo{insertResult: false} // already processed
	mc := &mockObserverMessenger{}
	h := NewSudaTokenHandlerForTest(repo, &mockPurchaseRepo{}, mc, noopBlockTime)

	from := common.HexToAddress("0x000000000000000000000000000000000000cafe")
	to := common.HexToAddress("0x000000000000000000000000000000000000beef")
	ev := makeEvent(from, to, 1000, "0xaaa")

	err := h.handleTransfer(context.Background(), ev)
	require.NoError(t, err)

	assert.Equal(t, 1, repo.insertTransactionCalls, "InsertTransaction must be called")
	assert.Equal(t, 0, repo.insertDonationCalls, "InsertDonation must NOT be called on duplicate")
	assert.Equal(t, 0, mc.notifyCallsTotal(), "NotifyUserEvent must NOT be called on duplicate")
	assert.Equal(t, 1, repo.deletePendingCalls, "DeletePending must be called for cleanup")
}

// TestHandleTransfer_RegularTransfer — обычный P2P transfer шлёт SUDA_RECEIVED
// + SUDA_SENT при известных адресах обоих участников.
func TestHandleTransfer_RegularTransfer(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	fromAddr := "0x000000000000000000000000000000000000000A"
	toAddr := "0x000000000000000000000000000000000000000B"

	repo := &mockObserverRepo{insertResult: true}
	repo.findUserFn = func(address string) (*uuid.UUID, error) {
		switch address {
		case common.HexToAddress(fromAddr).Hex():
			return &fromID, nil
		case common.HexToAddress(toAddr).Hex():
			return &toID, nil
		}
		return nil, repository.ErrNotFound
	}
	mc := &mockObserverMessenger{}
	h := NewSudaTokenHandlerForTest(repo, &mockPurchaseRepo{}, mc, noopBlockTime)

	ev := makeEvent(
		common.HexToAddress(fromAddr),
		common.HexToAddress(toAddr),
		1_000_000_000_000_000_000, // 1 SUDA
		"0xbbb",
	)
	require.NoError(t, h.handleTransfer(context.Background(), ev))

	assert.Equal(t, 0, repo.insertDonationCalls)
	assert.Equal(t, 1, repo.deletePendingCalls)
	assert.Equal(t, 2, mc.notifyCallsTotal(), "must send SUDA_RECEIVED + SUDA_SENT")
	assert.Contains(t, mc.events, EventSudaReceived)
	assert.Contains(t, mc.events, EventSudaSent)
}

// TestHandleTransfer_PurchaseDetected — если purchaseRepo возвращает match,
// тип транзакции = PURCHASE и шлётся PURCHASE_COMPLETED (без SUDA_RECEIVED).
func TestHandleTransfer_PurchaseDetected(t *testing.T) {
	toID := uuid.New()
	purchaseID := uuid.New()

	repo := &mockObserverRepo{insertResult: true}
	repo.findUserFn = func(address string) (*uuid.UUID, error) {
		return &toID, nil
	}
	purchRepo := &mockPurchaseRepo{
		match: &purchaserepo.TxHashMatch{
			ID:           purchaseID,
			UserID:       toID,
			PackageCode:  "MEDIUM",
			SudaAmountWei: big.NewInt(500),
		},
	}
	mc := &mockObserverMessenger{}
	h := NewSudaTokenHandlerForTest(repo, purchRepo, mc, noopBlockTime)

	ev := makeEvent(
		common.HexToAddress("0x0000000000000000000000000000000000000001"),
		common.HexToAddress("0x0000000000000000000000000000000000000002"),
		500, "0xccc",
	)
	require.NoError(t, h.handleTransfer(context.Background(), ev))

	assert.Equal(t, 0, repo.insertDonationCalls)
	assert.Equal(t, 1, mc.notifyCallsTotal(), "only PURCHASE_COMPLETED, no SUDA_RECEIVED")
	assert.Contains(t, mc.events, EventPurchaseCompleted)
	assert.NotContains(t, mc.events, EventSudaReceived)
}

// TestHandleTransfer_Donation_P2P — pending.ExpectedType=DONATION для P2P:
// InsertDonation вызван, DONATION_SENT + DONATION_RECEIVED оба шлются.
func TestHandleTransfer_Donation_P2P(t *testing.T) {
	fromID := uuid.New()
	toID := uuid.New()
	chatID := uuid.New()
	fromAddr := common.HexToAddress("0x000000000000000000000000000000000000000A")
	toAddr := common.HexToAddress("0x000000000000000000000000000000000000000B")

	repo := &mockObserverRepo{
		insertResult: true,
		pendingMeta: &repository.PendingMeta{
			ChatID:          chatID,
			ExpectedType:    "DONATION",
			DonationMessage: "thanks!",
		},
	}
	repo.findUserFn = func(address string) (*uuid.UUID, error) {
		if address == fromAddr.Hex() {
			return &fromID, nil
		}
		if address == toAddr.Hex() {
			return &toID, nil
		}
		return nil, repository.ErrNotFound
	}
	mc := &mockObserverMessenger{}
	h := NewSudaTokenHandlerForTest(repo, &mockPurchaseRepo{}, mc, noopBlockTime)

	ev := makeEvent(fromAddr, toAddr, 1000, "0xddd")
	require.NoError(t, h.handleTransfer(context.Background(), ev))

	assert.Equal(t, 1, repo.insertDonationCalls, "InsertDonation must be called")
	assert.Equal(t, 2, mc.notifyCallsTotal(), "DONATION_SENT + DONATION_RECEIVED")
	assert.Contains(t, mc.events, donation.EventDonationSent)
	assert.Contains(t, mc.events, donation.EventDonationReceived)
	assert.NotContains(t, mc.events, EventSudaReceived)

	var insertedParams repository.DonationInsertParams
	_ = insertedParams // verified by call count above
}

// TestHandleTransfer_Donation_Channel — донат каналу: InsertDonation вызван,
// шлётся только DONATION_SENT (DONATION_RECEIVED нет, т.к. получатель — не юзер).
func TestHandleTransfer_Donation_Channel(t *testing.T) {
	fromID := uuid.New()
	channelID := uuid.New()
	chatID := uuid.New()
	fromAddr := common.HexToAddress("0x000000000000000000000000000000000000000A")
	toAddr := common.HexToAddress("0x000000000000000000000000000000000000000C")

	repo := &mockObserverRepo{
		insertResult: true,
		pendingMeta: &repository.PendingMeta{
			ChatID:          chatID,
			ExpectedType:    "DONATION",
			DonationMessage: "keep going",
		},
	}
	repo.findUserFn = func(address string) (*uuid.UUID, error) {
		if address == fromAddr.Hex() {
			return &fromID, nil
		}
		return nil, repository.ErrNotFound // toAddr is a channel wallet, not a user
	}
	repo.findChannelFn = func(address string) (*uuid.UUID, error) {
		if address == toAddr.Hex() {
			return &channelID, nil
		}
		return nil, repository.ErrNotFound
	}

	mc := &mockObserverMessenger{}
	h := NewSudaTokenHandlerForTest(repo, &mockPurchaseRepo{}, mc, noopBlockTime)

	ev := makeEvent(fromAddr, toAddr, 5000, "0xeee")
	require.NoError(t, h.handleTransfer(context.Background(), ev))

	assert.Equal(t, 1, repo.insertDonationCalls, "InsertDonation must be called")
	assert.Equal(t, 1, mc.notifyCallsTotal(), "only DONATION_SENT, no DONATION_RECEIVED for channel")
	assert.Contains(t, mc.events, donation.EventDonationSent)
	assert.NotContains(t, mc.events, donation.EventDonationReceived)
}

// TestHandleTransfer_BlockTimeFnError — если blockTimeFn возвращает ошибку,
// handleTransfer пробрасывает её наружу.
func TestHandleTransfer_BlockTimeFnError(t *testing.T) {
	repo := &mockObserverRepo{insertResult: true}
	mc := &mockObserverMessenger{}
	h := NewSudaTokenHandlerForTest(repo, &mockPurchaseRepo{}, mc, func(_ context.Context, _ uint64) (time.Time, error) {
		return time.Time{}, errors.New("rpc error")
	})

	ev := makeEvent(
		common.HexToAddress("0xA"),
		common.HexToAddress("0xB"),
		100, "0xfff",
	)
	err := h.handleTransfer(context.Background(), ev)
	assert.Error(t, err)
	assert.Equal(t, 0, repo.insertTransactionCalls, "should fail before InsertTransaction")
}
