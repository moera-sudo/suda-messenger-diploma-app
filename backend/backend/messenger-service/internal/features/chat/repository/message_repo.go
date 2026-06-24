package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	appcrypto "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/crypto"
)

type messageRepo struct {
	db     *pgxpool.Pool
	cipher appcrypto.ContentCipher
}

func (r *messageRepo) SaveMessage(ctx context.Context, msg *chat.Message) error {
	// Encrypt content at rest for private chats (DIRECT/GROUP). System messages
	// and public channels stay plaintext. msg.Content is left untouched in memory
	// (plaintext), so WS broadcast and push previews keep working as before.
	plaintext := msg.Content
	stored := plaintext
	if msg.Type != chat.MsgTypeSystem &&
		(msg.ChatType == chat.TypeDirect || msg.ChatType == chat.TypeGroup) {
		stored = r.cipher.EncryptContent(plaintext)
	}

	// search_vector is built from the PLAINTEXT ($9). The DB trigger that built it
	// from `content` was dropped (migration 000016) — it would index ciphertext.
	// $9 is NULL for SYSTEM messages, so to_tsvector(NULL) -> NULL and the vector
	// ends up NULL (matching the old trigger). Using a dedicated nullable param
	// (not reusing $4) avoids Postgres deducing conflicting types for $4 (42P08).
	var searchSrc *string
	if msg.Type != chat.MsgTypeSystem {
		searchSrc = &plaintext
	}

	return r.db.QueryRow(ctx,
		`INSERT INTO messenger_messages
			(chat_id, sender_id, content, type, status, attachment_media_id,
			 reply_to_message_id, forwarded_from_chat, forwarded_from_msg, search_vector)
		 VALUES ($1, $2, $3, $4, 'SENT', $5, $6, $7, $8,
			 setweight(to_tsvector('russian', $9), 'A') ||
			 setweight(to_tsvector('simple',  $9), 'B'))
		 RETURNING id, created_at, status`,
		msg.ChatID, msg.SenderID, stored, msg.Type,
		msg.AttachmentMediaID, msg.ReplyToMessageID,
		msg.ForwardedFromChat, msg.ForwardedFromMsg, searchSrc,
	).Scan(&msg.ID, &msg.CreatedAt, &msg.Status)
}

func (r *messageRepo) GetMessageByID(ctx context.Context, messageID int64) (*chat.Message, error) {
	var m chat.Message
	err := r.db.QueryRow(ctx,
		`SELECT id, chat_id, sender_id, content, type, status,
		        attachment_media_id, reply_to_message_id,
		        forwarded_from_chat, forwarded_from_msg,
		        edited_at, deleted_at, created_at
		 FROM messenger_messages WHERE id = $1`, messageID,
	).Scan(
		&m.ID, &m.ChatID, &m.SenderID, &m.Content, &m.Type, &m.Status,
		&m.AttachmentMediaID, &m.ReplyToMessageID,
		&m.ForwardedFromChat, &m.ForwardedFromMsg,
		&m.EditedAt, &m.DeletedAt, &m.CreatedAt,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("message not found: %d", messageID)
	}
	if err != nil {
		return nil, err
	}
	m.Content = r.cipher.DecryptContent(m.Content)
	return &m, nil
}

// GetMessagePreview returns a short, hydrated snippet of a message by id —
// used to render reply / forwarded-from blocks. Returns (nil, nil) if the
// referenced message no longer exists.
func (r *messageRepo) GetMessagePreview(ctx context.Context, messageID int64) (*chat.MessagePreview, error) {
	var (
		p          chat.MessagePreview
		senderName string
		content    string
		deleted    bool
	)
	err := r.db.QueryRow(ctx,
		`SELECT m.id, m.sender_id, COALESCE(u.display_name, u.username, ''),
		        m.content, m.type, (m.deleted_at IS NOT NULL)
		 FROM messenger_messages m
		 LEFT JOIN messenger_users u ON u.id = m.sender_id
		 WHERE m.id = $1`, messageID,
	).Scan(&p.MessageID, &p.SenderID, &senderName, &content, &p.Type, &deleted)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	p.SenderName = senderName
	p.Deleted = deleted
	if !deleted {
		p.Content = chat.PreviewSnippet(r.cipher.DecryptContent(content))
	}
	return &p, nil
}

