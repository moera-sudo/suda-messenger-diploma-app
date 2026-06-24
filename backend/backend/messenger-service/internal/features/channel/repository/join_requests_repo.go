package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel"
)

type joinRequestRepo struct {
	db *pgxpool.Pool
}

func NewJoinRequestRepo(db *pgxpool.Pool) JoinRequestRepo { return &joinRequestRepo{db: db} }

// Upsert creates (or re-opens) a PENDING row for (channel, user). A previously
// rejected request/invite is reset to PENDING.
func (r *joinRequestRepo) Upsert(ctx context.Context, channelID, userID uuid.UUID, kind string, createdBy uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO messenger_channel_join_requests (channel_id, user_id, kind, status, created_by)
		 VALUES ($1, $2, $3, 'PENDING', $4)
		 ON CONFLICT (channel_id, user_id) DO UPDATE
		    SET kind = EXCLUDED.kind,
		        status = 'PENDING',
		        created_by = EXCLUDED.created_by,
		        created_at = NOW(),
		        decided_by = NULL,
		        decided_at = NULL`,
		channelID, userID, kind, createdBy,
	)
	return err
}

// Get returns the row for (channel, user), or (nil, nil) when none exists.
func (r *joinRequestRepo) Get(ctx context.Context, channelID, userID uuid.UUID) (*channel.JoinRequest, error) {
	var jr channel.JoinRequest
	err := r.db.QueryRow(ctx,
		`SELECT channel_id, user_id, kind, status, created_by, created_at
		 FROM messenger_channel_join_requests
		 WHERE channel_id = $1 AND user_id = $2`,
		channelID, userID,
	).Scan(&jr.ChannelID, &jr.UserID, &jr.Kind, &jr.Status, &jr.CreatedBy, &jr.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get join request: %w", err)
	}
	return &jr, nil
}

func (r *joinRequestRepo) SetStatus(ctx context.Context, channelID, userID uuid.UUID, status string, decidedBy uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE messenger_channel_join_requests
		 SET status = $1, decided_by = $2, decided_at = NOW()
		 WHERE channel_id = $3 AND user_id = $4`,
		status, decidedBy, channelID, userID,
	)
	return err
}

func (r *joinRequestRepo) Delete(ctx context.Context, channelID, userID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM messenger_channel_join_requests WHERE channel_id = $1 AND user_id = $2`,
		channelID, userID,
	)
	return err
}

// ListPending returns enriched pending rows of a given kind for a channel.
func (r *joinRequestRepo) ListPending(ctx context.Context, channelID uuid.UUID, kind string, limit, offset int) ([]channel.PendingRequestInfo, error) {
	rows, err := r.db.Query(ctx,
		`SELECT jr.user_id, u.username, u.display_name, u.avatar_media_id, jr.created_at
		 FROM messenger_channel_join_requests jr
		 JOIN messenger_users u ON u.id = jr.user_id
		 WHERE jr.channel_id = $1 AND jr.kind = $2 AND jr.status = 'PENDING'
		 ORDER BY jr.created_at DESC
		 LIMIT $3 OFFSET $4`,
		channelID, kind, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("list pending: %w", err)
	}
	defer rows.Close()
	out := make([]channel.PendingRequestInfo, 0)
	for rows.Next() {
		var p channel.PendingRequestInfo
		if err := rows.Scan(&p.UserID, &p.Username, &p.DisplayName, &p.AvatarMediaID, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan pending: %w", err)
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ListMyInvites returns channels that have a PENDING INVITE for userID.
func (r *joinRequestRepo) ListMyInvites(ctx context.Context, userID uuid.UUID) ([]channel.MyInviteInfo, error) {
	rows, err := r.db.Query(ctx,
		`SELECT c.id, c.name, c.username, c.avatar_media_id, COALESCE(c.description, ''), c.subscriber_count, jr.created_at
		 FROM messenger_channel_join_requests jr
		 JOIN messenger_chats c ON c.id = jr.channel_id
		 WHERE jr.user_id = $1 AND jr.kind = 'INVITE' AND jr.status = 'PENDING'
		 ORDER BY jr.created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("list my invites: %w", err)
	}
	defer rows.Close()
	out := make([]channel.MyInviteInfo, 0)
	for rows.Next() {
		var i channel.MyInviteInfo
		if err := rows.Scan(&i.ChannelID, &i.Name, &i.Username, &i.AvatarMediaID, &i.Description, &i.SubscriberCount, &i.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan invite: %w", err)
		}
		out = append(out, i)
	}
	return out, rows.Err()
}

// ListChannelAdminIDs returns OWNER/ADMIN user ids of a channel (for notifications).
func (r *joinRequestRepo) ListChannelAdminIDs(ctx context.Context, channelID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx,
		`SELECT user_id FROM messenger_chat_members
		 WHERE chat_id = $1 AND role IN ('OWNER', 'ADMIN')`,
		channelID,
	)
	if err != nil {
		return nil, fmt.Errorf("list admins: %w", err)
	}
	defer rows.Close()
	out := make([]uuid.UUID, 0)
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		out = append(out, id)
	}
	return out, rows.Err()
}
