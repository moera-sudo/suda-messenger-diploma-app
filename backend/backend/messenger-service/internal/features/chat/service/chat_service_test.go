package service_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/service"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/mocks"
)

var (
	user1  = uuid.MustParse("11111111-1111-1111-1111-111111111111")
	user2  = uuid.MustParse("22222222-2222-2222-2222-222222222222")
	user3  = uuid.MustParse("33333333-3333-3333-3333-333333333333")
	chatID = uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
)

type testEnv struct {
	svc     *service.ChatService
	chats   *mocks.MockChatRepo
	msgs    *mocks.MockMessageRepo
	members *mocks.MockMemberRepo
	mods    *mocks.MockModerationRepo
	redis   *redis.Client
	mr      *miniredis.Miniredis
}

func setup(t *testing.T) *testEnv {
	mr := miniredis.RunT(t)
	rc := redis.NewClient(&redis.Options{Addr: mr.Addr()})

	chats := mocks.NewMockChatRepo()
	msgs := mocks.NewMockMessageRepo()
	members := mocks.NewMockMemberRepo()
	mods := mocks.NewMockModerationRepo()

	svc := service.NewChatService(chats, msgs, members, mods, nil, nil, nil, rc, nil, nil, nil)

	return &testEnv{
		svc: svc, chats: chats, msgs: msgs,
		members: members, mods: mods, redis: rc, mr: mr,
	}
}

// ═════════════════════════════════════════════════════════════
//  CreateChat
// ═════════════════════════════════════════════════════════════

func TestCreateChat_Saved_New(t *testing.T) {
	e := setup(t)
	e.chats.GetSavedChatFn = func(_ context.Context, _ uuid.UUID) (*chat.Chat, error) {
		return nil, nil
	}

	result, err := e.svc.CreateChat(context.Background(), user1, &chat.CreateChatReq{Type: chat.TypeSaved})

	require.NoError(t, err)
	assert.Equal(t, chat.TypeSaved, result.Type)
	assert.Equal(t, "Saved Messages", result.Name)
	assert.Equal(t, 1, e.chats.CreateChatCalls)
}

func TestCreateChat_Saved_Exists(t *testing.T) {
	e := setup(t)
	existing := &chat.Chat{ID: uuid.New(), Type: chat.TypeSaved, Name: "Saved Messages"}
	e.chats.GetSavedChatFn = func(_ context.Context, _ uuid.UUID) (*chat.Chat, error) {
		return existing, nil
	}

	result, err := e.svc.CreateChat(context.Background(), user1, &chat.CreateChatReq{Type: chat.TypeSaved})

	require.NoError(t, err)
	assert.Equal(t, existing.ID, result.ID)
	assert.Equal(t, 0, e.chats.CreateChatCalls)
}

func TestCreateChat_Direct_New(t *testing.T) {
	e := setup(t)
	e.chats.GetDirectChatFn = func(_ context.Context, _, _ uuid.UUID) (*chat.Chat, error) {
		return nil, nil
	}

	result, err := e.svc.CreateChat(context.Background(), user1, &chat.CreateChatReq{
		Type: chat.TypeDirect, TargetID: user2,
	})

	require.NoError(t, err)
	assert.Equal(t, chat.TypeDirect, result.Type)
	assert.Equal(t, 1, e.chats.CreateChatCalls)
}

