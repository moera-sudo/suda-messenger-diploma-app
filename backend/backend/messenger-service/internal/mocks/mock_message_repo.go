package mocks

import (
	"context"

	"github.com/google/uuid"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)

type MockMessageRepo struct {
	SaveMessageFn       func(ctx context.Context, msg *chat.Message) error
	GetMessageByIDFn    func(ctx context.Context, messageID int64) (*chat.Message, error)
	GetMessagePreviewFn func(ctx context.Context, messageID int64) (*chat.MessagePreview, error)
	GetMessagesFn      func(ctx context.Context, chatID, userID uuid.UUID, limit, offset int) ([]chat.Message, error)
	EditMessageFn      func(ctx context.Context, messageID int64, content string) error
	SoftDeleteMsgFn    func(ctx context.Context, messageID int64) error
	HideMessageFn      func(ctx context.Context, userID uuid.UUID, messageID int64) error
	SearchMessagesFn   func(ctx context.Context, chatID uuid.UUID, query string, limit int) ([]chat.Message, error)
	GetChatMediaFn     func(ctx context.Context, chatID uuid.UUID) (*chat.ChatMediaResponse, error)

	SaveMessageCalls    int
	EditMessageCalls    int
	SoftDeleteMsgCalls  int
	HideMessageCalls    int
}

func NewMockMessageRepo() *MockMessageRepo { return &MockMessageRepo{} }

func (m *MockMessageRepo) SaveMessage(ctx context.Context, msg *chat.Message) error {
	m.SaveMessageCalls++
	if m.SaveMessageFn != nil {
		return m.SaveMessageFn(ctx, msg)
	}
	msg.ID = int64(m.SaveMessageCalls)
	msg.Status = chat.StatusSent
	return nil
}

func (m *MockMessageRepo) GetMessageByID(ctx context.Context, messageID int64) (*chat.Message, error) {
	if m.GetMessageByIDFn != nil {
		return m.GetMessageByIDFn(ctx, messageID)
	}
	return nil, nil
}

func (m *MockMessageRepo) GetMessagePreview(ctx context.Context, messageID int64) (*chat.MessagePreview, error) {
	if m.GetMessagePreviewFn != nil {
		return m.GetMessagePreviewFn(ctx, messageID)
	}
	return nil, nil
}

func (m *MockMessageRepo) GetMessages(ctx context.Context, chatID, userID uuid.UUID, limit, offset int) ([]chat.Message, error) {
	if m.GetMessagesFn != nil {
		return m.GetMessagesFn(ctx, chatID, userID, limit, offset)
	}
	return nil, nil
}

func (m *MockMessageRepo) EditMessage(ctx context.Context, messageID int64, content string) error {
	m.EditMessageCalls++
	if m.EditMessageFn != nil {
		return m.EditMessageFn(ctx, messageID, content)
	}
	return nil
}

func (m *MockMessageRepo) SoftDeleteMessage(ctx context.Context, messageID int64) error {
	m.SoftDeleteMsgCalls++
	if m.SoftDeleteMsgFn != nil {
		return m.SoftDeleteMsgFn(ctx, messageID)
	}
	return nil
}

func (m *MockMessageRepo) HideMessageForUser(ctx context.Context, userID uuid.UUID, messageID int64) error {
	m.HideMessageCalls++
	if m.HideMessageFn != nil {
		return m.HideMessageFn(ctx, userID, messageID)
	}
	return nil
}

func (m *MockMessageRepo) SearchMessages(ctx context.Context, chatID uuid.UUID, query string, limit int) ([]chat.Message, error) {
	if m.SearchMessagesFn != nil {
		return m.SearchMessagesFn(ctx, chatID, query, limit)
	}
	return nil, nil
}

func (m *MockMessageRepo) GetChatMedia(ctx context.Context, chatID uuid.UUID) (*chat.ChatMediaResponse, error) {
	if m.GetChatMediaFn != nil {
		return m.GetChatMediaFn(ctx, chatID)
	}
	return &chat.ChatMediaResponse{}, nil
}