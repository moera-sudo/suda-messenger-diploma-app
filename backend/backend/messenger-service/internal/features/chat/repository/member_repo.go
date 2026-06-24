package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)

type memberRepo struct {
	db *pgxpool.Pool
}

func (r *memberRepo) GetChatMembers(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx,
		`SELECT user_id FROM messenger_chat_members WHERE chat_id = $1`, chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

func (r *memberRepo) GetChatMembersDetailed(ctx context.Context, chatID uuid.UUID) ([]chat.ChatMemberInfoResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT cm.user_id, u.username, u.display_name, u.avatar_media_id, cm.role, u.last_seen_at
		 FROM messenger_chat_members cm
		 JOIN messenger_users u ON cm.user_id = u.id
		 WHERE cm.chat_id = $1
		 ORDER BY
			CASE cm.role WHEN 'OWNER' THEN 0 WHEN 'ADMIN' THEN 1 ELSE 2 END,
			u.display_name`, chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []chat.ChatMemberInfoResponse
	for rows.Next() {
		var m chat.ChatMemberInfoResponse
		if err := rows.Scan(&m.UserID, &m.Username, &m.DisplayName, &m.AvatarMediaID, &m.Role, &m.LastSeenAt); err != nil {
			return nil, err
		}
		members = append(members, m)
	}
	return members, rows.Err()
}

func (r *memberRepo) GetMemberRole(ctx context.Context, chatID, userID uuid.UUID) (string, error) {
	var role string
	err := r.db.QueryRow(ctx,
		`SELECT role FROM messenger_chat_members WHERE chat_id = $1 AND user_id = $2`,
		chatID, userID,
	).Scan(&role)
	if err != nil {
		return "", fmt.Errorf("member not found: %w", err)
	}
	return role, nil
}

func (r *memberRepo) IsMember(ctx context.Context, chatID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM messenger_chat_members WHERE chat_id = $1 AND user_id = $2)`,
		chatID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *memberRepo) IsMemberByMessageID(ctx context.Context, userID uuid.UUID, messageID int64) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM messenger_messages m
			JOIN messenger_chat_members cm ON cm.chat_id = m.chat_id
			WHERE m.id = $1 AND cm.user_id = $2
		)`, messageID, userID,
	).Scan(&exists)
	return exists, err
}

func (r *memberRepo) AddMember(ctx context.Context, chatID, userID uuid.UUID, role string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO messenger_chat_members (chat_id, user_id, role)
		 VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
		chatID, userID, role,
	)
	return err
}

func (r *memberRepo) RemoveMember(ctx context.Context, chatID, userID uuid.UUID) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM messenger_chat_members WHERE chat_id = $1 AND user_id = $2`,
		chatID, userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("member not found in chat")
	}
	return nil
}

func (r *memberRepo) UpdateMemberRole(ctx context.Context, chatID, userID uuid.UUID, role string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE messenger_chat_members SET role = $1 WHERE chat_id = $2 AND user_id = $3`,
		role, chatID, userID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("member not found in chat")
	}
	return nil
}

func (r *memberRepo) UpdateReadCursor(ctx context.Context, userID, chatID uuid.UUID, lastReadID int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE messenger_chat_members SET last_read_message_id = $1
		 WHERE chat_id = $2 AND user_id = $3 AND last_read_message_id < $1`,
		lastReadID, chatID, userID,
	)
	return err
}

func (r *memberRepo) GetMemberCount(ctx context.Context, chatID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM messenger_chat_members WHERE chat_id = $1`, chatID,
	).Scan(&count)
	return count, err
}

func (r *memberRepo) GetMessageReaders(ctx context.Context, chatID uuid.UUID, messageID int64) ([]chat.MessageReaderInfoResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT cm.user_id, u.username, u.display_name, cm.joined_at
		 FROM messenger_chat_members cm
		 JOIN messenger_users u ON cm.user_id = u.id
		 WHERE cm.chat_id = $1
		   AND cm.last_read_message_id >= $2
		 ORDER BY u.display_name`,
		chatID, messageID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var readers []chat.MessageReaderInfoResponse
	for rows.Next() {
		var r chat.MessageReaderInfoResponse
		if err := rows.Scan(&r.UserID, &r.Username, &r.DisplayName, &r.ReadAt); err != nil {
			return nil, err
		}
		readers = append(readers, r)
	}
	return readers, rows.Err()
}