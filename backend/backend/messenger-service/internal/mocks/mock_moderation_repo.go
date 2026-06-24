package mocks

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)

type MockModerationRepo struct {
	MuteChatFn         func(ctx context.Context, userID, chatID uuid.UUID, mutedUntil *time.Time) error
	UnmuteChatFn       func(ctx context.Context, userID, chatID uuid.UUID) error
	IsChatMutedFn      func(ctx context.Context, userID, chatID uuid.UUID) (bool, error)
	ClearAndHideChatFn func(ctx context.Context, userID, chatID uuid.UUID) error
	BlockUserFn       func(ctx context.Context, blockerID, blockedID uuid.UUID) error
	UnblockUserFn     func(ctx context.Context, blockerID, blockedID uuid.UUID) error
	IsBlockedFn       func(ctx context.Context, blockerID, blockedID uuid.UUID) (bool, error)
	GetBlockedUsersFn func(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetBlockedUsersDetailedFn func(ctx context.Context, userID uuid.UUID) ([]chat.BlockedUserInfo, error)
	PinMessageFn      func(ctx context.Context, chatID uuid.UUID, messageID int64, pinnedBy uuid.UUID) error
	UnpinMessageFn    func(ctx context.Context, chatID uuid.UUID, messageID int64) error
	GetPinnedMsgsFn   func(ctx context.Context, chatID uuid.UUID) ([]chat.PinnedMessageInfoResponse, error)
	SetContactNameFn  func(ctx context.Context, ownerID, contactID uuid.UUID, name string) error
	RemoveContactFn   func(ctx context.Context, ownerID, contactID uuid.UUID) error
	GetContactNameFn  func(ctx context.Context, ownerID, contactID uuid.UUID) (string, error)
	GetContactsFn     func(ctx context.Context, ownerID uuid.UUID) ([]chat.Contact, error)

	BlockUserCalls   int
	UnblockUserCalls int
	PinMessageCalls  int
}

func NewMockModerationRepo() *MockModerationRepo { return &MockModerationRepo{} }

func (m *MockModerationRepo) MuteChat(ctx context.Context, userID, chatID uuid.UUID, mutedUntil *time.Time) error {
	if m.MuteChatFn != nil { return m.MuteChatFn(ctx, userID, chatID, mutedUntil) }
	return nil
}

func (m *MockModerationRepo) UnmuteChat(ctx context.Context, userID, chatID uuid.UUID) error {
	if m.UnmuteChatFn != nil { return m.UnmuteChatFn(ctx, userID, chatID) }
	return nil
}

func (m *MockModerationRepo) IsChatMuted(ctx context.Context, userID, chatID uuid.UUID) (bool, error) {
	if m.IsChatMutedFn != nil { return m.IsChatMutedFn(ctx, userID, chatID) }
	return false, nil
}

func (m *MockModerationRepo) ClearAndHideChat(ctx context.Context, userID, chatID uuid.UUID) error {
	if m.ClearAndHideChatFn != nil { return m.ClearAndHideChatFn(ctx, userID, chatID) }
	return nil
}

func (m *MockModerationRepo) BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	m.BlockUserCalls++
	if m.BlockUserFn != nil { return m.BlockUserFn(ctx, blockerID, blockedID) }
	return nil
}

func (m *MockModerationRepo) UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	m.UnblockUserCalls++
	if m.UnblockUserFn != nil { return m.UnblockUserFn(ctx, blockerID, blockedID) }
	return nil
}

func (m *MockModerationRepo) IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) (bool, error) {
	if m.IsBlockedFn != nil { return m.IsBlockedFn(ctx, blockerID, blockedID) }
	return false, nil
}

func (m *MockModerationRepo) GetBlockedUsers(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	if m.GetBlockedUsersFn != nil { return m.GetBlockedUsersFn(ctx, userID) }
	return nil, nil
}

func (m *MockModerationRepo) GetBlockedUsersDetailed(ctx context.Context, userID uuid.UUID) ([]chat.BlockedUserInfo, error) {
	if m.GetBlockedUsersDetailedFn != nil { return m.GetBlockedUsersDetailedFn(ctx, userID) }
	return nil, nil
}

func (m *MockModerationRepo) PinMessage(ctx context.Context, chatID uuid.UUID, messageID int64, pinnedBy uuid.UUID) error {
	m.PinMessageCalls++
	if m.PinMessageFn != nil { return m.PinMessageFn(ctx, chatID, messageID, pinnedBy) }
	return nil
}

func (m *MockModerationRepo) UnpinMessage(ctx context.Context, chatID uuid.UUID, messageID int64) error {
	if m.UnpinMessageFn != nil { return m.UnpinMessageFn(ctx, chatID, messageID) }
	return nil
}

func (m *MockModerationRepo) GetPinnedMessages(ctx context.Context, chatID uuid.UUID) ([]chat.PinnedMessageInfoResponse, error) {
	if m.GetPinnedMsgsFn != nil { return m.GetPinnedMsgsFn(ctx, chatID) }
	return nil, nil
}

func (m *MockModerationRepo) SetContactName(ctx context.Context, ownerID, contactID uuid.UUID, name string) error {
	if m.SetContactNameFn != nil { return m.SetContactNameFn(ctx, ownerID, contactID, name) }
	return nil
}

func (m *MockModerationRepo) RemoveContactName(ctx context.Context, ownerID, contactID uuid.UUID) error {
	if m.RemoveContactFn != nil { return m.RemoveContactFn(ctx, ownerID, contactID) }
	return nil
}

func (m *MockModerationRepo) GetContactName(ctx context.Context, ownerID, contactID uuid.UUID) (string, error) {
	if m.GetContactNameFn != nil { return m.GetContactNameFn(ctx, ownerID, contactID) }
	return "", nil
}

func (m *MockModerationRepo) GetContacts(ctx context.Context, ownerID uuid.UUID) ([]chat.Contact, error) {
	if m.GetContactsFn != nil { return m.GetContactsFn(ctx, ownerID) }
	return nil, nil
}