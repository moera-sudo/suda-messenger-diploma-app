package friends

import "github.com/google/uuid"

// SendRequestReq — POST /users/friends/requests body.
type SendRequestReq struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}

// FriendInfo — enriched friend row returned by GET /users/friends.
// Includes minimal user info so client can render the list without a second roundtrip.
type FriendInfo struct {
	UserID        uuid.UUID  `json:"user_id"`
	Username      string     `json:"username"`
	DisplayName   string     `json:"display_name"`
	AvatarMediaID *uuid.UUID `json:"avatar_media_id,omitempty"`
	IsOnline      bool       `json:"is_online"`
	BecameAt      string     `json:"became_at"` // ISO time when request was accepted
}

// FriendRequestInfo — enriched request row for GET /users/friends/requests.
type FriendRequestInfo struct {
	RequestID       uuid.UUID  `json:"request_id"`
	Direction       string     `json:"direction"`         // "incoming" | "outgoing"
	Status          string     `json:"status"`            // PENDING / ACCEPTED / REJECTED / CANCELLED
	OtherUserID     uuid.UUID  `json:"other_user_id"`
	OtherUsername   string     `json:"other_username"`
	OtherDisplay    string     `json:"other_display_name"`
	OtherAvatarID   *uuid.UUID `json:"other_avatar_media_id,omitempty"`
	CreatedAt       string     `json:"created_at"`
}

// RelationStatusResponse — GET /users/friends/{user_id}/status.
type RelationStatusResponse struct {
	UserID    uuid.UUID `json:"user_id"`
	Relation  string    `json:"relation"`            // see RelationXxx constants
	RequestID *uuid.UUID `json:"request_id,omitempty"`
}

// WSPayload — payload of FRIEND_* WebSocket events.
type WSPayload struct {
	RequestID     uuid.UUID  `json:"request_id"`
	FromUserID    uuid.UUID  `json:"from_user_id"`
	FromUsername  string     `json:"from_username"`
	FromDisplay   string     `json:"from_display_name"`
	FromAvatarID  *uuid.UUID `json:"from_avatar_media_id,omitempty"`
}
