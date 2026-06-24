package service_test

// import (
// 	"context"
// 	"fmt"
// 	"testing"
// 	"time"

// 	"github.com/alicebob/miniredis/v2"
// 	"github.com/google/uuid"
// 	"github.com/redis/go-redis/v9"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/user"
// 	userService "github.com/moera-sudo/backend/backend/messenger-service/internal/features/user/service"
// 	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
// )

// var (
// 	userA = uuid.MustParse("11111111-1111-1111-1111-111111111111")
// 	userB = uuid.MustParse("22222222-2222-2222-2222-222222222222")
// )

// // ── Mock User Repository ────────────────────────────────────

// type MockUserRepo struct {
// 	GetUserByIDFn          func(ctx context.Context, id uuid.UUID) (*user.UserResponse, error)
// 	UpdateLastSeenFn       func(ctx context.Context, id uuid.UUID) error
// 	GetAllInterlocutorsFn  func(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error)
// 	GetLastSeenAtFn        func(ctx context.Context, id uuid.UUID) (*time.Time, error)
// 	UpdateLastSeenCalls    int
// 	GetUserByExactUsername func()
// 	GetWalletAddressFn func()
// 	SetWalletAddressFn func()
	
// }

// func (m *MockUserRepo) UpdateProfile(_ context.Context, _ uuid.UUID, _, _, _, _ string) error { return nil }
// func (m *MockUserRepo) UpdateAvatar(_ context.Context, _ uuid.UUID, _ uuid.UUID) error { return nil }
// func (m *MockUserRepo) UserLogout(_ context.Context, _ uuid.UUID, _ string) error { return nil }
// func (m *MockUserRepo) RemoveDeviceToken(_ context.Context, _ uuid.UUID, _ string) error { return nil }

// func (m *MockUserRepo) GetUserByID(ctx context.Context, id uuid.UUID) (*user.UserResponse, error) {
// 	if m.GetUserByIDFn != nil { return m.GetUserByIDFn(ctx, id) }
// 	return &user.UserResponse{ID: id, Username: "@test", DisplayName: "Test"}, nil
// }

// func (m *MockUserRepo) UpdateLastSeen(ctx context.Context, id uuid.UUID) error {
// 	m.UpdateLastSeenCalls++
// 	if m.UpdateLastSeenFn != nil { return m.UpdateLastSeenFn(ctx, id) }
// 	return nil
// }

// func (m *MockUserRepo) GetAllInterlocutorIDs(ctx context.Context, id uuid.UUID) ([]uuid.UUID, error) {
// 	if m.GetAllInterlocutorsFn != nil { return m.GetAllInterlocutorsFn(ctx, id) }
// 	return nil, nil
// }

// func (m *MockUserRepo) GetLastSeenAt(ctx context.Context, id uuid.UUID) (*time.Time, error) {
// 	if m.GetLastSeenAtFn != nil { return m.GetLastSeenAtFn(ctx, id) }
// 	t := time.Now().Add(-5 * time.Minute)
// 	return &t, nil
// }

// // ── Mock Moderation Repo ────────────────────────────────────

// type MockModRepo struct {
// 	IsBlockedFn func(ctx context.Context, a, b uuid.UUID) (bool, error)
// }

// func (m *MockModRepo) MuteChat(_ context.Context, _, _ uuid.UUID, _ *time.Time) error { return nil }
// func (m *MockModRepo) UnmuteChat(_ context.Context, _, _ uuid.UUID) error { return nil }
// func (m *MockModRepo) IsChatMuted(_ context.Context, _, _ uuid.UUID) (bool, error) {
// 	return false, nil
// }
// func (m *MockModRepo) BlockUser(_ context.Context, _, _ uuid.UUID) error { return nil }
// func (m *MockModRepo) UnblockUser(_ context.Context, _, _ uuid.UUID) error { return nil }
// func (m *MockModRepo) GetBlockedUsers(_ context.Context, _ uuid.UUID) ([]uuid.UUID, error) { return nil, nil }
// func (m *MockModRepo) PinMessage(_ context.Context, _ uuid.UUID, _ int64, _ uuid.UUID) error { return nil }
// func (m *MockModRepo) UnpinMessage(_ context.Context, _ uuid.UUID, _ int64) error { return nil }
// func (m *MockModRepo) GetPinnedMessages(_ context.Context, _ uuid.UUID) ([]chat.PinnedMessageInfoResponse, error) {
// 	return nil, nil
// }
// func (m *MockModRepo) SetContactName(_ context.Context, _, _ uuid.UUID, _ string) error { return nil }
// func (m *MockModRepo) RemoveContactName(_ context.Context, _, _ uuid.UUID) error { return nil }
// func (m *MockModRepo) GetContactName(_ context.Context, _, _ uuid.UUID) (string, error) { return "", nil }

// func (m *MockModRepo) GetContacts(_ context.Context, _ uuid.UUID) ([]chat.Contact, error) {
// 	return nil, nil
// }
// func (m *MockModRepo) IsBlocked(ctx context.Context, a, b uuid.UUID) (bool, error) {
// 	if m.IsBlockedFn != nil { return m.IsBlockedFn(ctx, a, b) }
// 	return false, nil
// }

