package service_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/donation/service"
	walletmodels "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet"
	walletrepo "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/repository"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/mocks"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/platform/blockchain"
	platgrpc "github.com/moera-sudo/backend/backend/transaction-service/internal/platform/grpc"
)

// ── Setup ──────────────────────────────────────────────────────

type testEnv struct {
	svc       *service.Service
	wallet    *mocks.MockWalletRepo
	sender    *mocks.MockTokenSender
	reader    *fakeReader
	encryptor *fakeEncryptor
	messenger *fakeMessenger
}

// fakeReader, fakeEncryptor, fakeMessenger — здесь inline для краткости.
// Они simple enough что выносить в internal/mocks/ не нужно.

type fakeReader struct {
	BalanceFn func(ctx context.Context, addr common.Address) (*big.Int, error)
}

func (r *fakeReader) SudaBalanceOf(ctx context.Context, addr common.Address) (*big.Int, error) {
	if r.BalanceFn != nil {
		return r.BalanceFn(ctx, addr)
	}
	// По умолчанию — огромный баланс, достаточный для любых тестов.
	return new(big.Int).Mul(big.NewInt(1_000_000), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), nil
}

type fakeEncryptor struct{}

// Decrypt — возвращает валидный hex-приватник (для NewUserSigner).
func (e *fakeEncryptor) Decrypt(_ uint8, _ string) ([]byte, error) {
	// Любой валидный 32-байтный hex приватник.
	return []byte("4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"), nil
}

type fakeMessenger struct {
	ResolveFn          func(ctx context.Context, username string) (*platgrpc.ResolvedUser, error)
	CheckMembershipFn  func(ctx context.Context, chatID string, userIDs []string) (*platgrpc.ChatMembershipResult, error)
	ResolveCalls       int
	CheckMembersCalls  int
}

func (m *fakeMessenger) ResolveUsername(ctx context.Context, username string) (*platgrpc.ResolvedUser, error) {
	m.ResolveCalls++
	if m.ResolveFn != nil {
		return m.ResolveFn(ctx, username)
	}
	return &platgrpc.ResolvedUser{Found: false}, nil
}

func (m *fakeMessenger) CheckChatMembership(ctx context.Context, chatID string, userIDs []string) (*platgrpc.ChatMembershipResult, error) {
	m.CheckMembersCalls++
	if m.CheckMembershipFn != nil {
		return m.CheckMembershipFn(ctx, chatID, userIDs)
	}
	return nil, nil
}

func setup(t *testing.T) *testEnv {
	wallet := mocks.NewMockWalletRepo()
	sender := mocks.NewMockTokenSender()
	reader := &fakeReader{}
	enc := &fakeEncryptor{}
	mc := &fakeMessenger{}

	svc := service.NewWithDeps(wallet, enc, reader, sender, mc)

	return &testEnv{
		svc: svc, wallet: wallet, sender: sender,
		reader: reader, encryptor: enc, messenger: mc,
	}
}

// helper: senderWallet с фиксированным адресом
func (e *testEnv) givenSenderWallet(userID uuid.UUID, address string) {
	e.wallet.GetUserWalletFn = func(_ context.Context, uid uuid.UUID) (*walletmodels.Wallet, error) {
		if uid != userID {
			return nil, walletrepo.ErrNotFound
		}
		return &walletmodels.Wallet{
			UserID:              userID,
			Address:             address,
			EncryptedPrivateKey: "fake",
			KeyVersion:          1,
		}, nil
	}
}

// ── Validation ────────────────────────────────────────────────

func TestDonate_NoRecipient_InvalidInput(t *testing.T) {
	e := setup(t)
	_, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID: uuid.New(),
		AmountWei:  "1000",
	})
	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestDonate_BothRecipients_InvalidInput(t *testing.T) {
	e := setup(t)
	_, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID:  uuid.New(),
		ToUsername:  "alice",
		ToChannelID: uuid.New(),
		AmountWei:   "1000",
	})
	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

func TestDonate_NegativeAmount_InvalidInput(t *testing.T) {
	e := setup(t)
	_, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID: uuid.New(),
		ToUsername: "alice",
		AmountWei:  "-5",
	})
	assert.ErrorIs(t, err, service.ErrInvalidInput)
}

// ── P2P ───────────────────────────────────────────────────────

