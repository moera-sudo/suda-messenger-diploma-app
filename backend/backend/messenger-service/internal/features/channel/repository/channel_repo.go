package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel"
)

type channelRepo struct {
	db *pgxpool.Pool
}

func (r *channelRepo) GetByUsername(ctx context.Context, username string) (uuid.UUID, error) {
	var id uuid.UUID
	err := r.db.QueryRow(ctx,
		`SELECT id FROM messenger_chats
		 WHERE username = $1 AND type = 'CHANNEL'`,
		username,
	).Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) {
		return uuid.Nil, fmt.Errorf("channel not found")
	}
	return id, err
}

func (r *channelRepo) UpdateSettings(ctx context.Context, channelID uuid.UUID, req *channel.UpdateChannelSettingsReq) error {
	var sets []string
	var args []interface{}
	argNum := 1

	if req.CommentsEnabled != nil {
		sets = append(sets, fmt.Sprintf("comments_enabled = $%d", argNum))
		args = append(args, *req.CommentsEnabled)
		argNum++
	}
	if req.Visibility != nil {
		sets = append(sets, fmt.Sprintf("visibility = $%d", argNum))
		args = append(args, *req.Visibility)
		argNum++
	}
	if req.Username != nil {
		sets = append(sets, fmt.Sprintf("username = $%d", argNum))
		args = append(args, *req.Username)
		argNum++
	}

	if len(sets) == 0 {
		return nil
	}

	args = append(args, channelID)
	query := fmt.Sprintf(
		`UPDATE messenger_chats SET %s WHERE id = $%d AND type = 'CHANNEL'`,
		strings.Join(sets, ", "), argNum,
	)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *channelRepo) IsUsernameTaken(ctx context.Context, username string, excludeChatID *uuid.UUID) (bool, error) {
	var exists bool
	if excludeChatID != nil {
		err := r.db.QueryRow(ctx,
			`SELECT EXISTS (SELECT 1 FROM messenger_chats WHERE username = $1 AND id != $2)`,
			username, *excludeChatID,
		).Scan(&exists)
		return exists, err
	}

	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM messenger_chats WHERE username = $1)`,
		username,
	).Scan(&exists)
	return exists, err
}

func (r *channelRepo) GetSubscriberCount(ctx context.Context, channelID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT subscriber_count FROM messenger_chats WHERE id = $1`,
		channelID,
	).Scan(&count)
	return count, err
}

func (r *channelRepo) GetSubscribers(ctx context.Context, channelID uuid.UUID, limit, offset int) ([]channel.SubscriberInfo, error) {
	rows, err := r.db.Query(ctx,
		`SELECT cm.user_id, u.username, u.display_name, cm.role, cm.joined_at
		 FROM messenger_chat_members cm
		 JOIN messenger_users u ON cm.user_id = u.id
		 WHERE cm.chat_id = $1
		 ORDER BY
		     CASE cm.role
		         WHEN 'OWNER' THEN 1
		         WHEN 'ADMIN' THEN 2
		         ELSE 3
		     END,
		     cm.joined_at ASC
		 LIMIT $2 OFFSET $3`,
		channelID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []channel.SubscriberInfo
	for rows.Next() {
		var s channel.SubscriberInfo
		if err := rows.Scan(&s.UserID, &s.Username, &s.DisplayName, &s.Role, &s.JoinedAt); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}
	return subs, rows.Err()
}