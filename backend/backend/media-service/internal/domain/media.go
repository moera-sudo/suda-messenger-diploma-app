package domain

import (
	"time"
	"github.com/google/uuid"
)

const (
	StatusPending = "PENDING"
	StatusReady   = "READY"
	StatusDeleted = "DELETED"
)

// Media kinds
const (
	KindAvatar       = "AVATAR"
	KindChatAvatar   = "CHAT_AVATAR"
	KindAttachment   = "ATTACHMENT"
	KindPostImage    = "POST_IMAGE"
	KindVoiceMessage = "VOICE_MESSAGE"
	KindVideoNote    = "VIDEO_NOTE"
	KindDocument     = "DOCUMENT"
)

// AllowedContentTypes — MIME prefix whitelist per kind
var AllowedContentTypes = map[string][]string{
	KindAvatar:       {"image/"},
	KindChatAvatar:   {"image/"},
	KindPostImage:    {"image/"},
	KindVoiceMessage: {"audio/"},
	KindVideoNote:    {"video/"},
	KindAttachment:   {}, // any type allowed
	KindDocument:     {}, // any type allowed
}

type Media struct {
	ID             string     `json:"id"              db:"id"`
	OwnerUserID    uuid.UUID     `json:"owner_user_id"   db:"owner_user_id"`
	Bucket         string     `json:"bucket"          db:"bucket"`
	ObjectKey      string     `json:"object_key"      db:"object_key"`
	OriginalName   *string    `json:"original_name"   db:"original_name"`
	Kind           string     `json:"kind"            db:"kind"`
	ContentType    string     `json:"content_type"    db:"content_type"`
	SizeBytes      *int64     `json:"size_bytes"      db:"size_bytes"`
	ChecksumSHA256 *string    `json:"checksum_sha256" db:"checksum_sha256"`
	Status         string     `json:"status"          db:"status"`
	IsPrivate      bool       `json:"is_private"      db:"is_private"`
	CreatedAt      time.Time  `json:"created_at"      db:"created_at"`
	ReadyAt        *time.Time `json:"ready_at"        db:"ready_at"`
	DeletedAt      *time.Time `json:"deleted_at"      db:"deleted_at"`
}

func (m *Media) IsReady() bool   { return m.Status == StatusReady }
func (m *Media) IsDeleted() bool { return m.Status == StatusDeleted }