func (r *messageRepo) GetMessages(ctx context.Context, chatID, userID uuid.UUID, limit, offset int) ([]chat.Message, error) {
	// LEFT JOIN the reply target (rm/ru) and the forward source (fm/fu) so each
	// message carries a hydrated preview and the client never depends on having
	// the referenced message in its local cache.
	rows, err := r.db.Query(ctx,
		`SELECT m.id, m.chat_id, m.sender_id, m.content, m.type, m.status,
		        m.attachment_media_id, m.reply_to_message_id,
		        m.forwarded_from_chat, m.forwarded_from_msg,
		        m.edited_at, m.deleted_at, m.created_at,
		        rm.id, rm.sender_id, COALESCE(ru.display_name, ru.username, ''),
		        rm.content, rm.type, (rm.deleted_at IS NOT NULL),
		        fm.id, fm.sender_id, COALESCE(fu.display_name, fu.username, ''),
		        fm.content, fm.type, (fm.deleted_at IS NOT NULL)
		 FROM messenger_messages m
		 LEFT JOIN messenger_messages rm ON rm.id = m.reply_to_message_id
		 LEFT JOIN messenger_users ru ON ru.id = rm.sender_id
		 LEFT JOIN messenger_messages fm ON fm.id = m.forwarded_from_msg
		 LEFT JOIN messenger_users fu ON fu.id = fm.sender_id
		 -- LEFT JOIN (not INNER) so public-channel readers who aren't members still see messages.
		 LEFT JOIN messenger_chat_members cmv ON cmv.chat_id = m.chat_id AND cmv.user_id = $2
		 WHERE m.chat_id = $1 AND m.deleted_at IS NULL
		   AND NOT EXISTS (
		       SELECT 1 FROM messenger_message_hidden h
		       WHERE h.message_id = m.id AND h.user_id = $2
		   )
		   -- "delete for me": hide history cleared by this member.
		   AND (cmv.cleared_at IS NULL OR m.created_at > cmv.cleared_at)
		 ORDER BY m.id DESC LIMIT $3 OFFSET $4`,
		chatID, userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []chat.Message
	for rows.Next() {
		var (
			m chat.Message

			rID       *int64
			rSender   *uuid.UUID
			rName     string
			rContent  *string
			rType     *string
			rDeleted  *bool

			fID      *int64
			fSender  *uuid.UUID
			fName    string
			fContent *string
			fType    *string
			fDeleted *bool
		)
		if err := rows.Scan(
			&m.ID, &m.ChatID, &m.SenderID, &m.Content, &m.Type, &m.Status,
			&m.AttachmentMediaID, &m.ReplyToMessageID,
			&m.ForwardedFromChat, &m.ForwardedFromMsg,
			&m.EditedAt, &m.DeletedAt, &m.CreatedAt,
			&rID, &rSender, &rName, &rContent, &rType, &rDeleted,
			&fID, &fSender, &fName, &fContent, &fType, &fDeleted,
		); err != nil {
			return nil, err
		}
		// Decrypt the message body and the hydrated reply/forward snippets.
		m.Content = r.cipher.DecryptContent(m.Content)
		if rContent != nil {
			d := r.cipher.DecryptContent(*rContent)
			rContent = &d
		}
		if fContent != nil {
			d := r.cipher.DecryptContent(*fContent)
			fContent = &d
		}
		m.ReplyPreview = buildPreview(rID, rSender, rName, rContent, rType, rDeleted)
		m.ForwardedFrom = buildPreview(fID, fSender, fName, fContent, fType, fDeleted)
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

// buildPreview assembles a MessagePreview from nullable joined columns.
// Returns nil when the referenced row is absent (no reply / no forward source).
func buildPreview(id *int64, sender *uuid.UUID, name string, content, msgType *string, deleted *bool) *chat.MessagePreview {
	if id == nil {
		return nil
	}
	p := &chat.MessagePreview{
		MessageID:  *id,
		SenderID:   sender,
		SenderName: name,
	}
	if msgType != nil {
		p.Type = *msgType
	}
	if deleted != nil && *deleted {
		p.Deleted = true
		return p
	}
	if content != nil {
		p.Content = chat.PreviewSnippet(*content)
	}
	return p
}

func (r *messageRepo) EditMessage(ctx context.Context, messageID int64, content string) error {
	// Preserve the message's encryption state: if the stored body was encrypted,
	// re-encrypt the edited text; otherwise keep it plaintext (legacy / channel).
	var wasEncrypted bool
	if err := r.db.QueryRow(ctx,
		`SELECT content LIKE 'enc:v1:%' FROM messenger_messages WHERE id = $1`,
		messageID,
	).Scan(&wasEncrypted); err != nil {
		return fmt.Errorf("message not found: %d", messageID)
	}

	stored := content
	if wasEncrypted {
		stored = r.cipher.EncryptContent(content)
	}

	// Rebuild search_vector from the plaintext ($3) — the DB trigger was dropped.
	tag, err := r.db.Exec(ctx,
		`UPDATE messenger_messages
		 SET content = $1, edited_at = NOW(),
		     search_vector = setweight(to_tsvector('russian', $3), 'A') ||
		                     setweight(to_tsvector('simple',  $3), 'B')
		 WHERE id = $2 AND deleted_at IS NULL`,
		stored, messageID, content,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("message not found or deleted: %d", messageID)
	}
	return nil
}

func (r *messageRepo) SoftDeleteMessage(ctx context.Context, messageID int64) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE messenger_messages SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
		messageID,
	)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("message not found or already deleted: %d", messageID)
	}
	return nil
}

