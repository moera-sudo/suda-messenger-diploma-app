package search

import (
	"time"

	"github.com/google/uuid"
)

const (
	ResultTypeUser    = "USER"
	ResultTypeGroup   = "GROUP"
	ResultTypeChannel = "CHANNEL"
	ResultTypeMessage = "MESSAGE"
)

type SearchResult struct {
	Type        string     `json:"type"`
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Subtitle    string     `json:"subtitle,omitempty"`
	ImageURL    string     `json:"image_url,omitempty"`
	Rank        float64    `json:"rank"`
	ChatID      *uuid.UUID `json:"chat_id,omitempty"`
	MessageID   *int64     `json:"message_id,omitempty"`
	Highlight   string     `json:"highlight,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
}