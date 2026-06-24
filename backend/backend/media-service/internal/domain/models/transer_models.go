package models

import "github.com/google/uuid"

type InitUploadRequest struct {
	OwnerUserID  uuid.UUID `json:"owner_user_id"`
	Kind         string `json:"kind" validate:"required"`
	ContentType  string `json:"content_type" validate:"required"`
	OriginalName string `json:"original_name"`
	IsPrivate    bool   `json:"is_private"`
	SizeBytes    *int64 `json:"size_bytes,omitempty"`
}

type InitUploadResponse struct {
	MediaID   string `json:"media_id"`
	UploadURL string `json:"upload_url"`
	ObjectKey string `json:"object_key"`
	ExpiresIn int    `json:"expires_in"`
}

type ConfirmUploadResponse struct {
	MediaID   string `json:"media_id"`
	Status    string `json:"status"`
	SizeBytes int64  `json:"size_bytes"`
}

type GetURLResponse struct {
	MediaID   string `json:"media_id"`
	URL       string `json:"url"`
	ExpiresIn int    `json:"expires_in"`
}