// // ── Test Setup ──────────────────────────────────────────────

// type testEnv struct {
// 	svc      *userService.Service
// 	userRepo *MockUserRepo
// 	modRepo  *MockModRepo
// 	redis    *redis.Client
// }

// func setup(t *testing.T) *testEnv {
// 	mr := miniredis.RunT(t)
// 	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})

// 	userRepo := &MockUserRepo{}
// 	modRepo := &MockModRepo{}

// 	// svc := userService.NewService(userRepo, nil, rc, nil, nil, nil, modRepo, nil)

// 	return &testEnv{svc: svc, userRepo: userRepo, modRepo: modRepo, redis: rc}
// }

// // ═════════════════════════════════════════════════════════════
// //  OnUserDisconnected — UpdateLastSeen
// // ═════════════════════════════════════════════════════════════

// func TestOnUserDisconnected_UpdatesLastSeen(t *testing.T) {
// 	e := setup(t)

// 	e.svc.OnUserDisconnected(userA)

// 	assert.Equal(t, 1, e.userRepo.UpdateLastSeenCalls)
// }

// func TestOnUserDisconnected_LastSeenDBError(t *testing.T) {
// 	e := setup(t)
// 	e.userRepo.UpdateLastSeenFn = func(_ context.Context, _ uuid.UUID) error {
// 		return fmt.Errorf("connection refused")
// 	}

// 	// Не паникует
// 	e.svc.OnUserDisconnected(userA)

// 	assert.Equal(t, 1, e.userRepo.UpdateLastSeenCalls)
// }

// // ═════════════════════════════════════════════════════════════
// //  GetUserStatus
// // ═════════════════════════════════════════════════════════════

// func TestGetUserStatus_Online(t *testing.T) {
// 	e := setup(t)

// 	// Ставим юзера онлайн в Redis
// 	e.redis.Set(context.Background(), fmt.Sprintf("user:%s:online", userB), "1", time.Minute)

// 	status, err := e.svc.GetUserStatus(context.Background(), userA, userB)

// 	require.NoError(t, err)
// 	assert.True(t, status.IsOnline)
// 	assert.Equal(t, userB, status.UserID)
// 	assert.Nil(t, status.LastSeenAt) // Онлайн — не показываем last_seen
// }

// func TestGetUserStatus_Offline(t *testing.T) {
// 	e := setup(t)

// 	lastSeen := time.Now().Add(-30 * time.Minute)
// 	e.userRepo.GetLastSeenAtFn = func(_ context.Context, _ uuid.UUID) (*time.Time, error) {
// 		return &lastSeen, nil
// 	}

// 	status, err := e.svc.GetUserStatus(context.Background(), userA, userB)

// 	require.NoError(t, err)
// 	assert.False(t, status.IsOnline)
// 	assert.NotNil(t, status.LastSeenAt)
// 	assert.Equal(t, lastSeen, *status.LastSeenAt)
// }

// func TestGetUserStatus_Blocked(t *testing.T) {
// 	e := setup(t)
// 	e.modRepo.IsBlockedFn = func(_ context.Context, a, b uuid.UUID) (bool, error) {
// 		return true, nil
// 	}

// 	status, err := e.svc.GetUserStatus(context.Background(), userA, userB)

// 	assert.Error(t, err)
// 	assert.Nil(t, status)
// 	assert.Contains(t, err.Error(), "forbidden")
// }

// // ═════════════════════════════════════════════════════════════
// //  GetUserWithBlockCheck
// // ═════════════════════════════════════════════════════════════

// func TestGetUserWithBlockCheck_Success(t *testing.T) {
// 	e := setup(t)

// 	resp, err := e.svc.GetUserWithBlockCheck(context.Background(), userA, userB)

// 	require.NoError(t, err)
// 	assert.Equal(t, userB, resp.ID)
// }

// func TestGetUserWithBlockCheck_Blocked(t *testing.T) {
// 	e := setup(t)
// 	e.modRepo.IsBlockedFn = func(_ context.Context, a, b uuid.UUID) (bool, error) {
// 		if a == userB && b == userA {
// 			return true, nil // userB blocked userA
// 		}
// 		return false, nil
// 	}

// 	resp, err := e.svc.GetUserWithBlockCheck(context.Background(), userA, userB)

// 	assert.Error(t, err)
// 	assert.Nil(t, resp)
// 	assert.Contains(t, err.Error(), "forbidden")
// }

// func TestGetUserWithBlockCheck_ReverseBlock(t *testing.T) {
// 	e := setup(t)
// 	e.modRepo.IsBlockedFn = func(_ context.Context, a, b uuid.UUID) (bool, error) {
// 		if a == userA && b == userB {
// 			return true, nil // userA blocked userB — тоже нельзя смотреть
// 		}
// 		return false, nil
// 	}

// 	resp, err := e.svc.GetUserWithBlockCheck(context.Background(), userA, userB)

// 	assert.Error(t, err)
// 	assert.Nil(t, resp)
// }