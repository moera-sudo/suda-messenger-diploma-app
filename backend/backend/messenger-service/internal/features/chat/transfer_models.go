package chat

import (
	"time"

	"github.com/google/uuid"
)

// *Chat Requests
type CreateChatReq struct {
	Type        string      `json:"type"        validate:"required,oneof=DIRECT GROUP SAVED CHANNEL"`
	Name        string      `json:"name,omitempty"`
	Description string      `json:"description,omitempty"`
	TargetID    uuid.UUID   `json:"target_id,omitempty"`
	MemberIDs   []uuid.UUID `json:"member_ids,omitempty"`

	// Для каналов
	Visibility string `json:"visibility,omitempty" validate:"omitempty,oneof=PUBLIC PRIVATE"`
	Username   string `json:"username,omitempty"   validate:"omitempty,min=3,max=32,alphanum"`
}

type UpdateChatReq struct {
	Name          *string    `json:"name,omitempty"`
	Description   *string    `json:"description,omitempty"`
	AvatarMediaID *uuid.UUID `json:"avatar_media_id,omitempty"`
}

//* Messages Requests

type SendMessageReq struct {
	ChatID            uuid.UUID  `json:"chat_id" validate:"required"`
	Content           string     `json:"content"`
	Type              string     `json:"type" validate:"required"`
	AttachmentMediaID *uuid.UUID `json:"attachment_media_id,omitempty"`
	ReplyToMessageID  *int64     `json:"reply_to_message_id,omitempty"`
	ClientSideID      string     `json:"client_side_id"`
}

type EditMessageReq struct {
	Content string `json:"content" validate:"required"`
}

type DeleteMessageReq struct {
	ForEveryone bool `json:"for_everyone"`
}

type ForwardMessageReq struct {
	FromChatID uuid.UUID `json:"from_chat_id" validate:"required"`
	ToChatID   uuid.UUID `json:"to_chat_id" validate:"required"`
	MessageID  int64     `json:"message_id" validate:"required"`
}

// * Member Requests
type AddMemberReq struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

type UpdateMemberRoleReq struct {
	Role string `json:"role" validate:"required,oneof=ADMIN MEMBER"`
}

type MuteChatReq struct {
	MutedUntil *time.Time `json:"muted_until,omitempty"` // nil = mute forever
}

// * Pin Requests
type PinMessageReq struct {
	MessageID int64 `json:"message_id" validate:"required"`
}

// * Block Requests
type BlockUserReq struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// * Reaction Requests
type ReactionReq struct {
	Emoji string `json:"emoji" validate:"required"`
}

// * Contacts Requests
type SetContactNameReq struct {
	UserID     uuid.UUID `json:"user_id" validate:"required"`
	CustomName string    `json:"custom_name" validate:"required,max=50"`
}

// * ResponseModels
type ChatListItemResponse struct {
	Chat
	LastMessageContent string     `json:"last_message_content"`
	LastMessageType    string     `json:"last_message_type"`
	LastMessageTime    time.Time  `json:"last_message_time"`
	UnreadCount        int        `json:"unread_count"`
	IsMuted            bool       `json:"is_muted"`
	InterlocutorID     *uuid.UUID `json:"interlocutor_id,omitempty"`
	InterlocutorName   string     `json:"interlocutor_name,omitempty"`
	InterlocutorAvatar string     `json:"interlocutor_avatar,omitempty"`
	// Interlocutor's avatar media id for DIRECT chats — resolve via
	// GET /media/{id}/url (cached). Lets the client render the direct-chat
	// avatar straight from the list, without a per-row GET /users/{id} call.
	InterlocutorAvatarMediaID *uuid.UUID `json:"interlocutor_avatar_media_id,omitempty"`
	MemberCount        int        `json:"member_count,omitempty"`
	OnlineCount        int        `json:"online_count,omitempty"`
	IsPinned           bool       `json:"is_pinned"`
	PinnedAt           *time.Time `json:"pinned_at,omitempty"`
	PinOrder           int        `json:"-"`
	// Block status for DIRECT chats (always false otherwise). The client can
	// disable the composer up-front instead of discovering a block via a 403
	// on send. can_write := !(blocked_by_me || blocked_me).
	BlockedByMe bool `json:"blocked_by_me"`
	BlockedMe   bool `json:"blocked_me"`
}

type ChatInfoResponse struct {
	Chat
	Members        []ChatMemberInfoResponse    `json:"members"`
	MemberCount    int                         `json:"member_count"`
	OnlineCount    int                         `json:"online_count"`
	PinnedMessages []PinnedMessageInfoResponse `json:"pinned_messages,omitempty"`
	// Block status for DIRECT chats (false otherwise). See ChatListItemResponse.
	BlockedByMe bool `json:"blocked_by_me"`
	BlockedMe   bool `json:"blocked_me"`
}

type ChatMemberInfoResponse struct {
	UserID        uuid.UUID  `json:"user_id"`
	Username      string     `json:"username"`
	DisplayName   string     `json:"display_name"`
	AvatarMediaID *uuid.UUID `json:"avatar_media_id,omitempty"`
	Role          string     `json:"role"`
	IsOnline      bool       `json:"is_online"`
	LastSeenAt    *time.Time `json:"last_seen_at,omitempty"`
}

// BlockedUserInfo — minimal public info for a blocked user, returned by
// GET /users/blocked so the client can render names/avatars in the Blocked tab
// (the regular GET /user/id/{id} returns 403 for blocked users).
type BlockedUserInfo struct {
	UserID        uuid.UUID  `json:"user_id"`
	Username      string     `json:"username"`
	DisplayName   string     `json:"display_name"`
	AvatarMediaID *uuid.UUID `json:"avatar_media_id,omitempty"`
}

type PinnedMessageInfoResponse struct {
	PinnedMessage
	Content  string     `json:"content"`
	SenderID *uuid.UUID `json:"sender_id"`
	Type     string     `json:"type"`
}

type ChatMediaResponse struct {
	Images    []MediaItemResponse `json:"images"`
	Videos    []MediaItemResponse `json:"videos"`
	Documents []MediaItemResponse `json:"documents"`
	Audio     []MediaItemResponse `json:"audio"`
	Links     []LinkItemResponse  `json:"links"`
}

type MediaItemResponse struct {
	MessageID    int64      `json:"message_id"`
	MediaID      uuid.UUID  `json:"media_id"`
	ContentType  string     `json:"content_type"`
	OriginalName *string    `json:"original_name,omitempty"`
	SenderID     *uuid.UUID `json:"sender_id"`
	CreatedAt    time.Time  `json:"created_at"`
}

type LinkItemResponse struct {
	MessageID int64      `json:"message_id"`
	URL       string     `json:"url"`
	SenderID  *uuid.UUID `json:"sender_id"`
	CreatedAt time.Time  `json:"created_at"`
}

type MessageReaderInfoResponse struct {
	UserID      uuid.UUID  `json:"user_id"`
	Username    string     `json:"username"`
	DisplayName string     `json:"display_name"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
}
