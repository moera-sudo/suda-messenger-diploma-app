package mocks

import (
	"context"

	"github.com/google/uuid"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)

type MockChatRepo struct {
	CreateChatFn    func(ctx context.Context, c *chat.Chat, memberIDs []uuid.UUID, creatorRole string) error
	GetChatByIDFn   func(ctx context.Context, chatID uuid.UUID) (*chat.Chat, error)
	GetChatTypeFn   func(ctx context.Context, chatID uuid.UUID) (string, error)
	GetDirectChatFn func(ctx context.Context, userA, userB uuid.UUID) (*chat.Chat, error)
	GetSavedChatFn  func(ctx context.Context, userID uuid.UUID) (*chat.Chat, error)
	GetUserChatsFn  func(ctx context.Context, userID uuid.UUID) ([]chat.ChatListItemResponse, error)
	UpdateChatFn        func(ctx context.Context, chatID uuid.UUID, req *chat.UpdateChatReq) error
	DeleteChatFn        func(ctx context.Context, chatID uuid.UUID) error
	IsUsernameTakenFn   func(ctx context.Context, username string) (bool, error)

	CreateChatCalls int
}

func NewMockChatRepo() *MockChatRepo { return &MockChatRepo{} }

func (m *MockChatRepo) CreateChat(ctx context.Context, c *chat.Chat, memberIDs []uuid.UUID, creatorRole string) error {
	m.CreateChatCalls++
	if m.CreateChatFn != nil {
		return m.CreateChatFn(ctx, c, memberIDs, creatorRole)
	}
	return nil
}

func (m *MockChatRepo) GetChatByID(ctx context.Context, chatID uuid.UUID) (*chat.Chat, error) {
	if m.GetChatByIDFn != nil {
		return m.GetChatByIDFn(ctx, chatID)
	}
	return nil, nil
}

func (m *MockChatRepo) GetChatType(ctx context.Context, chatID uuid.UUID) (string, error) {
	if m.GetChatTypeFn != nil {
		return m.GetChatTypeFn(ctx, chatID)
	}
	return chat.TypeGroup, nil
}

func (m *MockChatRepo) GetDirectChat(ctx context.Context, userA, userB uuid.UUID) (*chat.Chat, error) {
	if m.GetDirectChatFn != nil {
		return m.GetDirectChatFn(ctx, userA, userB)
	}
	return nil, nil
}

func (m *MockChatRepo) GetSavedChat(ctx context.Context, userID uuid.UUID) (*chat.Chat, error) {
	if m.GetSavedChatFn != nil {
		return m.GetSavedChatFn(ctx, userID)
	}
	return nil, nil
}

func (m *MockChatRepo) GetUserChats(ctx context.Context, userID uuid.UUID) ([]chat.ChatListItemResponse, error) {
	if m.GetUserChatsFn != nil {
		return m.GetUserChatsFn(ctx, userID)
	}
	return nil, nil
}

func (m *MockChatRepo) UpdateChat(ctx context.Context, chatID uuid.UUID, req *chat.UpdateChatReq) error {
	if m.UpdateChatFn != nil {
		return m.UpdateChatFn(ctx, chatID, req)
	}
	return nil
}

func (m *MockChatRepo) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	if m.DeleteChatFn != nil {
		return m.DeleteChatFn(ctx, chatID)
	}
	return nil
}

func (m *MockChatRepo) IsUsernameTaken(ctx context.Context, username string) (bool, error) {
	if m.IsUsernameTakenFn != nil {
		return m.IsUsernameTakenFn(ctx, username)
	}
	return false, nil
}