package preferences

import "github.com/google/uuid"

// UpdatePreferencesReq — все поля опциональные, обновляем только переданные
type UpdatePreferencesReq struct {
	LangCode *string `json:"lang_code,omitempty" validate:"omitempty,oneof=ru kz en"`
	// Theme is a free-form string. Backend only stores the value and propagates it to other
	// sessions via PREFERENCES_UPDATED; client decides how to render. Lifted from `oneof` to
	// allow custom themes ("suda", "teaChatsDark", ...) without backend changes.
	Theme *string `json:"theme,omitempty" validate:"omitempty,max=64"`

	LastSeenVisibility *string `json:"last_seen_visibility,omitempty" validate:"omitempty,oneof=EVERYONE CONTACTS NOBODY"`
	AvatarVisibility   *string `json:"avatar_visibility,omitempty"    validate:"omitempty,oneof=EVERYONE CONTACTS NOBODY"`
	WhoCanMessage      *string `json:"who_can_message,omitempty"      validate:"omitempty,oneof=EVERYONE CONTACTS NOBODY"`
	WhoCanAddToGroups  *string `json:"who_can_add_to_groups,omitempty" validate:"omitempty,oneof=EVERYONE CONTACTS NOBODY"` // NEW
	AcceptGiftsFrom    *string `json:"accept_gifts_from,omitempty"    validate:"omitempty,oneof=EVERYONE CONTACTS NOBODY"`  // NEW

	NotificationMessages   *bool `json:"notification_messages,omitempty"`
	NotificationGroups     *bool `json:"notification_groups,omitempty"`
	NotificationChannels   *bool `json:"notification_channels,omitempty"`
	NotificationSound      *bool `json:"notification_sound,omitempty"`
	NotificationVibration  *bool `json:"notification_vibration,omitempty"`
	PreviewInNotifications *bool `json:"preview_in_notifications,omitempty"` // NEW

	ChatFontSize      *string    `json:"chat_font_size,omitempty" validate:"omitempty,oneof=small medium large"`
	SendByEnter       *bool      `json:"send_by_enter,omitempty"`
	AutoDownloadMedia *bool      `json:"auto_download_media,omitempty"`
	WallpaperID       *uuid.UUID `json:"wallpaper_id,omitempty"` // NEW

	NFTShowcaseEnabled *bool `json:"nft_showcase_enabled,omitempty"` // NEW
}

type HideLastSeenReq struct {
	UserID uuid.UUID `json:"user_id" validate:"required"`
}