func (r *messageRepo) HideMessageForUser(ctx context.Context, userID uuid.UUID, messageID int64) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO messenger_message_hidden (user_id, message_id)
		 VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		userID, messageID,
	)
	return err
}

func (r *messageRepo) GetChatMedia(ctx context.Context, chatID uuid.UUID) (*chat.ChatMediaResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT m.id, m.attachment_media_id, m.type, m.content, m.sender_id, m.created_at,
		        me.content_type, me.original_name
		 FROM messenger_messages m
		 LEFT JOIN media me ON me.id = m.attachment_media_id
		 WHERE m.chat_id = $1 AND m.deleted_at IS NULL
		   AND (m.attachment_media_id IS NOT NULL OR m.content ~ 'https?://')
		 ORDER BY m.id DESC`, chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &chat.ChatMediaResponse{}
	for rows.Next() {
		var (
			msgID       int64
			mediaID     *uuid.UUID
			msgType     string
			content     string
			senderID    *uuid.UUID
			createdAt   time.Time
			contentType *string
			origName    *string
		)
		if err := rows.Scan(&msgID, &mediaID, &msgType, &content, &senderID, &createdAt, &contentType, &origName); err != nil {
			return nil, err
		}

		if mediaID != nil && contentType != nil {
			item := chat.MediaItemResponse{
				MessageID:    msgID,
				MediaID:      *mediaID,
				ContentType:  *contentType,
				OriginalName: origName,
				SenderID:     senderID,
				CreatedAt:    createdAt,
			}
			ct := *contentType
			switch {
			case strings.HasPrefix(ct, "image/"):
				resp.Images = append(resp.Images, item)
			case strings.HasPrefix(ct, "video/"):
				resp.Videos = append(resp.Videos, item)
			case strings.HasPrefix(ct, "audio/"):
				resp.Audio = append(resp.Audio, item)
			default:
				resp.Documents = append(resp.Documents, item)
			}
		}
	}
	return resp, rows.Err()
}
