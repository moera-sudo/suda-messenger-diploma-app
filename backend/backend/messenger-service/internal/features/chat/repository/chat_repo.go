package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	appcrypto "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/crypto"
)

type chatRepo struct {
	db     *pgxpool.Pool
	cipher appcrypto.ContentCipher
}

func (r *chatRepo) CreateChat(ctx context.Context, c *chat.Chat, memberIDs []uuid.UUID, creatorRole string) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer tx.Rollback(ctx)

	var visibilityArg interface{}
	if c.Visibility != "" {
		visibilityArg = c.Visibility
	}
	var commentsArg interface{}
	if c.Type == chat.TypeChannel {
		commentsArg = c.CommentsEnabled
	}
	err = tx.QueryRow(ctx,
		`INSERT INTO messenger_chats
		     (id, type, name, description, avatar_media_id, creator_id,
		      visibility, comments_enabled, username)
		 VALUES ($1, $2, $3, $4, $5, $6,
		         COALESCE($7, 'PRIVATE'), COALESCE($8, TRUE), $9)
		 RETURNING created_at`,
		c.ID, c.Type, c.Name, c.Description, c.AvatarMediaID, c.CreatorID,
		visibilityArg, commentsArg, c.Username,
	).Scan(&c.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert chat: %w", err)
	}

	for i, uid := range memberIDs {
		role := chat.RoleMember
		if i == 0 && creatorRole != "" {
			role = creatorRole
		}
		_, err := tx.Exec(ctx,
			`INSERT INTO messenger_chat_members (chat_id, user_id, role) VALUES ($1, $2, $3)`,
			c.ID, uid, role,
		)
		if err != nil {
			return fmt.Errorf("add member %s: %w", uid, err)
		}
	}

	return tx.Commit(ctx)
}

func (r *chatRepo) GetChatByID(ctx context.Context, chatID uuid.UUID) (*chat.Chat, error) {
	var c chat.Chat
	err := r.db.QueryRow(ctx,
		`SELECT id, type, name, description, avatar_media_id, creator_id, invite_link, created_at,
		        COALESCE(visibility, ''), COALESCE(subscriber_count, 0), COALESCE(comments_enabled, FALSE), username
		 FROM messenger_chats WHERE id = $1`, chatID,
	).Scan(&c.ID, &c.Type, &c.Name, &c.Description, &c.AvatarMediaID, &c.CreatorID, &c.InviteLink, &c.CreatedAt,
		&c.Visibility, &c.SubscriberCount, &c.CommentsEnabled, &c.Username)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &c, err
}

func (r *chatRepo) GetChatType(ctx context.Context, chatID uuid.UUID) (string, error) {
	var chatType string
	err := r.db.QueryRow(ctx,
		`SELECT type FROM messenger_chats WHERE id = $1`, chatID,
	).Scan(&chatType)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", fmt.Errorf("chat not found: %s", chatID)
	}
	return chatType, err
}

func (r *chatRepo) GetDirectChat(ctx context.Context, userA, userB uuid.UUID) (*chat.Chat, error) {
	var c chat.Chat
	err := r.db.QueryRow(ctx,
		`SELECT c.id, c.type, c.name, c.description, c.avatar_media_id, c.creator_id, c.invite_link, c.created_at,
		        COALESCE(c.visibility, ''), COALESCE(c.subscriber_count, 0), COALESCE(c.comments_enabled, FALSE), c.username
		 FROM messenger_chats c
		 JOIN messenger_chat_members m1 ON c.id = m1.chat_id
		 JOIN messenger_chat_members m2 ON c.id = m2.chat_id
		 WHERE c.type = 'DIRECT' AND m1.user_id = $1 AND m2.user_id = $2
		 LIMIT 1`, userA, userB,
	).Scan(&c.ID, &c.Type, &c.Name, &c.Description, &c.AvatarMediaID, &c.CreatorID, &c.InviteLink, &c.CreatedAt,
		&c.Visibility, &c.SubscriberCount, &c.CommentsEnabled, &c.Username)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &c, err
}

func (r *chatRepo) GetSavedChat(ctx context.Context, userID uuid.UUID) (*chat.Chat, error) {
	var c chat.Chat
	err := r.db.QueryRow(ctx,
		`SELECT c.id, c.type, c.name, c.description, c.avatar_media_id, c.creator_id, c.invite_link, c.created_at,
		        COALESCE(c.visibility, ''), COALESCE(c.subscriber_count, 0), COALESCE(c.comments_enabled, FALSE), c.username
		 FROM messenger_chats c
		 JOIN messenger_chat_members m ON c.id = m.chat_id
		 WHERE c.type = 'SAVED' AND m.user_id = $1
		 LIMIT 1`, userID,
	).Scan(&c.ID, &c.Type, &c.Name, &c.Description, &c.AvatarMediaID, &c.CreatorID, &c.InviteLink, &c.CreatedAt,
		&c.Visibility, &c.SubscriberCount, &c.CommentsEnabled, &c.Username)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	return &c, err
}

