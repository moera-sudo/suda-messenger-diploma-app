package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel"
)

type ChannelRepo interface {
	GetByUsername(ctx context.Context, username string) (uuid.UUID, error)
	UpdateSettings(ctx context.Context, channelID uuid.UUID, req *channel.UpdateChannelSettingsReq) error
	IsUsernameTaken(ctx context.Context, username string, excludeChatID *uuid.UUID) (bool, error)
	GetSubscriberCount(ctx context.Context, channelID uuid.UUID) (int, error)
	GetSubscribers(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]channel.SubscriberInfo, error)
}

type CommentsRepo interface {
	Create(ctx context.Context, c *channel.ChannelComment) error
	GetByID(ctx context.Context, commentID int64) (*channel.ChannelComment, error)
	GetByPost(ctx context.Context, postID int64, limit, offset int) ([]channel.ChannelComment, error)
	Update(ctx context.Context, commentID int64, content string) error
	SoftDelete(ctx context.Context, commentID int64) error
	CountByPost(ctx context.Context, postID int64) (int, error)
}

type AppsRepo interface {
	Link(ctx context.Context, link *channel.ChannelAppLink) error
	Unlink(ctx context.Context, channelID, appID uuid.UUID) error
	GetByChannel(ctx context.Context, channelID uuid.UUID) ([]channel.ChannelAppLink, error)
}

// JoinRequestRepo — join requests (user→channel) and invites (admin→user) for
// private channels. One row per (channel, user); see messenger_channel_join_requests.
type JoinRequestRepo interface {
	Upsert(ctx context.Context, channelID, userID uuid.UUID, kind string, createdBy uuid.UUID) error
	Get(ctx context.Context, channelID, userID uuid.UUID) (*channel.JoinRequest, error)
	SetStatus(ctx context.Context, channelID, userID uuid.UUID, status string, decidedBy uuid.UUID) error
	Delete(ctx context.Context, channelID, userID uuid.UUID) error
	ListPending(ctx context.Context, channelID uuid.UUID, kind string, limit, offset int) ([]channel.PendingRequestInfo, error)
	ListMyInvites(ctx context.Context, userID uuid.UUID) ([]channel.MyInviteInfo, error)
	ListChannelAdminIDs(ctx context.Context, channelID uuid.UUID) ([]uuid.UUID, error)
}

func NewChannelRepo(db *pgxpool.Pool) ChannelRepo { return &channelRepo{db: db} }
func NewCommentsRepo(db *pgxpool.Pool) CommentsRepo { return &commentsRepo{db: db} }
func NewAppsRepo(db *pgxpool.Pool) AppsRepo { return &appsRepo{db: db} }