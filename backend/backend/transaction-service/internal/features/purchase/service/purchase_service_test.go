package service_test

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase"
	purchaserepo "github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase/service"
	walletmodels "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet"
	walletrepo "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/mocks"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/blockchain"
)

// ── Setup ──────────────────────────────────────────────────────

type testEnv struct {
	svc       *service.Service
	repo      *mocks.MockPurchaseRepo
	wallet    *mocks.MockWalletRepo
	sender    *mocks.MockTokenSender
	treasury  *mocks.MockTreasurySigner
	redis     *redis.Client
	miniredis *miniredis.Miniredis
}

func setup(t *testing.T) *testEnv {
	mr := miniredis.RunT(t)
	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	repo := mocks.NewMockPurchaseRepo()
	wallet := mocks.NewMockWalletRepo()
	sender := mocks.NewMockTokenSender()
	treasury := mocks.NewMockTreasurySigner()

	// simulateProcessing → no-op для скорости тестов
	noopSim := func(_ context.Context) error { return nil }

	svc := service.NewWithDeps(repo, wallet, rc, sender, treasury, noopSim)

	return &testEnv{
		svc: svc, repo: repo, wallet: wallet,
		sender: sender, treasury: treasury,
		redis: rc, miniredis: mr,
	}
}

// helper: положить готовый purchase в mock repo
func (e *testEnv) givenPurchase(p *purchase.Purchase) {
	e.repo.GetByIDFn = func(_ context.Context, id uuid.UUID) (*purchase.Purchase, error) {
		if id != p.ID {
			return nil, purchaserepo.ErrNotFound
		}
		return p, nil
	}
}

// ── Initiate ───────────────────────────────────────────────────

func TestInitiate_Success(t *testing.T) {
	e := setup(t)
	userID := uuid.New()
	createdID := uuid.New()
	e.repo.CreateFn = func(_ context.Context, p purchaserepo.CreateParams) (uuid.UUID, error) {
		assert.Equal(t, userID, p.UserID)
		assert.Equal(t, "MEDIUM", p.PackageCode)
		assert.Equal(t, "CARD", p.PaymentMethod)
		return createdID, nil
	}

	res, err := e.svc.Initiate(context.Background(), service.InitiateInput{
		UserID:        userID,
		PackageCode:   "MEDIUM",
		PaymentMethod: "CARD",
	})

	require.NoError(t, err)
	assert.Equal(t, createdID, res.PurchaseID)
	assert.Equal(t, "MEDIUM", res.PackageCode)
	assert.Equal(t, purchase.StatusPending, res.Status)
	assert.Equal(t, "8.99", res.FiatAmount)
	assert.Equal(t, 1, e.repo.CreateCalls)
}

func TestInitiate_DefaultPaymentMethod(t *testing.T) {
	e := setup(t)
	e.svc.Initiate(context.Background(), service.InitiateInput{
		UserID:      uuid.New(),
		PackageCode: "SMALL",
	})
	// Проверяем что Create вызвался — payment_method обрабатывается внутри (CARD по умолчанию).
	assert.Equal(t, 1, e.repo.CreateCalls)
}

func TestInitiate_UnknownPackage(t *testing.T) {
	e := setup(t)
	res, err := e.svc.Initiate(context.Background(), service.InitiateInput{
		UserID:      uuid.New(),
		PackageCode: "UNKNOWN_PACKAGE",
	})
	assert.ErrorIs(t, err, service.ErrPackageUnknown)
	assert.Nil(t, res)
	assert.Equal(t, 0, e.repo.CreateCalls)
}

func TestInitiate_InvalidPaymentMethod(t *testing.T) {
	e := setup(t)
	res, err := e.svc.Initiate(context.Background(), service.InitiateInput{
		UserID:        uuid.New(),
		PackageCode:   "SMALL",
		PaymentMethod: "BITCOIN", // не из CARD/APPLE_PAY/GOOGLE_PAY
	})
	assert.ErrorIs(t, err, service.ErrInvalidInput)
	assert.Nil(t, res)
	assert.Equal(t, 0, e.repo.CreateCalls)
}

// ── Confirm ────────────────────────────────────────────────────

