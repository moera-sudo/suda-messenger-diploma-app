package service_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/search"
	searchService "github.com/moera-sudo/backend/backend/messenger-service/internal/features/search/service"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)

var user1 = uuid.MustParse("11111111-1111-1111-1111-111111111111")
var chatID = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")

// ── Mock Search Repo ────────────────────────────────────────

type MockSearchRepo struct {
	SearchUsersFn          func(ctx context.Context, query string, limit int) ([]search.SearchResult, error)
	SearchChatsFn          func(ctx context.Context, query string, limit int) ([]search.SearchResult, error)
	SearchMessagesGlobalFn func(ctx context.Context, userID uuid.UUID, query string, limit int) ([]search.SearchResult, error)
	SearchMessagesInChatFn func(ctx context.Context, chatID uuid.UUID, query string, limit int) ([]search.SearchResult, error)
}

func (m *MockSearchRepo) SearchUsers(ctx context.Context, q string, l int) ([]search.SearchResult, error) {
	if m.SearchUsersFn != nil { return m.SearchUsersFn(ctx, q, l) }
	return []search.SearchResult{}, nil
}

func (m *MockSearchRepo) SearchChats(ctx context.Context, q string, l int) ([]search.SearchResult, error) {
	if m.SearchChatsFn != nil { return m.SearchChatsFn(ctx, q, l) }
	return []search.SearchResult{}, nil
}

func (m *MockSearchRepo) SearchMessagesGlobal(ctx context.Context, uid uuid.UUID, q string, l int) ([]search.SearchResult, error) {
	if m.SearchMessagesGlobalFn != nil { return m.SearchMessagesGlobalFn(ctx, uid, q, l) }
	return []search.SearchResult{}, nil
}

func (m *MockSearchRepo) SearchMessagesInChat(ctx context.Context, cid uuid.UUID, q string, l int) ([]search.SearchResult, error) {
	if m.SearchMessagesInChatFn != nil { return m.SearchMessagesInChatFn(ctx, cid, q, l) }
	return []search.SearchResult{}, nil
}

// ── Mock Member Repo ────────────────────────────────────────

type MockMemberRepo struct {
	IsMemberFn func(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
}

func (m *MockMemberRepo) IsMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	if m.IsMemberFn != nil { return m.IsMemberFn(ctx, chatID, userID) }
	return true, nil
}

// Stub all other MemberRepo methods
func (m *MockMemberRepo) GetChatMembers(_ context.Context, _ uuid.UUID) ([]uuid.UUID, error) { return nil, nil }

func (m *MockMemberRepo) GetChatMembersDetailed(_ context.Context, _ uuid.UUID) ([]chat.ChatMemberInfoResponse, error) {
	return nil, nil
}
func (m *MockMemberRepo) GetMemberRole(_ context.Context, _, _ uuid.UUID) (string, error) { return "", nil }
func (m *MockMemberRepo) IsMemberByMessageID(_ context.Context, _ uuid.UUID, _ int64) (bool, error) { return true, nil }
func (m *MockMemberRepo) AddMember(_ context.Context, _, _ uuid.UUID, _ string) error { return nil }
func (m *MockMemberRepo) RemoveMember(_ context.Context, _, _ uuid.UUID) error { return nil }
func (m *MockMemberRepo) UpdateMemberRole(_ context.Context, _, _ uuid.UUID, _ string) error { return nil }
func (m *MockMemberRepo) UpdateReadCursor(_ context.Context, _, _ uuid.UUID, _ int64) error { return nil }
func (m *MockMemberRepo) GetMemberCount(_ context.Context, _ uuid.UUID) (int, error) { return 0, nil }
func (m *MockMemberRepo) GetMessageReaders(_ context.Context, _ uuid.UUID, _ int64) ([]chat.MessageReaderInfoResponse, error) {
	return nil, nil
}
// ═════════════════════════════════════════════════════════════
//  GlobalSearch
// ═════════════════════════════════════════════════════════════

func TestGlobalSearch_TooShortQuery(t *testing.T) {
	repo := &MockSearchRepo{}
	members := &MockMemberRepo{}
	svc := searchService.NewService(repo, members)

	resp, err := svc.GlobalSearch(context.Background(), user1, "a")

	require.NoError(t, err)
	assert.Empty(t, resp.Users)
	assert.Empty(t, resp.Chats)
	assert.Empty(t, resp.Messages)
}

