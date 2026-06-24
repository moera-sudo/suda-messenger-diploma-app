package service_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/gating"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/gating/service"
	walletrepo "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/mocks"
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"
)

// ── Setup ──────────────────────────────────────────────────────

type fakeMessenger struct {
	CheckPermFn func(ctx context.Context, userID, channelID, permission string) (*platgrpc.ChannelPermissionResult, error)
}

func (m *fakeMessenger) CheckChannelPermission(ctx context.Context, userID, channelID, permission string) (*platgrpc.ChannelPermissionResult, error) {
	if m.CheckPermFn != nil {
		return m.CheckPermFn(ctx, userID, channelID, permission)
	}
	return &platgrpc.ChannelPermissionResult{Granted: true, Reason: "ok"}, nil
}

type testEnv struct {
	svc       *service.Service
	wallet    *mocks.MockWalletRepo
	messenger *fakeMessenger
}

func setup(_ *testing.T) *testEnv {
	wallet := mocks.NewMockWalletRepo()
	mc := &fakeMessenger{}
	svc := service.NewWithDeps(wallet, mc)
	return &testEnv{svc: svc, wallet: wallet, messenger: mc}
}

// ── CreateRule ─────────────────────────────────────────────────

func TestCreateRule_Success(t *testing.T) {
	e := setup(t)
	owner := uuid.New()
	chatID := uuid.New()

	err := e.svc.CreateRule(context.Background(), service.RuleInput{
		OwnerID:           owner,
		ChatID:            chatID,
		MinSudaBalanceWei: "100000000000000000000", // 100 SUDA
	})
	require.NoError(t, err)
	assert.Equal(t, 1, e.wallet.UpsertGatingRuleCalls)
}

func TestCreateRule_NotOwner_Forbidden(t *testing.T) {
	e := setup(t)
	e.messenger.CheckPermFn = func(_ context.Context, _, _, _ string) (*platgrpc.ChannelPermissionResult, error) {
		return &platgrpc.ChannelPermissionResult{Granted: false, Reason: "not_owner"}, nil
	}

	err := e.svc.CreateRule(context.Background(), service.RuleInput{
		OwnerID:           uuid.New(),
		ChatID:            uuid.New(),
		MinSudaBalanceWei: "100",
	})
	assert.ErrorIs(t, err, gating.ErrForbidden)
	assert.Equal(t, 0, e.wallet.UpsertGatingRuleCalls)
}

func TestCreateRule_NegativeBalance_InvalidInput(t *testing.T) {
	e := setup(t)
	err := e.svc.CreateRule(context.Background(), service.RuleInput{
		OwnerID:           uuid.New(),
		ChatID:            uuid.New(),
		MinSudaBalanceWei: "-1",
	})
	assert.ErrorIs(t, err, gating.ErrInvalidInput)
	assert.Equal(t, 0, e.wallet.UpsertGatingRuleCalls)
}

func TestCreateRule_NonNumeric_InvalidInput(t *testing.T) {
	e := setup(t)
	err := e.svc.CreateRule(context.Background(), service.RuleInput{
		OwnerID:           uuid.New(),
		ChatID:            uuid.New(),
		MinSudaBalanceWei: "abc",
	})
	assert.ErrorIs(t, err, gating.ErrInvalidInput)
}

// ── GetRule ────────────────────────────────────────────────────

func TestGetRule_Exists(t *testing.T) {
	e := setup(t)
	chatID := uuid.New()
	minBal, _ := new(big.Int).SetString("250000000000000000000", 10)
	e.wallet.GetGatingRuleFn = func(_ context.Context, cid uuid.UUID) (*walletrepo.GatingRule, error) {
		assert.Equal(t, chatID, cid)
		return &walletrepo.GatingRule{ChatID: chatID, MinSudaBalance: minBal}, nil
	}

	rule, err := e.svc.GetRule(context.Background(), chatID)
	require.NoError(t, err)
	assert.True(t, rule.HasRule)
	assert.Equal(t, "250000000000000000000", rule.MinSudaBalanceWei)
}

func TestGetRule_NotExists_OpenChannel(t *testing.T) {
	e := setup(t)
	chatID := uuid.New()
	e.wallet.GetGatingRuleFn = func(_ context.Context, _ uuid.UUID) (*walletrepo.GatingRule, error) {
		return nil, walletrepo.ErrNotFound
	}

	rule, err := e.svc.GetRule(context.Background(), chatID)
	require.NoError(t, err)
	assert.False(t, rule.HasRule)
	assert.Equal(t, "0", rule.MinSudaBalanceWei)
}

// ── DeleteRule ─────────────────────────────────────────────────

func TestDeleteRule_Success(t *testing.T) {
	e := setup(t)
	err := e.svc.DeleteRule(context.Background(), uuid.New(), uuid.New())
	require.NoError(t, err)
	assert.Equal(t, 1, e.wallet.DeleteGatingRuleCalls)
}

func TestDeleteRule_NotOwner_Forbidden(t *testing.T) {
	e := setup(t)
	e.messenger.CheckPermFn = func(_ context.Context, _, _, _ string) (*platgrpc.ChannelPermissionResult, error) {
		return &platgrpc.ChannelPermissionResult{Granted: false, Reason: "not_owner"}, nil
	}

	err := e.svc.DeleteRule(context.Background(), uuid.New(), uuid.New())
	assert.ErrorIs(t, err, gating.ErrForbidden)
	assert.Equal(t, 0, e.wallet.DeleteGatingRuleCalls)
}