func TestCreateChat_Direct_Blocked(t *testing.T) {
	e := setup(t)
	e.mods.IsBlockedFn = func(_ context.Context, blockerID, blockedID uuid.UUID) (bool, error) {
		return true, nil
	}

	result, err := e.svc.CreateChat(context.Background(), user1, &chat.CreateChatReq{
		Type: chat.TypeDirect, TargetID: user2,
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "blocked")
}

func TestCreateChat_Group(t *testing.T) {
	e := setup(t)

	result, err := e.svc.CreateChat(context.Background(), user1, &chat.CreateChatReq{
		Type: chat.TypeGroup, Name: "Test Group", MemberIDs: []uuid.UUID{user2, user3},
	})

	require.NoError(t, err)
	assert.Equal(t, chat.TypeGroup, result.Type)
	assert.Equal(t, "Test Group", result.Name)
}

func TestCreateChat_Group_NoName(t *testing.T) {
	e := setup(t)

	result, err := e.svc.CreateChat(context.Background(), user1, &chat.CreateChatReq{
		Type: chat.TypeGroup, MemberIDs: []uuid.UUID{user2},
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "group name is required")
}

func TestCreateChat_InvalidType(t *testing.T) {
	e := setup(t)

	result, err := e.svc.CreateChat(context.Background(), user1, &chat.CreateChatReq{Type: "INVALID"})

	assert.Error(t, err)
	assert.Nil(t, result)
}

// ═════════════════════════════════════════════════════════════
//  SendMessage
// ═════════════════════════════════════════════════════════════

func TestSendMessage_Success(t *testing.T) {
	e := setup(t)
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeGroup, nil
	}

	msg, err := e.svc.SendMessage(context.Background(), user1, &chat.SendMessageReq{
		ChatID: chatID, Content: "Hello!", Type: chat.MsgTypeText,
	})

	require.NoError(t, err)
	assert.Equal(t, "Hello!", msg.Content)
	assert.Equal(t, chat.StatusSent, msg.Status)
	assert.Equal(t, 1, e.msgs.SaveMessageCalls)
}

func TestSendMessage_NotMember(t *testing.T) {
	e := setup(t)
	e.members.IsMemberFn = func(_ context.Context, _, _ uuid.UUID) (bool, error) {
		return false, nil
	}
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeGroup, nil
	}

	msg, err := e.svc.SendMessage(context.Background(), user1, &chat.SendMessageReq{
		ChatID: chatID, Content: "Hello!", Type: chat.MsgTypeText,
	})

	assert.Error(t, err)
	assert.Nil(t, msg)
	assert.Contains(t, err.Error(), "not a member")
	assert.Equal(t, 0, e.msgs.SaveMessageCalls)
}

func TestSendMessage_BlockedInDirect(t *testing.T) {
	e := setup(t)
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeDirect, nil
	}
	e.members.GetChatMembersFn = func(_ context.Context, _ uuid.UUID) ([]uuid.UUID, error) {
		return []uuid.UUID{user1, user2}, nil
	}
	e.mods.IsBlockedFn = func(_ context.Context, _, _ uuid.UUID) (bool, error) {
		return true, nil
	}

	msg, err := e.svc.SendMessage(context.Background(), user1, &chat.SendMessageReq{
		ChatID: chatID, Content: "Hello!", Type: chat.MsgTypeText,
	})

	assert.Error(t, err)
	assert.Nil(t, msg)
	assert.Contains(t, err.Error(), "blocked")
}

// ═════════════════════════════════════════════════════════════
//  EditMessage
// ═════════════════════════════════════════════════════════════

func TestEditMessage_Success(t *testing.T) {
	e := setup(t)
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		return &chat.Message{ID: 1, ChatID: chatID, SenderID: &user1, Type: chat.MsgTypeText}, nil
	}

	err := e.svc.EditMessage(context.Background(), user1, 1, &chat.EditMessageReq{Content: "Edited!"})

	require.NoError(t, err)
	assert.Equal(t, 1, e.msgs.EditMessageCalls)
}

func TestEditMessage_NotAuthor(t *testing.T) {
	e := setup(t)
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		return &chat.Message{ID: 1, ChatID: chatID, SenderID: &user1, Type: chat.MsgTypeText}, nil
	}

	err := e.svc.EditMessage(context.Background(), user2, 1, &chat.EditMessageReq{Content: "Hacked!"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "can only edit own")
	assert.Equal(t, 0, e.msgs.EditMessageCalls)
}

func TestEditMessage_SystemMessage(t *testing.T) {
	e := setup(t)
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		return &chat.Message{ID: 1, ChatID: chatID, SenderID: &user1, Type: chat.MsgTypeSystem}, nil
	}

	err := e.svc.EditMessage(context.Background(), user1, 1, &chat.EditMessageReq{Content: "Hack"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "system")
}

// ═════════════════════════════════════════════════════════════
//  DeleteMessage
// ═════════════════════════════════════════════════════════════

func TestDeleteMessage_ForMe(t *testing.T) {
	e := setup(t)
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		return &chat.Message{ID: 1, ChatID: chatID, SenderID: &user2, Type: chat.MsgTypeText}, nil
	}

	err := e.svc.DeleteMessage(context.Background(), user1, 1, &chat.DeleteMessageReq{ForEveryone: false})

	require.NoError(t, err)
	assert.Equal(t, 1, e.msgs.HideMessageCalls)
	assert.Equal(t, 0, e.msgs.SoftDeleteMsgCalls)
}