func TestDonate_P2P_RecipientNotFound(t *testing.T) {
	e := setup(t)
	e.messenger.ResolveFn = func(_ context.Context, _ string) (*platgrpc.ResolvedUser, error) {
		return &platgrpc.ResolvedUser{Found: false}, nil
	}
	_, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID: uuid.New(),
		ToUsername: "ghost",
		AmountWei:  "1000000000000000000",
	})
	assert.ErrorIs(t, err, service.ErrRecipientNotFound)
}

func TestDonate_P2P_RecipientHasNoWallet(t *testing.T) {
	e := setup(t)
	e.messenger.ResolveFn = func(_ context.Context, _ string) (*platgrpc.ResolvedUser, error) {
		return &platgrpc.ResolvedUser{Found: true, UserID: uuid.New().String(), WalletAddress: ""}, nil
	}
	_, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID: uuid.New(),
		ToUsername: "newbie",
		AmountWei:  "1000000000000000000",
	})
	assert.ErrorIs(t, err, service.ErrRecipientHasNoWallet)
}

func TestDonate_P2P_Success(t *testing.T) {
	e := setup(t)
	from := uuid.New()
	to := uuid.New()
	e.messenger.ResolveFn = func(_ context.Context, _ string) (*platgrpc.ResolvedUser, error) {
		return &platgrpc.ResolvedUser{Found: true, UserID: to.String(), WalletAddress: "0x000000000000000000000000000000000000bEEf"}, nil
	}
	e.givenSenderWallet(from, "0x000000000000000000000000000000000000CAFE")

	res, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID: from,
		ToUsername: "alice",
		AmountWei:  "1000000000000000000",
		Message:    "thanks",
	})

	require.NoError(t, err)
	assert.False(t, res.IsChannel)
	assert.Equal(t, to.String(), res.ToUserID)
	assert.NotEmpty(t, res.TxHash)
	assert.Equal(t, 1, e.sender.SendCalls)
	assert.Equal(t, 1, e.wallet.WriteAuditCalls)
	assert.Equal(t, 1, e.wallet.InsertPendingDonationCalls)
}

func TestDonate_SelfDonation(t *testing.T) {
	e := setup(t)
	from := uuid.New()
	myAddr := "0x000000000000000000000000000000000000CAFE"
	e.messenger.ResolveFn = func(_ context.Context, _ string) (*platgrpc.ResolvedUser, error) {
		return &platgrpc.ResolvedUser{Found: true, UserID: uuid.New().String(), WalletAddress: myAddr}, nil
	}
	e.givenSenderWallet(from, myAddr)

	_, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID: from,
		ToUsername: "myself",
		AmountWei:  "1000",
	})
	assert.ErrorIs(t, err, service.ErrSelfDonation)
	assert.Equal(t, 0, e.sender.SendCalls)
}

func TestDonate_InsufficientBalance(t *testing.T) {
	e := setup(t)
	from := uuid.New()
	e.messenger.ResolveFn = func(_ context.Context, _ string) (*platgrpc.ResolvedUser, error) {
		return &platgrpc.ResolvedUser{Found: true, UserID: uuid.New().String(), WalletAddress: "0xbeef"}, nil
	}
	e.givenSenderWallet(from, "0xcafe")
	e.reader.BalanceFn = func(_ context.Context, _ common.Address) (*big.Int, error) {
		return big.NewInt(100), nil // меньше чем 1000
	}

	_, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID: from,
		ToUsername: "alice",
		AmountWei:  "1000",
	})
	assert.ErrorIs(t, err, service.ErrInsufficientBalance)
	assert.Equal(t, 0, e.sender.SendCalls)
}

func TestDonate_SenderWalletNotFound(t *testing.T) {
	e := setup(t)
	from := uuid.New()
	e.messenger.ResolveFn = func(_ context.Context, _ string) (*platgrpc.ResolvedUser, error) {
		return &platgrpc.ResolvedUser{Found: true, UserID: uuid.New().String(), WalletAddress: "0xbeef"}, nil
	}
	e.wallet.GetUserWalletFn = func(_ context.Context, _ uuid.UUID) (*walletmodels.Wallet, error) {
		return nil, walletrepo.ErrNotFound
	}

	_, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID: from,
		ToUsername: "alice",
		AmountWei:  "1000",
	})
	assert.ErrorIs(t, err, service.ErrWalletNotFound)
}