func TestGlobalSearch_ReturnsAllCategories(t *testing.T) {
	repo := &MockSearchRepo{
		SearchUsersFn: func(_ context.Context, _ string, _ int) ([]search.SearchResult, error) {
			return []search.SearchResult{
				{Type: search.ResultTypeUser, ID: "1", Title: "Alice", Rank: 0.9},
			}, nil
		},
		SearchChatsFn: func(_ context.Context, _ string, _ int) ([]search.SearchResult, error) {
			return []search.SearchResult{
				{Type: search.ResultTypeGroup, ID: "2", Title: "Developers", Rank: 0.8},
			}, nil
		},
		SearchMessagesGlobalFn: func(_ context.Context, _ uuid.UUID, _ string, _ int) ([]search.SearchResult, error) {
			return []search.SearchResult{
				{Type: search.ResultTypeMessage, ID: "100", Title: "Alice", Subtitle: "Hello!", Rank: 0.7},
			}, nil
		},
	}
	members := &MockMemberRepo{}
	svc := searchService.NewService(repo, members)

	resp, err := svc.GlobalSearch(context.Background(), user1, "test")

	require.NoError(t, err)
	assert.Len(t, resp.Users, 1)
	assert.Len(t, resp.Chats, 1)
	assert.Len(t, resp.Messages, 1)
	assert.Equal(t, "Alice", resp.Users[0].Title)
}

func TestGlobalSearch_DBError_ReturnsPartial(t *testing.T) {
	repo := &MockSearchRepo{
		SearchUsersFn: func(_ context.Context, _ string, _ int) ([]search.SearchResult, error) {
			return nil, fmt.Errorf("connection refused")
		},
		SearchChatsFn: func(_ context.Context, _ string, _ int) ([]search.SearchResult, error) {
			return []search.SearchResult{{Type: search.ResultTypeGroup, ID: "1", Title: "Test"}}, nil
		},
	}
	members := &MockMemberRepo{}
	svc := searchService.NewService(repo, members)

	resp, err := svc.GlobalSearch(context.Background(), user1, "test")

	require.NoError(t, err)
	assert.Empty(t, resp.Users) // Ошибка — пустой массив, не nil
	assert.Len(t, resp.Chats, 1)
}

// ═════════════════════════════════════════════════════════════
//  SearchInChat
// ═════════════════════════════════════════════════════════════

func TestSearchInChat_Success(t *testing.T) {
	repo := &MockSearchRepo{
		SearchMessagesInChatFn: func(_ context.Context, _ uuid.UUID, _ string, _ int) ([]search.SearchResult, error) {
			return []search.SearchResult{
				{Type: search.ResultTypeMessage, ID: "10", Title: "Alice", Highlight: "Hello <b>world</b>!"},
				{Type: search.ResultTypeMessage, ID: "5", Title: "Bob", Highlight: "<b>World</b> peace"},
			}, nil
		},
	}
	members := &MockMemberRepo{}
	svc := searchService.NewService(repo, members)

	resp, err := svc.SearchInChat(context.Background(), user1, chatID, "world")

	require.NoError(t, err)
	assert.Len(t, resp.Messages, 2)
	assert.Equal(t, 2, resp.Total)
	assert.Contains(t, resp.Messages[0].Highlight, "<b>")
}

func TestSearchInChat_NotMember(t *testing.T) {
	repo := &MockSearchRepo{}
	members := &MockMemberRepo{
		IsMemberFn: func(_ context.Context, _, _ uuid.UUID) (bool, error) {
			return false, nil
		},
	}
	svc := searchService.NewService(repo, members)

	resp, err := svc.SearchInChat(context.Background(), user1, chatID, "test")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestSearchInChat_TooShort(t *testing.T) {
	repo := &MockSearchRepo{}
	members := &MockMemberRepo{}
	svc := searchService.NewService(repo, members)

	resp, err := svc.SearchInChat(context.Background(), user1, chatID, "x")

	require.NoError(t, err)
	assert.Empty(t, resp.Messages)
	assert.Equal(t, 0, resp.Total)
}