func TestDeleteMessage_ForEveryone_Author(t *testing.T) {
	e := setup(t)
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		return &chat.Message{ID: 1, ChatID: chatID, SenderID: &user1, Type: chat.MsgTypeText}, nil
	}

	err := e.svc.DeleteMessage(context.Background(), user1, 1, &chat.DeleteMessageReq{ForEveryone: true})

	require.NoError(t, err)
	assert.Equal(t, 1, e.msgs.SoftDeleteMsgCalls)
	assert.Equal(t, 0, e.msgs.HideMessageCalls)
}

func TestDeleteMessage_ForEveryone_NotAuthorNotAdmin(t *testing.T) {
	e := setup(t)
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		return &chat.Message{ID: 1, ChatID: chatID, SenderID: &user2, Type: chat.MsgTypeText}, nil
	}
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeGroup, nil
	}
	e.members.GetMemberRoleFn = func(_ context.Context, _, _ uuid.UUID) (string, error) {
		return chat.RoleMember, nil
	}

	err := e.svc.DeleteMessage(context.Background(), user1, 1, &chat.DeleteMessageReq{ForEveryone: true})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "insufficient rights")
}

func TestDeleteMessage_ForEveryone_AdminCanDelete(t *testing.T) {
	e := setup(t)
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		return &chat.Message{ID: 1, ChatID: chatID, SenderID: &user2, Type: chat.MsgTypeText}, nil
	}
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeGroup, nil
	}
	e.members.GetMemberRoleFn = func(_ context.Context, _, userID uuid.UUID) (string, error) {
		if userID == user1 {
			return chat.RoleAdmin, nil
		}
		return chat.RoleMember, nil
	}

	err := e.svc.DeleteMessage(context.Background(), user1, 1, &chat.DeleteMessageReq{ForEveryone: true})

	require.NoError(t, err)
	assert.Equal(t, 1, e.msgs.SoftDeleteMsgCalls)
}

// ═════════════════════════════════════════════════════════════
//  ForwardMessage
// ═════════════════════════════════════════════════════════════

func TestForwardMessage_Success(t *testing.T) {
	e := setup(t)
	fromChat := uuid.New()
	toChat := uuid.New()

	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeGroup, nil
	}
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		return &chat.Message{ID: 1, ChatID: fromChat, Content: "Original", Type: chat.MsgTypeText}, nil
	}

	msg, err := e.svc.ForwardMessage(context.Background(), user1, &chat.ForwardMessageReq{
		FromChatID: fromChat, ToChatID: toChat, MessageID: 1,
	})

	require.NoError(t, err)
	assert.Equal(t, "Original", msg.Content)
	assert.Equal(t, &fromChat, msg.ForwardedFromChat)
}

// ═════════════════════════════════════════════════════════════
//  Members
// ═════════════════════════════════════════════════════════════

func TestAddMember_Success(t *testing.T) {
	e := setup(t)
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeGroup, nil
	}
	e.members.GetMemberRoleFn = func(_ context.Context, _, _ uuid.UUID) (string, error) {
		return chat.RoleOwner, nil
	}

	err := e.svc.AddMember(context.Background(), chatID, user1, user3)

	require.NoError(t, err)
	assert.Equal(t, 1, e.members.AddMemberCalls)
}

func TestAddMember_NotAdmin(t *testing.T) {
	e := setup(t)
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeGroup, nil
	}
	e.members.GetMemberRoleFn = func(_ context.Context, _, _ uuid.UUID) (string, error) {
		return chat.RoleMember, nil
	}

	err := e.svc.AddMember(context.Background(), chatID, user1, user3)

	assert.Error(t, err)
	assert.Equal(t, 0, e.members.AddMemberCalls)
}

func TestRemoveMember_CannotRemoveOwner(t *testing.T) {
	e := setup(t)
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeGroup, nil
	}
	e.members.GetMemberRoleFn = func(_ context.Context, _, userID uuid.UUID) (string, error) {
		if userID == user1 {
			return chat.RoleAdmin, nil
		}
		return chat.RoleOwner, nil // target is owner
	}

	err := e.svc.RemoveMember(context.Background(), chatID, user1, user2)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot remove the chat owner")
	assert.Equal(t, 0, e.members.RemoveMemberCalls)
}

func TestLeaveChat_Direct_Forbidden(t *testing.T) {
	e := setup(t)
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeDirect, nil
	}

	err := e.svc.LeaveChat(context.Background(), chatID, user1)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cannot leave")
}