func (r *chatRepo) GetUserChats(ctx context.Context, userID uuid.UUID) ([]chat.ChatListItemResponse, error) {
	query := `
		SELECT
			c.id, c.type, c.name, c.description, c.avatar_media_id, c.creator_id, c.invite_link, c.created_at,
			COALESCE(lm.content, '') as last_msg,
			COALESCE(lm.type, '') as last_msg_type,
			COALESCE(lm.created_at, c.created_at) as last_time,
			(SELECT COUNT(*) FROM messenger_messages msg
			 WHERE msg.chat_id = c.id AND msg.id > cm.last_read_message_id
			   AND msg.sender_id != $1 AND msg.deleted_at IS NULL
			   AND (cm.cleared_at IS NULL OR msg.created_at > cm.cleared_at)
			) as unread_count,
			cm.is_muted,
			u.id as interlocutor_id,
			COALESCE(ct.custom_name, u.display_name, '') as interlocutor_name,
			COALESCE(u.avatar_media_id::text, '') as interlocutor_avatar,
			u.avatar_media_id as interlocutor_avatar_media_id,
			(SELECT COUNT(*) FROM messenger_chat_members WHERE chat_id = c.id) as member_count,
			EXISTS (SELECT 1 FROM messenger_blocked_users b
			        WHERE b.blocker_id = $1 AND b.blocked_id = u.id) as blocked_by_me,
			EXISTS (SELECT 1 FROM messenger_blocked_users b
			        WHERE b.blocker_id = u.id AND b.blocked_id = $1) as blocked_me
		FROM messenger_chats c
		JOIN messenger_chat_members cm ON c.id = cm.chat_id
		LEFT JOIN LATERAL (
			SELECT content, type, created_at FROM messenger_messages
			WHERE chat_id = c.id AND deleted_at IS NULL
			  AND (cm.cleared_at IS NULL OR created_at > cm.cleared_at)
			ORDER BY id DESC LIMIT 1
		) lm ON TRUE
		LEFT JOIN messenger_chat_members cm2 ON c.type = 'DIRECT' AND cm2.chat_id = c.id AND cm2.user_id != $1
		LEFT JOIN messenger_users u ON cm2.user_id = u.id
		LEFT JOIN messenger_contacts ct ON ct.owner_id = $1 AND ct.contact_id = u.id
		WHERE cm.user_id = $1
		  -- "delete for me": hide the chat until a message arrives after hidden_at.
		  AND (cm.hidden_at IS NULL OR (lm.created_at IS NOT NULL AND lm.created_at > cm.hidden_at))
		ORDER BY last_time DESC`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []chat.ChatListItemResponse
	for rows.Next() {
		var item chat.ChatListItemResponse
		err := rows.Scan(
			&item.ID, &item.Type, &item.Name, &item.Description,
			&item.AvatarMediaID, &item.CreatorID, &item.InviteLink, &item.CreatedAt,
			&item.LastMessageContent, &item.LastMessageType, &item.LastMessageTime,
			&item.UnreadCount, &item.IsMuted,
			&item.InterlocutorID, &item.InterlocutorName, &item.InterlocutorAvatar,
			&item.InterlocutorAvatarMediaID,
			&item.MemberCount,
			&item.BlockedByMe, &item.BlockedMe,
		)
		if err != nil {
			return nil, err
		}
		// Decrypt the last-message preview (DIRECT/GROUP bodies are encrypted at rest).
		item.LastMessageContent = r.cipher.DecryptContent(item.LastMessageContent)
		// last_message_content stays raw (JSON for system/transfer messages); the
		// client renders a localized preview using last_message_type.
		if item.Type == chat.TypeDirect && item.InterlocutorName != "" {
			item.Name = item.InterlocutorName
		}
		chats = append(chats, item)
	}
	return chats, rows.Err()
}

func (r *chatRepo) UpdateChat(ctx context.Context, chatID uuid.UUID, req *chat.UpdateChatReq) error {
	if req.Name != nil {
		if _, err := r.db.Exec(ctx, `UPDATE messenger_chats SET name = $1 WHERE id = $2`, *req.Name, chatID); err != nil {
			return err
		}
	}
	if req.Description != nil {
		if _, err := r.db.Exec(ctx, `UPDATE messenger_chats SET description = $1 WHERE id = $2`, *req.Description, chatID); err != nil {
			return err
		}
	}
	if req.AvatarMediaID != nil {
		if _, err := r.db.Exec(ctx, `UPDATE messenger_chats SET avatar_media_id = $1 WHERE id = $2`, *req.AvatarMediaID, chatID); err != nil {
			return err
		}
	}
	return nil
}

func (r *chatRepo) DeleteChat(ctx context.Context, chatID uuid.UUID) error {
	_, err := r.db.Exec(ctx, `DELETE FROM messenger_chats WHERE id = $1`, chatID)
	return err
}

func (r *chatRepo) IsUsernameTaken(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (SELECT 1 FROM messenger_chats WHERE username = $1)`,
		username,
	).Scan(&exists)
	return exists, err
}
