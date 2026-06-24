package channel

import (
	"time"

	"github.com/google/uuid"
)

const (
	VisibilityPublic  = "PUBLIC"
	VisibilityPrivate = "PRIVATE"

	RoleSubscriber = "SUBSCRIBER"

	// Join request kinds / statuses (messenger_channel_join_requests).
	JoinKindRequest = "REQUEST" // user -> private channel (admin approves)
	JoinKindInvite  = "INVITE"  // admin -> user (user accepts)

	JoinStatusPending  = "PENDING"
	JoinStatusApproved = "APPROVED"
	JoinStatusRejected = "REJECTED"
)

// GatingError — типизированная ошибка: юзер не прошёл token-gating канала.
// Несёт reason и числа баланса, чтобы handler отдал фронту понятный ответ
// ("вход стоит N SUDA, у вас M").
type GatingError struct {
	Reason         string // wallet_not_found | insufficient_balance | nft_required
	MinBalanceWei  string
	UserBalanceWei string
}

func (e *GatingError) Error() string {
	return "gating requirement not met: " + e.Reason
}

type ChannelComment struct {
	ID                int64      `json:"id"                            db:"id"`
	ChannelID         uuid.UUID  `json:"channel_id"                    db:"channel_id"`
	PostID            int64      `json:"post_id"                       db:"post_id"`
	SenderID          uuid.UUID  `json:"sender_id"                     db:"sender_id"`
	Content           string     `json:"content"                       db:"content"`
	ReplyToCommentID  *int64     `json:"reply_to_comment_id,omitempty" db:"reply_to_comment_id"`
	EditedAt          *time.Time `json:"edited_at,omitempty"           db:"edited_at"`
	DeletedAt         *time.Time `json:"-"                              db:"deleted_at"`
	CreatedAt         time.Time  `json:"created_at"                    db:"created_at"`

	// Дополнительные поля для UI (заполняются в JOIN)
	SenderUsername    string `json:"sender_username,omitempty"    db:"-"`
	SenderDisplayName string `json:"sender_display_name,omitempty" db:"-"`

	// Hydrated preview of the parent comment (when this is a reply), so the
	// client renders the reply block even if the parent is outside the loaded
	// window. nil when not a reply / parent deleted.
	ReplyPreview *CommentReplyPreview `json:"reply_preview,omitempty" db:"-"`
}

// CommentReplyPreview — короткая выжимка о родительском комментарии для рендера reply.
type CommentReplyPreview struct {
	SenderDisplayName string `json:"sender_display_name"`
	Content           string `json:"content"`
}

type ChannelAppLink struct {
	ChannelID uuid.UUID `json:"channel_id" db:"channel_id"`
	AppID     uuid.UUID `json:"app_id"     db:"app_id"`
	AddedBy   uuid.UUID `json:"added_by"   db:"added_by"`
	AddedAt   time.Time `json:"added_at"   db:"added_at"`
}

// SubscriberInfo — для списка подписчиков (только админам)
type SubscriberInfo struct {
	UserID       uuid.UUID `json:"user_id"`
	Username     string    `json:"username"`
	DisplayName  string    `json:"display_name"`
	Role         string    `json:"role"`
	JoinedAt     time.Time `json:"joined_at"`
	IsOnline     bool      `json:"is_online"`
}

// JoinRequest — raw row of messenger_channel_join_requests.
type JoinRequest struct {
	ChannelID uuid.UUID
	UserID    uuid.UUID
	Kind      string
	Status    string
	CreatedBy uuid.UUID
	CreatedAt time.Time
}

// PendingRequestInfo — enriched pending join REQUEST for the admin moderation list.
type PendingRequestInfo struct {
	UserID        uuid.UUID  `json:"user_id"`
	Username      string     `json:"username"`
	DisplayName   string     `json:"display_name"`
	AvatarMediaID *uuid.UUID `json:"avatar_media_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// MyInviteInfo — a pending INVITE shown to the invited user (channel card).
type MyInviteInfo struct {
	ChannelID       uuid.UUID  `json:"channel_id"`
	Name            string     `json:"name"`
	Username        *string    `json:"username,omitempty"`
	AvatarMediaID   *uuid.UUID `json:"avatar_media_id,omitempty"`
	Description     string     `json:"description,omitempty"`
	SubscriberCount int        `json:"subscriber_count"`
	CreatedAt       time.Time  `json:"created_at"`
}