// ═════════════════════════════════════════════════════════════
//  Block
// ═════════════════════════════════════════════════════════════

func TestBlockUser_Success(t *testing.T) {
	e := setup(t)

	err := e.svc.BlockUser(context.Background(), user1, user2)

	require.NoError(t, err)
	assert.Equal(t, 1, e.mods.BlockUserCalls)
}

// ═════════════════════════════════════════════════════════════
//  Pin
// ═════════════════════════════════════════════════════════════

func TestPinMessage_Success(t *testing.T) {
	e := setup(t)
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeDirect, nil
	}
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		return &chat.Message{ID: 1, ChatID: chatID}, nil
	}

	err := e.svc.PinMessage(context.Background(), chatID, user1, 1)

	require.NoError(t, err)
	assert.Equal(t, 1, e.mods.PinMessageCalls)
}

func TestPinMessage_GroupRequiresAdmin(t *testing.T) {
	e := setup(t)
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeGroup, nil
	}
	e.members.GetMemberRoleFn = func(_ context.Context, _, _ uuid.UUID) (string, error) {
		return chat.RoleMember, nil
	}

	err := e.svc.PinMessage(context.Background(), chatID, user1, 1)

	assert.Error(t, err)
	assert.Equal(t, 0, e.mods.PinMessageCalls)
}

// ═════════════════════════════════════════════════════════════
//  Contacts
// ═════════════════════════════════════════════════════════════

func TestSetContactName_CannotSetSelf(t *testing.T) {
	e := setup(t)

	err := e.svc.SetContactName(context.Background(), user1, user1, "Me")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "yourself")
}

func TestSetContactName_Success(t *testing.T) {
	e := setup(t)

	err := e.svc.SetContactName(context.Background(), user1, user2, "Best Friend")

	require.NoError(t, err)
}

// ═════════════════════════════════════════════════════════════
//  GetChatMessages
// ═════════════════════════════════════════════════════════════

func TestGetChatMessages_NotMember(t *testing.T) {
	e := setup(t)
	e.members.IsMemberFn = func(_ context.Context, _, _ uuid.UUID) (bool, error) {
		return false, nil
	}

	msgs, err := e.svc.GetChatMessages(context.Background(), user1, chatID, 50, 0)

	assert.Error(t, err)
	assert.Nil(t, msgs)
	assert.Contains(t, err.Error(), "not a member")
}

func TestGetChatMessages_Success(t *testing.T) {
	e := setup(t)
	// requireReadAccess loads the chat: a GROUP falls through to the membership check.
	e.chats.GetChatByIDFn = func(_ context.Context, _ uuid.UUID) (*chat.Chat, error) {
		return &chat.Chat{ID: chatID, Type: chat.TypeGroup}, nil
	}
	e.msgs.GetMessagesFn = func(_ context.Context, _, _ uuid.UUID, _, _ int) ([]chat.Message, error) {
		return []chat.Message{
			{ID: 2, Content: "Second"},
			{ID: 1, Content: "First"},
		}, nil
	}

	msgs, err := e.svc.GetChatMessages(context.Background(), user1, chatID, 50, 0)

	require.NoError(t, err)
	assert.Len(t, msgs, 2)
}


// ═════════════════════════════════════════════════════════════
//  Message Readers
// ═════════════════════════════════════════════════════════════

func TestGetMessageReaders_Success(t *testing.T) {
	e := setup(t)
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		sender := user1
		return &chat.Message{ID: 1, ChatID: chatID, SenderID: &sender, Type: chat.MsgTypeText}, nil
	}
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeGroup, nil
	}
	e.members.GetMessageReadersFn = func(_ context.Context, _ uuid.UUID, _ int64) ([]chat.MessageReaderInfoResponse, error) {
		return []chat.MessageReaderInfoResponse{
			{UserID: user1, Username: "@sender", DisplayName: "Sender"},
			{UserID: user2, Username: "@reader1", DisplayName: "Reader 1"},
			{UserID: user3, Username: "@reader2", DisplayName: "Reader 2"},
		}, nil
	}

	readers, err := e.svc.GetMessageReaders(context.Background(), user2, 1)

	require.NoError(t, err)
	// Автор исключён из списка
	assert.Len(t, readers, 2)
	for _, r := range readers {
		assert.NotEqual(t, user1, r.UserID)
	}
}

