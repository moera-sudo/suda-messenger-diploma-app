package mocks

import (
	"context"

	"github.com/google/uuid"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)

type MockMemberRepo struct {
	GetChatMembersFn         func(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error)
	GetChatMembersDetailedFn func(ctx context.Context, chatID uuid.UUID) ([]chat.ChatMemberInfoResponse, error)
	GetMemberRoleFn          func(ctx context.Context, chatID, userID uuid.UUID) (string, error)
	IsMemberFn               func(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
	IsMemberByMessageIDFn    func(ctx context.Context, userID uuid.UUID, messageID int64) (bool, error)
	AddMemberFn              func(ctx context.Context, chatID, userID uuid.UUID, role string) error
	RemoveMemberFn           func(ctx context.Context, chatID, userID uuid.UUID) error
	UpdateMemberRoleFn       func(ctx context.Context, chatID, userID uuid.UUID, role string) error
	UpdateReadCursorFn       func(ctx context.Context, userID, chatID uuid.UUID, lastReadID int64) error
	GetMemberCountFn         func(ctx context.Context, chatID uuid.UUID) (int, error)
	GetMessageReadersFn      func(ctx context.Context, chatID uuid.UUID, messageID int64) ([]chat.MessageReaderInfoResponse, error)

	AddMemberCalls    int
	RemoveMemberCalls int
	UpdateRoleCalls   int
	IsMemberCalls     int
}

func NewMockMemberRepo() *MockMemberRepo { return &MockMemberRepo{} }

func (m *MockMemberRepo) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	if m.GetChatMembersFn != nil {
		return m.GetChatMembersFn(ctx, chatID)
	}
	return nil, nil
}

func (m *MockMemberRepo) GetChatMembersDetailed(ctx context.Context, chatID uuid.UUID) ([]chat.ChatMemberInfoResponse, error) {
	if m.GetChatMembersDetailedFn != nil {
		return m.GetChatMembersDetailedFn(ctx, chatID)
	}
	return nil, nil
}

func (m *MockMemberRepo) GetMemberRole(ctx context.Context, chatID, userID uuid.UUID) (string, error) {
	if m.GetMemberRoleFn != nil {
		return m.GetMemberRoleFn(ctx, chatID, userID)
	}
	return chat.RoleMember, nil
}

func (m *MockMemberRepo) IsMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	m.IsMemberCalls++
	if m.IsMemberFn != nil {
		return m.IsMemberFn(ctx, chatID, userID)
	}
	return true, nil
}

func (m *MockMemberRepo) IsMemberByMessageID(ctx context.Context, userID uuid.UUID, messageID int64) (bool, error) {
	if m.IsMemberByMessageIDFn != nil {
		return m.IsMemberByMessageIDFn(ctx, userID, messageID)
	}
	return true, nil
}

func (m *MockMemberRepo) AddMember(ctx context.Context, chatID, userID uuid.UUID, role string) error {
	m.AddMemberCalls++
	if m.AddMemberFn != nil {
		return m.AddMemberFn(ctx, chatID, userID, role)
	}
	return nil
}

func (m *MockMemberRepo) RemoveMember(ctx context.Context, chatID, userID uuid.UUID) error {
	m.RemoveMemberCalls++
	if m.RemoveMemberFn != nil {
		return m.RemoveMemberFn(ctx, chatID, userID)
	}
	return nil
}

func (m *MockMemberRepo) UpdateMemberRole(ctx context.Context, chatID, userID uuid.UUID, role string) error {
	m.UpdateRoleCalls++
	if m.UpdateMemberRoleFn != nil {
		return m.UpdateMemberRoleFn(ctx, chatID, userID, role)
	}
	return nil
}

func (m *MockMemberRepo) UpdateReadCursor(ctx context.Context, userID, chatID uuid.UUID, lastReadID int64) error {
	if m.UpdateReadCursorFn != nil {
		return m.UpdateReadCursorFn(ctx, userID, chatID, lastReadID)
	}
	return nil
}

func (m *MockMemberRepo) GetMemberCount(ctx context.Context, chatID uuid.UUID) (int, error) {
	if m.GetMemberCountFn != nil {
		return m.GetMemberCountFn(ctx, chatID)
	}
	return 2, nil
}

func (m *MockMemberRepo) GetMessageReaders(ctx context.Context, chatID uuid.UUID, messageID int64) ([]chat.MessageReaderInfoResponse, error) {
	if m.GetMessageReadersFn != nil {
		return m.GetMessageReadersFn(ctx, chatID, messageID)
	}
	return nil, nil
}