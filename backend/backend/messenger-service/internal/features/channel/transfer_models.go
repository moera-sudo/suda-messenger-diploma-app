package channel

import "github.com/google/uuid"

type CreateChannelReq struct {
	Name        string  `json:"name"        validate:"required,min=1,max=100"`
	Description string  `json:"description" validate:"max=500"`
	Visibility  string  `json:"visibility"  validate:"omitempty,oneof=PUBLIC PRIVATE"`
	Username    string  `json:"username"    validate:"omitempty,min=3,max=32,alphanum"`
}

type UpdateChannelSettingsReq struct {
	CommentsEnabled *bool   `json:"comments_enabled,omitempty"`
	Visibility      *string `json:"visibility,omitempty"       validate:"omitempty,oneof=PUBLIC PRIVATE"`
	Username        *string `json:"username,omitempty"         validate:"omitempty,min=3,max=32,alphanum"`
}

type AddCommentReq struct {
	Content          string `json:"content"                       validate:"required,min=1,max=4000"`
	ReplyToCommentID *int64 `json:"reply_to_comment_id,omitempty"`
}

type EditCommentReq struct {
	Content string `json:"content" validate:"required,min=1,max=4000"`
}

type LinkAppReq struct {
	AppID uuid.UUID `json:"app_id" validate:"required"`
}

type ChannelInfoResponse struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Description     string     `json:"description,omitempty"`
	Username        *string    `json:"username,omitempty"`
	AvatarMediaID   *uuid.UUID `json:"avatar_media_id,omitempty"`
	Visibility      string     `json:"visibility"`
	SubscriberCount int        `json:"subscriber_count"`
	CommentsEnabled bool       `json:"comments_enabled"`
	IsSubscribed    bool       `json:"is_subscribed"`
	UserRole        string     `json:"user_role,omitempty"` // OWNER/ADMIN/SUBSCRIBER
	CreatedAt       string     `json:"created_at"`
}

type CommentsListResponse struct {
	Comments []ChannelComment `json:"comments"`
	Total    int              `json:"total"`
}

// InviteUserReq — admin invites a user to a private channel. Either user_id or
// username (handle without @) must be provided; the handler resolves username.
type InviteUserReq struct {
	UserID   *uuid.UUID `json:"user_id,omitempty"`
	Username string     `json:"username,omitempty"`
}

// ChannelSettingsResponse — current editable settings, for prefilling the admin
// Settings form (GET /channels/{id}/settings).
type ChannelSettingsResponse struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Description     string     `json:"description,omitempty"`
	Username        *string    `json:"username,omitempty"`
	AvatarMediaID   *uuid.UUID `json:"avatar_media_id,omitempty"`
	Visibility      string     `json:"visibility"`
	CommentsEnabled bool       `json:"comments_enabled"`
}

// ChannelView — channel card + the viewer's relationship to it. Returned by
// GET /channels/{id}/view and GET /channels/by-username/{username}. Lets the
// client choose the CTA (Open / Subscribe / Request to join / Accept invite)
// and know whether posts are readable.
type ChannelView struct {
	ID              uuid.UUID  `json:"id"`
	Name            string     `json:"name"`
	Description     string     `json:"description,omitempty"`
	Username        *string    `json:"username,omitempty"`
	AvatarMediaID   *uuid.UUID `json:"avatar_media_id,omitempty"`
	Visibility      string     `json:"visibility"`
	SubscriberCount int        `json:"subscriber_count"`
	CommentsEnabled bool       `json:"comments_enabled"`
	IsMember        bool       `json:"is_member"`
	MyRole          string     `json:"my_role,omitempty"`
	// JoinStatus: "none" | "pending_request" | "pending_invite" | "member".
	JoinStatus string `json:"join_status"`
	// CanRead: true if the viewer may read posts (public channel, or member).
	CanRead   bool   `json:"can_read"`
	CreatedAt string `json:"created_at"`
}