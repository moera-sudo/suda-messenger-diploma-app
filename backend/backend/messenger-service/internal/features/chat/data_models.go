package chat

import (
	"time"

	"github.com/google/uuid"
)

// * Chat types
const (
	TypeDirect  = "DIRECT"
	TypeGroup   = "GROUP"
	TypeChannel = "CHANNEL"
	TypeSaved   = "SAVED"
)

// * Chat Member Roles
const (
	RoleOwner  = "OWNER"
	RoleAdmin  = "ADMIN"
	RoleMember = "MEMBER"
)

// * Message types
const (
	MsgTypeText         = "TEXT"
	MsgTypeImage        = "IMAGE"
	MsgTypeVideo        = "VIDEO"
	MsgTypeAudio        = "AUDIO"
	MsgTypeFile         = "FILE"
	MsgTypeVoice        = "VOICE"
	MsgTypeSystem       = "SYSTEM"
	MsgTypeSudaTransfer = "SUDA_TRANSFER"
	MsgTypeDonation     = "DONATION"
)

//* Message Statuses

const (
	StatusSent = "SENT"
	StatusRead = "READ"
)

//* System Message Actions

const (
	SystemUserJoined      = "USER_JOINED"
	SystemUserLeft        = "USER_LEFT"
	SystemUserRemoved     = "USER_REMOVED"
	SystemGroupCreated    = "GROUP_CREATED"
	SystemGroupRenamed    = "GROUP_RENAMED"
	SystemAvatarChanged   = "AVATAR_CHANGED"
	SystemMessagePinned   = "MESSAGE_PINNED"
	SystemMessageUnpinned = "MESSAGE_UNPINNED"
)

//* Chat&Message Data Models

type Chat struct {
	ID              uuid.UUID  `json:"id"              db:"id"`
	Type            string     `json:"type"            db:"type"`
	Name            string     `json:"name,omitempty"  db:"name"`
	Description     *string    `json:"description,omitempty" db:"description"`
	AvatarMediaID   *uuid.UUID `json:"avatar_media_id,omitempty" db:"avatar_media_id"`
	CreatorID       *uuid.UUID `json:"creator_id,omitempty" db:"creator_id"`
	InviteLink      *string    `json:"invite_link,omitempty" db:"invite_link"`
	CreatedAt       time.Time  `json:"created_at"      db:"created_at"`
	Visibility      string     `json:"visibility,omitempty"        db:"visibility"`
	SubscriberCount int        `json:"subscriber_count,omitempty"  db:"subscriber_count"`
	CommentsEnabled bool       `json:"comments_enabled,omitempty"  db:"comments_enabled"`
	Username        *string    `json:"username,omitempty"          db:"username"`
}

type ChatMember struct {
	ChatID               uuid.UUID  `json:"chat_id"               db:"chat_id"`
	UserID               uuid.UUID  `json:"user_id"               db:"user_id"`
	Role                 string     `json:"role"                  db:"role"`
	JoinedAt             time.Time  `json:"joined_at"             db:"joined_at"`
	LastReadMessageID    int64      `json:"last_read_message_id"  db:"last_read_message_id"`
	IsMuted              bool       `json:"is_muted"              db:"is_muted"`
	MutedUntil           *time.Time `json:"muted_until,omitempty" db:"muted_until"`
	NotificationsEnabled bool       `json:"notifications_enabled" db:"notifications_enabled"`
}

type Message struct {
	ID                int64      `json:"id"                            db:"id"`
	ChatID            uuid.UUID  `json:"chat_id"                      db:"chat_id"`
	SenderID          *uuid.UUID `json:"sender_id,omitempty"          db:"sender_id"`
	Content           string     `json:"content"                      db:"content"`
	Type              string     `json:"type"                         db:"type"`
	Status            string     `json:"status"                       db:"status"`
	AttachmentMediaID *uuid.UUID `json:"attachment_media_id,omitempty" db:"attachment_media_id"`
	ReplyToMessageID  *int64     `json:"reply_to_message_id,omitempty" db:"reply_to_message_id"`
	ForwardedFromChat *uuid.UUID `json:"forwarded_from_chat,omitempty" db:"forwarded_from_chat"`
	ForwardedFromMsg  *int64     `json:"forwarded_from_msg,omitempty"  db:"forwarded_from_msg"`
	EditedAt          *time.Time `json:"edited_at,omitempty"          db:"edited_at"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty"         db:"deleted_at"`
	CreatedAt         time.Time  `json:"created_at"                   db:"created_at"`
	ClientSideID      string     `json:"client_side_id,omitempty"     db:"-"`
	// ChatType is set by the service before SaveMessage so the repo can decide
	// whether to encrypt content (DIRECT/GROUP only). Not persisted.
	ChatType string `json:"-" db:"-"`

	// Hydrated previews (not stored — assembled on read / before broadcast) so
	// the client can render the reply / forwarded-from block without having the
	// referenced message in its local cache.
	ReplyPreview  *MessagePreview `json:"reply_preview,omitempty"  db:"-"`
	ForwardedFrom *MessagePreview `json:"forwarded_from,omitempty" db:"-"`
}

// MessagePreview — короткая выжимка о связанном сообщении (reply / forward source).
type MessagePreview struct {
	MessageID  int64      `json:"message_id"`
	SenderID   *uuid.UUID `json:"sender_id,omitempty"`
	SenderName string     `json:"sender_name,omitempty"`
	Content    string     `json:"content"`
	Type       string     `json:"type"`
	Deleted    bool       `json:"deleted,omitempty"`
}

type PinnedMessage struct {
	ChatID    uuid.UUID  `json:"chat_id"    db:"chat_id"`
	MessageID int64      `json:"message_id" db:"message_id"`
	PinnedBy  *uuid.UUID `json:"pinned_by"  db:"pinned_by"`
	PinnedAt  time.Time  `json:"pinned_at"  db:"pinned_at"`
}

type BlockedUser struct {
	BlockerID uuid.UUID `json:"blocker_id" db:"blocker_id"`
	BlockedID uuid.UUID `json:"blocked_id" db:"blocked_id"`
	BlockedAt time.Time `json:"blocked_at" db:"blocked_at"`
}

type Reaction struct {
	ID        uuid.UUID `json:"id"         db:"id"`
	MessageID int64     `json:"message_id" db:"message_id"`
	UserID    uuid.UUID `json:"user_id"    db:"user_id"`
	Emoji     string    `json:"emoji"      db:"emoji"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Contact struct {
	OwnerID    uuid.UUID `json:"owner_id"    db:"owner_id"`
	ContactID  uuid.UUID `json:"contact_id"  db:"contact_id"`
	CustomName string    `json:"custom_name" db:"custom_name"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"  db:"updated_at"`
}
