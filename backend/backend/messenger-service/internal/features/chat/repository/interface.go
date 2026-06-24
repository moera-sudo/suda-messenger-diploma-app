package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	appcrypto "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/crypto"
)

type ChatRepo interface {
	CreateChat(ctx context.Context, c *chat.Chat, memberIDs []uuid.UUID, creatorRole string) error
	GetChatByID(ctx context.Context, chatID uuid.UUID) (*chat.Chat, error)
	GetChatType(ctx context.Context, chatID uuid.UUID) (string, error)
	GetDirectChat(ctx context.Context, userA, userB uuid.UUID) (*chat.Chat, error)
	GetSavedChat(ctx context.Context, userID uuid.UUID) (*chat.Chat, error)
	GetUserChats(ctx context.Context, userID uuid.UUID) ([]chat.ChatListItemResponse, error)
	UpdateChat(ctx context.Context, chatID uuid.UUID, req *chat.UpdateChatReq) error
	DeleteChat(ctx context.Context, chatID uuid.UUID) error
	IsUsernameTaken(ctx context.Context, username string) (bool, error)
}

type MessageRepo interface {
	SaveMessage(ctx context.Context, msg *chat.Message) error
	GetMessageByID(ctx context.Context, messageID int64) (*chat.Message, error)
	GetMessagePreview(ctx context.Context, messageID int64) (*chat.MessagePreview, error)
	GetMessages(ctx context.Context, chatID, userID uuid.UUID, limit, offset int) ([]chat.Message, error)
	EditMessage(ctx context.Context, messageID int64, content string) error
	SoftDeleteMessage(ctx context.Context, messageID int64) error
	HideMessageForUser(ctx context.Context, userID uuid.UUID, messageID int64) error
	GetChatMedia(ctx context.Context, chatID uuid.UUID) (*chat.ChatMediaResponse, error)
}

type MemberRepo interface {
	GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error)
	GetChatMembersDetailed(ctx context.Context, chatID uuid.UUID) ([]chat.ChatMemberInfoResponse, error)
	GetMemberRole(ctx context.Context, chatID, userID uuid.UUID) (string, error)
	IsMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error)
	IsMemberByMessageID(ctx context.Context, userID uuid.UUID, messageID int64) (bool, error)
	AddMember(ctx context.Context, chatID, userID uuid.UUID, role string) error
	RemoveMember(ctx context.Context, chatID, userID uuid.UUID) error
	UpdateMemberRole(ctx context.Context, chatID, userID uuid.UUID, role string) error
	UpdateReadCursor(ctx context.Context, userID, chatID uuid.UUID, lastReadID int64) error
	GetMemberCount(ctx context.Context, chatID uuid.UUID) (int, error)
	GetMessageReaders(ctx context.Context, chatID uuid.UUID, messageID int64) ([]chat.MessageReaderInfoResponse, error)
}

type ModerationRepo interface {
	MuteChat(ctx context.Context, userID, chatID uuid.UUID, mutedUntil *time.Time) error
	UnmuteChat(ctx context.Context, userID, chatID uuid.UUID) error
	IsChatMuted(ctx context.Context, userID, chatID uuid.UUID) (bool, error)

	// ClearAndHideChat implements "delete chat for me": clears the member's
	// visible history and hides the chat from their list (reappears on new message).
	ClearAndHideChat(ctx context.Context, userID, chatID uuid.UUID) error

	BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error
	IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) (bool, error)
	GetBlockedUsers(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
	GetBlockedUsersDetailed(ctx context.Context, userID uuid.UUID) ([]chat.BlockedUserInfo, error)

	PinMessage(ctx context.Context, chatID uuid.UUID, messageID int64, pinnedBy uuid.UUID) error
	UnpinMessage(ctx context.Context, chatID uuid.UUID, messageID int64) error
	GetPinnedMessages(ctx context.Context, chatID uuid.UUID) ([]chat.PinnedMessageInfoResponse, error)

	SetContactName(ctx context.Context, ownerID, contactID uuid.UUID, name string) error
	RemoveContactName(ctx context.Context, ownerID, contactID uuid.UUID) error
	GetContactName(ctx context.Context, ownerID, contactID uuid.UUID) (string, error)
	GetContacts(ctx context.Context, ownerID uuid.UUID) ([]chat.Contact, error)
}

// * Constuctors

func NewChatRepo(db *pgxpool.Pool, cipher appcrypto.ContentCipher) ChatRepo {
	return &chatRepo{db: db, cipher: cipher}
}
func NewMessageRepo(db *pgxpool.Pool, cipher appcrypto.ContentCipher) MessageRepo {
	return &messageRepo{db: db, cipher: cipher}
}
func NewMemberRepo(db *pgxpool.Pool) MemberRepo { return &memberRepo{db: db} }
func NewModerationRepo(db *pgxpool.Pool, cipher appcrypto.ContentCipher) ModerationRepo {
	return &moderationRepo{db: db, cipher: cipher}
}
