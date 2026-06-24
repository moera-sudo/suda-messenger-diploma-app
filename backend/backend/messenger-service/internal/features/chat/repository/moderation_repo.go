package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	appcrypto "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/crypto"
)

type moderationRepo struct {
	db     *pgxpool.Pool
	cipher appcrypto.ContentCipher
}

// * Mute/Unmute methods
func (r *moderationRepo) MuteChat(ctx context.Context, userID, chatID uuid.UUID, mutedUntil *time.Time) error {
	_, err := r.db.Exec(ctx,
		`UPDATE messenger_chat_members SET is_muted = TRUE, muted_until = $1
		 WHERE chat_id = $2 AND user_id = $3`,
		mutedUntil, chatID, userID,
	)
	return err
}

func (r *moderationRepo) UnmuteChat(ctx context.Context, userID, chatID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE messenger_chat_members SET is_muted = FALSE, muted_until = NULL
		 WHERE chat_id = $1 AND user_id = $2`,
		chatID, userID,
	)
	return err
}

// ClearAndHideChat implements "delete chat for me": hides the chat from the
// member's list and clears their visible history. The chat reappears (showing
// only newer messages) once a message arrives after hidden_at.
func (r *moderationRepo) ClearAndHideChat(ctx context.Context, userID, chatID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`UPDATE messenger_chat_members SET cleared_at = NOW(), hidden_at = NOW()
		 WHERE chat_id = $1 AND user_id = $2`,
		chatID, userID,
	)
	return err
}

func (r *moderationRepo) IsChatMuted(ctx context.Context, userID, chatID uuid.UUID) (bool, error) {
	var isMuted bool
	var mutedUntil *time.Time
	err := r.db.QueryRow(ctx,
		`SELECT is_muted, muted_until FROM messenger_chat_members
		 WHERE chat_id = $1 AND user_id = $2`,
		chatID, userID,
	).Scan(&isMuted, &mutedUntil)
	if err != nil {
		return false, nil
	}
	if !isMuted {
		return false, nil
	}
	// Если muted_until задан и прошёл — не замучен
	if mutedUntil != nil && mutedUntil.Before(time.Now()) {
		return false, nil
	}
	return true, nil
}

// * Block Methods
func (r *moderationRepo) BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO messenger_blocked_users (blocker_id, blocked_id)
		 VALUES ($1, $2) ON CONFLICT DO NOTHING`,
		blockerID, blockedID,
	)
	return err
}

func (r *moderationRepo) UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM messenger_blocked_users WHERE blocker_id = $1 AND blocked_id = $2`,
		blockerID, blockedID,
	)
	return err
}

func (r *moderationRepo) IsBlocked(ctx context.Context, blockerID, blockedID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS (
			SELECT 1 FROM messenger_blocked_users WHERE blocker_id = $1 AND blocked_id = $2
		)`, blockerID, blockedID,
	).Scan(&exists)
	return exists, err
}

func (r *moderationRepo) GetBlockedUsers(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.db.Query(ctx,
		`SELECT blocked_id FROM messenger_blocked_users WHERE blocker_id = $1`, userID,
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

// GetBlockedUsersDetailed returns blocked users joined with their public profile
// (username, display name, avatar) so the Blocked tab can render names without
// hitting GET /user/id/{id}, which is 403 for blocked users.
func (r *moderationRepo) GetBlockedUsersDetailed(ctx context.Context, userID uuid.UUID) ([]chat.BlockedUserInfo, error) {
	rows, err := r.db.Query(ctx,
		`SELECT u.id, u.username, u.display_name, u.avatar_media_id
		 FROM   messenger_blocked_users b
		 JOIN   messenger_users u ON u.id = b.blocked_id
		 WHERE  b.blocker_id = $1
		 ORDER BY u.username`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]chat.BlockedUserInfo, 0)
	for rows.Next() {
		var bu chat.BlockedUserInfo
		if err := rows.Scan(&bu.UserID, &bu.Username, &bu.DisplayName, &bu.AvatarMediaID); err != nil {
			return nil, err
		}
		out = append(out, bu)
	}
	return out, rows.Err()
}

// * Pin methods
func (r *moderationRepo) PinMessage(ctx context.Context, chatID uuid.UUID, messageID int64, pinnedBy uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO messenger_pinned_messages (chat_id, message_id, pinned_by)
		 VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`,
		chatID, messageID, pinnedBy,
	)
	return err
}

func (r *moderationRepo) UnpinMessage(ctx context.Context, chatID uuid.UUID, messageID int64) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM messenger_pinned_messages WHERE chat_id = $1 AND message_id = $2`,
		chatID, messageID,
	)
	return err
}

func (r *moderationRepo) GetPinnedMessages(ctx context.Context, chatID uuid.UUID) ([]chat.PinnedMessageInfoResponse, error) {
	rows, err := r.db.Query(ctx,
		`SELECT p.chat_id, p.message_id, p.pinned_by, p.pinned_at,
		        m.content, m.sender_id, m.type
		 FROM messenger_pinned_messages p
		 JOIN messenger_messages m ON p.message_id = m.id
		 WHERE p.chat_id = $1 AND m.deleted_at IS NULL
		 ORDER BY p.pinned_at DESC`, chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var pins []chat.PinnedMessageInfoResponse
	for rows.Next() {
		var p chat.PinnedMessageInfoResponse
		if err := rows.Scan(&p.ChatID, &p.MessageID, &p.PinnedBy, &p.PinnedAt, &p.Content, &p.SenderID, &p.Type); err != nil {
			return nil, err
		}
		p.Content = r.cipher.DecryptContent(p.Content)
		pins = append(pins, p)
	}
	return pins, rows.Err()
}

// * Contacts
func (r *moderationRepo) SetContactName(ctx context.Context, ownerID, contactID uuid.UUID, name string) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO messenger_contacts (owner_id, contact_id, custom_name)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (owner_id, contact_id) DO UPDATE SET custom_name = $3, updated_at = NOW()`,
		ownerID, contactID, name,
	)
	return err
}

func (r *moderationRepo) RemoveContactName(ctx context.Context, ownerID, contactID uuid.UUID) error {
	_, err := r.db.Exec(ctx,
		`DELETE FROM messenger_contacts WHERE owner_id = $1 AND contact_id = $2`,
		ownerID, contactID,
	)
	return err
}

func (r *moderationRepo) GetContactName(ctx context.Context, ownerID, contactID uuid.UUID) (string, error) {
	var name string
	err := r.db.QueryRow(ctx,
		`SELECT custom_name FROM messenger_contacts WHERE owner_id = $1 AND contact_id = $2`,
		ownerID, contactID,
	).Scan(&name)
	if err != nil {
		return "", nil
	}
	return name, nil
}

func (r *moderationRepo) GetContacts(ctx context.Context, ownerID uuid.UUID) ([]chat.Contact, error) {
	rows, err := r.db.Query(ctx,
		`SELECT owner_id, contact_id, custom_name, created_at, updated_at
		 FROM messenger_contacts WHERE owner_id = $1 ORDER BY custom_name`, ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []chat.Contact
	for rows.Next() {
		var c chat.Contact
		if err := rows.Scan(&c.OwnerID, &c.ContactID, &c.CustomName, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, rows.Err()
}
