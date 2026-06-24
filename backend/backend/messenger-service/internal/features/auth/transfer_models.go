package auth

import (
	"time"
)

type RegisterReq struct {
	Username    string `json:"username" validate:"required,min=3,max=30,startswith=@"`
	DisplayName string `json:"name" validate:"required,max=50"`
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=6"`
}

type VerifyReq struct {
	Email string `json:"email" validate:"required,email"`
	Code  string `json:"code" validate:"required,len=6"`
}

type LoginReq struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	UserAgent string `json:"user_agent" validate:"required"`
	ClientIP  string `json:"client_ip" validate:"required"`
}

type RefreshSessionReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
	UserAgent    string `json:"user_agent" validate:"required"`
	ClientIP     string `json:"client_ip" validate:"required"`
}

type ForgotPasswordReq struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordReq struct {
	Email       string `json:"email" validate:"required,email"`
	Code        string `json:"code" validate:"required,len=6"`
	NewPassword string `json:"password" validate:"required,min=6"`
}


type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type TokenResponse struct {
	AccessToken      string        `json:"access_token"`
	RefreshToken     string        `json:"refresh_token"`
	AccessExpiresIn  int64         `json:"access_expires_in"`
	RefreshExpiresIn time.Duration `json:"refresh_expires_in"`
}

// WebappSessionReq — body of POST /auth/webapp-session.
//
// Mini Apps (Wallet, Marketplace) open inside a WebView with `?initData=...`
// in the URL. The SPA POSTs that very string here, the backend validates the
// HMAC, and returns a short-lived access JWT that the SPA uses for ordinary
// API calls. See backend/docs/WebView-Integration-Guide.md.
type WebappSessionReq struct {
	InitData string `json:"init_data" validate:"required"`
}

// WebappSessionResponse — short-lived access token for a WebView mini-app.
// No refresh token: when the access token expires, the WebView simply re-runs
// /webapp-session with the same initData string from its URL (if still
// within the 30-minute validity window) or asks the host app to reopen with
// a fresh initData.
type WebappSessionResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
	UserID      string `json:"user_id"`
	Username    string `json:"username"`
	// AvatarMediaID lets the WebView SPA render the user's avatar via the media
	// service (GET /api/v1/media/{id}/url). Empty when the user has no avatar.
	AvatarMediaID string `json:"avatar_media_id,omitempty"`
}
