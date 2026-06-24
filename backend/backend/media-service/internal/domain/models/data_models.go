package models

import (
	"time"

	"github.com/google/uuid"
)

type MediaMetadata struct {
	ID           string     `json:"id"`
	OwnerUserID  uuid.UUID  `json:"owner_user_id"`
	Kind         string     `json:"kind"`
	ContentType  string     `json:"content_type"`
	SizeBytes    *int64     `json:"size_bytes"`
	OriginalName *string    `json:"original_name"`
	Status       string     `json:"status"`
	IsPrivate    bool       `json:"is_private"`
	CreatedAt    time.Time  `json:"created_at"`
	ReadyAt      *time.Time `json:"ready_at"`
}
