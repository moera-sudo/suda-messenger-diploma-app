package preferences

import (
	"time"

	"github.com/google/uuid"
)

// Privacy visibility levels
const (
	VisibilityEveryone = "EVERYONE"
	VisibilityContacts = "CONTACTS"
	VisibilityNobody   = "NOBODY"
)

// Languages
const (
	LangRu = "ru"
	LangKz = "kz"
	LangEn = "en"
)

// Themes
const (
	ThemeLight  = "light"
	ThemeDark   = "dark"
	ThemeSystem = "system"
)

// Font sizes
const (
	FontSmall  = "small"
	FontMedium = "medium"
	FontLarge  = "large"
)

type UserPreferences struct {
	UserID uuid.UUID `json:"user_id" db:"user_id"`

	LangCode string `json:"lang_code" db:"lang_code"`
	Theme    string `json:"theme"     db:"theme"`

	// Privacy
	LastSeenVisibility   string `json:"last_seen_visibility"   db:"last_seen_visibility"`
	AvatarVisibility     string `json:"avatar_visibility"      db:"avatar_visibility"`
	WhoCanMessage        string `json:"who_can_message"        db:"who_can_message"`
	WhoCanAddToGroups    string `json:"who_can_add_to_groups"  db:"who_can_add_to_groups"`     // NEW
	AcceptGiftsFrom      string `json:"accept_gifts_from"      db:"accept_gifts_from"`         // NEW

	// Notifications
	NotificationMessages    bool `json:"notification_messages"    db:"notification_messages"`
	NotificationGroups      bool `json:"notification_groups"      db:"notification_groups"`
	NotificationChannels    bool `json:"notification_channels"    db:"notification_channels"`
	NotificationSound       bool `json:"notification_sound"       db:"notification_sound"`
	NotificationVibration   bool `json:"notification_vibration"   db:"notification_vibration"`
	PreviewInNotifications  bool `json:"preview_in_notifications" db:"preview_in_notifications"` // NEW

	// Chat UI
	ChatFontSize      string     `json:"chat_font_size"      db:"chat_font_size"`
	SendByEnter       bool       `json:"send_by_enter"       db:"send_by_enter"`
	AutoDownloadMedia bool       `json:"auto_download_media" db:"auto_download_media"`
	WallpaperID       *uuid.UUID `json:"wallpaper_id"        db:"wallpaper_id"`        // NEW

	// Blockchain features
	NFTShowcaseEnabled bool `json:"nft_showcase_enabled" db:"nft_showcase_enabled"`     // NEW

	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

func Default(userID uuid.UUID) *UserPreferences {
	return &UserPreferences{
		UserID:                 userID,
		LangCode:               LangEn,
		Theme:                  ThemeSystem,
		LastSeenVisibility:     VisibilityEveryone,
		AvatarVisibility:       VisibilityEveryone,
		WhoCanMessage:          VisibilityEveryone,
		WhoCanAddToGroups:      VisibilityEveryone,
		AcceptGiftsFrom:        VisibilityEveryone,
		NotificationMessages:   true,
		NotificationGroups:     true,
		NotificationChannels:   true,
		NotificationSound:      true,
		NotificationVibration:  true,
		PreviewInNotifications: true,
		ChatFontSize:           FontMedium,
		SendByEnter:            true,
		AutoDownloadMedia:      true,
		NFTShowcaseEnabled:     false,
		WallpaperID:            nil,
	}
}