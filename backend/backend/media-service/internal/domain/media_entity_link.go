package domain

import "time"

const (
	EntityTypeMessage    = "MESSAGE"
	EntityTypeChatAvatar = "CHAT_AVATAR"
	EntityTypeUserAvatar = "USER_AVATAR"
)

type MediaEntityLink struct {
	ID         string    `json:"id"          db:"id"`
	MediaID    string    `json:"media_id"    db:"media_id"`
	EntityType string    `json:"entity_type" db:"entity_type"`
	EntityID   string    `json:"entity_id"   db:"entity_id"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"`
}