func newPendingPurchase(userID uuid.UUID) *purchase.Purchase {
	amount := new(big.Int).Mul(big.NewInt(500), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))
	return &purchase.Purchase{
		ID:            uuid.New(),
		UserID:        userID,
		PackageCode:   "MEDIUM",
		SudaAmountWei: amount,
		FiatAmount:    "8.99",
		FiatCurrency:  "USD",
		Status:        purchase.StatusPending,
		PaymentMethod: "CARD",
		CreatedAt:     time.Now(),
	}
}

func TestConfirm_Success(t *testing.T) {
	e := setup(t)
	userID := uuid.New()
	p := newPendingPurchase(userID)
	e.givenPurchase(p)
	e.wallet.GetUserWalletFn = func(_ context.Context, _ uuid.UUID) (*walletmodels.Wallet, error) {
		return &walletmodels.Wallet{
			UserID:  userID,
			Address: "0x000000000000000000000000000000000000beef",
		}, nil
	}

	res, err := e.svc.Confirm(context.Background(), service.ConfirmInput{
		PurchaseID: p.ID,
		UserID:     userID,
	})

	require.NoError(t, err)
	assert.Equal(t, purchase.StatusCompleted, res.Status)
	assert.NotEmpty(t, res.TxHash)

	assert.Equal(t, 1, e.repo.MarkProcessingCalls, "MarkProcessing должен вызваться 1 раз")
	assert.Equal(t, 1, e.repo.MarkCompletedCalls, "MarkCompleted должен вызваться")
	assert.Equal(t, 0, e.repo.MarkFailedCalls, "MarkFailed не должен вызваться")
	assert.Equal(t, 1, e.sender.SendCalls, "SendTokenTransfer должен вызваться")
	assert.Equal(t, 1, e.wallet.WriteAuditCalls)
	assert.Equal(t, 1, e.wallet.InsertPendingCalls)
}

func TestConfirm_RateLimited(t *testing.T) {
	e := setup(t)
	userID := uuid.New()
	p := newPendingPurchase(userID)
	e.givenPurchase(p)
	e.wallet.GetUserWalletFn = func(_ context.Context, _ uuid.UUID) (*walletmodels.Wallet, error) {
		return &walletmodels.Wallet{UserID: userID, Address: "0xbeef"}, nil
	}

	// Первый confirm — успех
	_, err := e.svc.Confirm(context.Background(), service.ConfirmInput{PurchaseID: p.ID, UserID: userID})
	require.NoError(t, err)

	// Второй confirm в окне 5s — должен попасть в rate-limit
	_, err = e.svc.Confirm(context.Background(), service.ConfirmInput{PurchaseID: p.ID, UserID: userID})
	assert.ErrorIs(t, err, service.ErrRateLimited)
	// MarkProcessing не должен вызваться второй раз
	assert.Equal(t, 1, e.repo.MarkProcessingCalls)
}

func TestConfirm_NotFound(t *testing.T) {
	e := setup(t)
	e.repo.GetByIDFn = func(_ context.Context, _ uuid.UUID) (*purchase.Purchase, error) {
		return nil, purchaserepo.ErrNotFound
	}

	res, err := e.svc.Confirm(context.Background(), service.ConfirmInput{
		PurchaseID: uuid.New(),
		UserID:     uuid.New(),
	})
	assert.ErrorIs(t, err, service.ErrNotFound)
	assert.Nil(t, res)
	assert.Equal(t, 0, e.repo.MarkProcessingCalls)
}

func TestConfirm_NotOwner(t *testing.T) {
	e := setup(t)
	ownerID := uuid.New()
	otherID := uuid.New()
	p := newPendingPurchase(ownerID)
	e.givenPurchase(p)

	res, err := e.svc.Confirm(context.Background(), service.ConfirmInput{
		PurchaseID: p.ID,
		UserID:     otherID, // не владелец
	})
	assert.ErrorIs(t, err, service.ErrForbidden)
	assert.Nil(t, res)
	assert.Equal(t, 0, e.repo.MarkProcessingCalls)
}