// ── P2P + chat membership ──────────────────────────────────────

func TestDonate_P2P_SenderNotInChat(t *testing.T) {
	e := setup(t)
	from := uuid.New()
	to := uuid.New()
	chatID := uuid.New()
	e.messenger.ResolveFn = func(_ context.Context, _ string) (*platgrpc.ResolvedUser, error) {
		return &platgrpc.ResolvedUser{Found: true, UserID: to.String(), WalletAddress: "0xbeef"}, nil
	}
	e.givenSenderWallet(from, "0xcafe")
	e.messenger.CheckMembershipFn = func(_ context.Context, _ string, _ []string) (*platgrpc.ChatMembershipResult, error) {
		return &platgrpc.ChatMembershipResult{
			ChatExists: true,
			Members: []platgrpc.UserMembership{
				{UserID: from.String(), IsMember: false},
				{UserID: to.String(), IsMember: true},
			},
		}, nil
	}

	_, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID: from,
		ToUsername: "alice",
		AmountWei:  "1000",
		ChatID:     chatID,
	})
	assert.ErrorIs(t, err, service.ErrSenderNotInChat)
	assert.Equal(t, 0, e.sender.SendCalls)
}

func TestDonate_P2P_RecipientNotInChat(t *testing.T) {
	e := setup(t)
	from := uuid.New()
	to := uuid.New()
	chatID := uuid.New()
	e.messenger.ResolveFn = func(_ context.Context, _ string) (*platgrpc.ResolvedUser, error) {
		return &platgrpc.ResolvedUser{Found: true, UserID: to.String(), WalletAddress: "0xbeef"}, nil
	}
	e.givenSenderWallet(from, "0xcafe")
	e.messenger.CheckMembershipFn = func(_ context.Context, _ string, _ []string) (*platgrpc.ChatMembershipResult, error) {
		return &platgrpc.ChatMembershipResult{
			ChatExists: true,
			Members: []platgrpc.UserMembership{
				{UserID: from.String(), IsMember: true},
				{UserID: to.String(), IsMember: false},
			},
		}, nil
	}

	_, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID: from,
		ToUsername: "alice",
		AmountWei:  "1000",
		ChatID:     chatID,
	})
	assert.ErrorIs(t, err, service.ErrRecipientNotInChat)
}

// ── Channel donations ──────────────────────────────────────────

func TestDonate_Channel_Success(t *testing.T) {
	e := setup(t)
	from := uuid.New()
	channelID := uuid.New()
	e.wallet.GetChannelWalletFn = func(_ context.Context, cid uuid.UUID) (*walletmodels.ChannelWallet, error) {
		if cid != channelID {
			return nil, walletrepo.ErrNotFound
		}
		return &walletmodels.ChannelWallet{
			ChannelID: channelID,
			Address:   "0x000000000000000000000000000000000000Bee5",
		}, nil
	}
	e.givenSenderWallet(from, "0x000000000000000000000000000000000000CAFE")

	res, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID:  from,
		ToChannelID: channelID,
		AmountWei:   "5000000000000000000",
		Message:     "keep going",
	})

	require.NoError(t, err)
	assert.True(t, res.IsChannel)
	assert.Equal(t, channelID.String(), res.ToChannelID)
	assert.Equal(t, 1, e.sender.SendCalls)
	assert.Equal(t, 1, e.wallet.InsertPendingDonationCalls)
}

func TestDonate_ChannelWalletNotFound(t *testing.T) {
	e := setup(t)
	e.wallet.GetChannelWalletFn = func(_ context.Context, _ uuid.UUID) (*walletmodels.ChannelWallet, error) {
		return nil, walletrepo.ErrNotFound
	}

	_, err := e.svc.Donate(context.Background(), service.DonationInput{
		FromUserID:  uuid.New(),
		ToChannelID: uuid.New(),
		AmountWei:   "1000",
	})
	assert.ErrorIs(t, err, service.ErrChannelWalletNotFound)
}

// Дополнительная проверка: signer/sender используется с blockchain.Signer интерфейсом.
var _ blockchain.Signer = (*mocks.MockTreasurySigner)(nil)