func TestGetMessageReaders_DirectChat_Forbidden(t *testing.T) {
	e := setup(t)
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		return &chat.Message{ID: 1, ChatID: chatID, Type: chat.MsgTypeText}, nil
	}
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeDirect, nil
	}

	readers, err := e.svc.GetMessageReaders(context.Background(), user1, 1)

	assert.Error(t, err)
	assert.Nil(t, readers)
	assert.Contains(t, err.Error(), "only available for groups")
}

func TestGetMessageReaders_NotMember(t *testing.T) {
	e := setup(t)
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		return &chat.Message{ID: 1, ChatID: chatID, Type: chat.MsgTypeText}, nil
	}
	e.members.IsMemberFn = func(_ context.Context, _, _ uuid.UUID) (bool, error) {
		return false, nil
	}

	readers, err := e.svc.GetMessageReaders(context.Background(), user1, 1)

	assert.Error(t, err)
	assert.Nil(t, readers)
	assert.Contains(t, err.Error(), "not a member")
}

// ═════════════════════════════════════════════════════════════
//  Block Enforcement — AddMember with existing member block
// ═════════════════════════════════════════════════════════════

func TestAddMember_BlockedByExistingMember(t *testing.T) {
	e := setup(t)
	e.chats.GetChatTypeFn = func(_ context.Context, _ uuid.UUID) (string, error) {
		return chat.TypeGroup, nil
	}
	// Requester is admin
	e.members.GetMemberRoleFn = func(_ context.Context, _, userID uuid.UUID) (string, error) {
		if userID == user1 {
			return chat.RoleOwner, nil
		}
		return chat.RoleMember, nil
	}
	// No block between requester and target
	// But block between target (user3) and existing member (user2)
	e.mods.IsBlockedFn = func(_ context.Context, blockerID, blockedID uuid.UUID) (bool, error) {
		if (blockerID == user3 && blockedID == user2) || (blockerID == user2 && blockedID == user3) {
			return true, nil
		}
		return false, nil
	}
	e.members.GetChatMembersFn = func(_ context.Context, _ uuid.UUID) ([]uuid.UUID, error) {
		return []uuid.UUID{user1, user2}, nil
	}

	err := e.svc.AddMember(context.Background(), chatID, user1, user3)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "block relationship")
	assert.Equal(t, 0, e.members.AddMemberCalls)
}

// ═════════════════════════════════════════════════════════════
//  Block Enforcement — Forward to blocked direct
// ═════════════════════════════════════════════════════════════

func TestForwardMessage_BlockedInTargetDirect(t *testing.T) {
	e := setup(t)
	fromChat := uuid.New()
	toChat := uuid.New()

	e.chats.GetChatTypeFn = func(_ context.Context, id uuid.UUID) (string, error) {
		if id == toChat {
			return chat.TypeDirect, nil
		}
		return chat.TypeGroup, nil
	}
	e.members.GetChatMembersFn = func(_ context.Context, cid uuid.UUID) ([]uuid.UUID, error) {
		if cid == toChat {
			return []uuid.UUID{user1, user2}, nil
		}
		return []uuid.UUID{user1}, nil
	}
	// Блок именно между user1 и user2
	e.mods.IsBlockedFn = func(_ context.Context, blockerID, blockedID uuid.UUID) (bool, error) {
		if blockerID == user2 && blockedID == user1 {
			return true, nil
		}
		return false, nil
	}
	e.msgs.GetMessageByIDFn = func(_ context.Context, _ int64) (*chat.Message, error) {
		return &chat.Message{ID: 1, ChatID: fromChat, Content: "Test", Type: chat.MsgTypeText}, nil
	}

	msg, err := e.svc.ForwardMessage(context.Background(), user1, &chat.ForwardMessageReq{
		FromChatID: fromChat, ToChatID: toChat, MessageID: 1,
	})

	assert.Error(t, err)
	assert.Nil(t, msg)
	assert.Contains(t, err.Error(), "blocked")
}
// ═════════════════════════════════════════════════════════════
//  Block Enforcement — SetContactName
// ═════════════════════════════════════════════════════════════

func TestSetContactName_BlockedByTarget(t *testing.T) {
	e := setup(t)
	// user2 blocked user1
	e.mods.IsBlockedFn = func(_ context.Context, blockerID, blockedID uuid.UUID) (bool, error) {
		if blockerID == user2 && blockedID == user1 {
			return true, nil
		}
		return false, nil
	}

	err := e.svc.SetContactName(context.Background(), user1, user2, "Nickname")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "blocked")
}