func TestConfirm_AlreadyHandled(t *testing.T) {
	e := setup(t)
	userID := uuid.New()
	p := newPendingPurchase(userID)
	p.Status = purchase.StatusCompleted // уже завершена
	e.givenPurchase(p)

	res, err := e.svc.Confirm(context.Background(), service.ConfirmInput{
		PurchaseID: p.ID,
		UserID:     userID,
	})
	assert.ErrorIs(t, err, service.ErrAlreadyHandled)
	assert.Nil(t, res)
	assert.Equal(t, 0, e.repo.MarkProcessingCalls)
}

func TestConfirm_MarkProcessingCAS_AlreadyHandled(t *testing.T) {
	// CAS-гонка: GetByID видит PENDING, но MarkProcessing возвращает ErrNotFound
	// (status уже сменился) → ErrAlreadyHandled.
	e := setup(t)
	userID := uuid.New()
	p := newPendingPurchase(userID)
	e.givenPurchase(p)
	e.repo.MarkProcessingFn = func(_ context.Context, _ uuid.UUID) error {
		return purchaserepo.ErrNotFound
	}

	_, err := e.svc.Confirm(context.Background(), service.ConfirmInput{
		PurchaseID: p.ID,
		UserID:     userID,
	})
	assert.ErrorIs(t, err, service.ErrAlreadyHandled)
}

func TestConfirm_NoWallet(t *testing.T) {
	e := setup(t)
	userID := uuid.New()
	p := newPendingPurchase(userID)
	e.givenPurchase(p)
	e.wallet.GetUserWalletFn = func(_ context.Context, _ uuid.UUID) (*walletmodels.Wallet, error) {
		return nil, walletrepo.ErrNotFound
	}

	res, err := e.svc.Confirm(context.Background(), service.ConfirmInput{
		PurchaseID: p.ID,
		UserID:     userID,
	})
	assert.ErrorIs(t, err, service.ErrWalletNotFound)
	assert.Nil(t, res)
	// MarkFailed должен быть вызван с reason=no_wallet
	assert.Equal(t, 1, e.repo.MarkFailedCalls)
	assert.Equal(t, 0, e.sender.SendCalls, "broadcast не должен запуститься")
}

func TestConfirm_BroadcastFailed(t *testing.T) {
	e := setup(t)
	userID := uuid.New()
	p := newPendingPurchase(userID)
	e.givenPurchase(p)
	e.wallet.GetUserWalletFn = func(_ context.Context, _ uuid.UUID) (*walletmodels.Wallet, error) {
		return &walletmodels.Wallet{UserID: userID, Address: "0xbeef"}, nil
	}
	e.sender.SendFn = func(_ context.Context, _ blockchain.Signer, _ common.Address, _ *big.Int) (string, error) {
		return "", errors.New("rpc down")
	}

	res, err := e.svc.Confirm(context.Background(), service.ConfirmInput{
		PurchaseID: p.ID,
		UserID:     userID,
	})
	assert.Error(t, err)
	assert.Nil(t, res)
	// MarkFailed должен быть вызван с reason=broadcast_failed:...
	assert.Equal(t, 1, e.repo.MarkFailedCalls)
	assert.Equal(t, 0, e.repo.MarkCompletedCalls)
}

// ── ListPackages / GetHistory ─────────────────────────────────

func TestListPackages_HasFour(t *testing.T) {
	e := setup(t)
	pkgs := e.svc.ListPackages(context.Background())
	assert.Len(t, pkgs, 4)
	codes := make(map[string]bool)
	for _, p := range pkgs {
		codes[p.Code] = true
	}
	assert.True(t, codes["SMALL"])
	assert.True(t, codes["MEDIUM"])
	assert.True(t, codes["LARGE"])
	assert.True(t, codes["MEGA"])
}

func TestGetHistory_OK(t *testing.T) {
	e := setup(t)
	userID := uuid.New()
	e.repo.ListForUserFn = func(_ context.Context, uid uuid.UUID, _ int, _ int) ([]purchase.Purchase, error) {
		assert.Equal(t, userID, uid)
		return []purchase.Purchase{*newPendingPurchase(userID)}, nil
	}

	items, err := e.svc.GetHistory(context.Background(), userID, 50, 0)
	require.NoError(t, err)
	assert.Len(t, items, 1)
}
