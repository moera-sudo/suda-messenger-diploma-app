package user

import (
	"time"

	"github.com/google/uuid"
)

type UpdateProfileReq struct {
	DisplayName string `json:"display_name" validate:"max=50"`
	Bio         string `json:"bio" validate:"max=500"`
	FirstName   string `json:"first_name" validate:"max=60"`
	LastName    string `json:"last_name" validate:"max=60"`
}

type UpdateAvatarReq struct {
	MediaID uuid.UUID `json:"media_id" validate:"required"`
}

type LogoutReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	FCMToken     string `json:"fcm_token"`
}

type RegisterDeviceReq struct {
	FCMToken   string `json:"fcm_token" validate:"required"`
	DeviceName string `json:"device_name"`
}

// ChangePasswordReq is the body of PUT /user/password (authenticated change).
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

// SessionInfo describes a single active refresh session (returned by GET /user/sessions).
//
// IsCurrent is true for the session that issued the access token used in this request:
// the access token carries a "sid" claim, the gateway forwards it as X-Session-ID, and
// the handler matches it against each session's ID.
type SessionInfo struct {
	ID         uuid.UUID  `json:"id"`
	UserAgent  string     `json:"user_agent"`
	ClientIP   string     `json:"client_ip"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`
	DeviceName string     `json:"device_name,omitempty"`
	IsCurrent  bool       `json:"is_current"`
}

// TerminateOtherSessionsReq carries the caller's current refresh_token so that the
// backend can keep that one session alive while revoking every other for the user.
type TerminateOtherSessionsReq struct {
	CurrentRefreshToken string `json:"current_refresh_token" validate:"required"`
}

type UserStatusResponse struct {
	UserID     uuid.UUID  `json:"user_id"`
	IsOnline   bool       `json:"is_online"`
	LastSeenAt *time.Time `json:"last_seen_at,omitempty"`
}

type UserResponse struct {
	ID            uuid.UUID  `json:"id"`
	Username      string     `json:"username"`
	DisplayName   string     `json:"display_name"`
	FirstName     string     `json:"first_name"`
	LastName      string     `json:"last_name"`
	Bio           string     `json:"bio"`
	AvatarMediaID *uuid.UUID `json:"avatar_media_id,omitempty"`
	WalletAddress string     `json:"wallet_address"`
	IsVerified    bool       `json:"is_verified"`
}
