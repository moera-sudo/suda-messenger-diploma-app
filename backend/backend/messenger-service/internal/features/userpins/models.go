package userpins

import (
	"time"

	"github.com/google/uuid"
)

const (
	PinTypeSidebar  = "SIDEBAR"
	PinTypeChatlist = "CHATLIST"
	PinTypeAppHub   = "APP_HUB"

	TargetTypeChat = "CHAT"
	TargetTypeApp  = "APP"
)

type UserPin struct {
	ID         int64     `json:"id"          db:"id"`
	UserID     uuid.UUID `json:"user_id"     db:"user_id"`
	PinType    string    `json:"pin_type"    db:"pin_type"`
	TargetType string    `json:"target_type" db:"target_type"`
	TargetID   uuid.UUID `json:"target_id"   db:"target_id"`
	SortOrder  int       `json:"sort_order"  db:"sort_order"`
	CreatedAt  time.Time `json:"created_at"  db:"created_at"`
}

// PinnedChat — закреплённый чат для построения списка чатов:
// id чата + время закрепления (для пометки is_pinned и сортировки по дате).
type PinnedChat struct {
	ChatID   uuid.UUID
	PinnedAt time.Time
}

// IsValidCombination проверяет валидность пары pin_type + target_type
func IsValidCombination(pinType, targetType string) bool {
	switch pinType {
	case PinTypeChatlist:
		return targetType == TargetTypeChat
	case PinTypeSidebar:
		return targetType == TargetTypeChat || targetType == TargetTypeApp
	case PinTypeAppHub:
		return targetType == TargetTypeApp
	}
	return false
}