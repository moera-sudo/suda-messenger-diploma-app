package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)

var (
	ErrForbidden         = errors.New("forbidden")
	ErrNotMember         = errors.New("forbidden: not a member of this chat")
	ErrNotAdmin          = errors.New("forbidden: admin rights required")
	ErrNotOwner          = errors.New("forbidden: owner rights required")
	ErrCannotSend        = errors.New("forbidden: cannot send messages to this chat")
	ErrBlocked           = errors.New("forbidden: user is blocked")
	ErrCannotEditOther   = errors.New("forbidden: can only edit own messages")
	ErrCannotDeleteOther = errors.New("forbidden: insufficient rights to delete this message")
	// ErrGatingRequired — public channel has a token-gating rule the viewer doesn't meet.
	ErrGatingRequired = errors.New("forbidden: token gating — meet the requirement to view this channel")
	// ErrInvalidReplyTarget — reply_to_message_id points to a missing/deleted
	// message or one from another chat (client must not reference such ids).
	ErrInvalidReplyTarget = errors.New("Invalid reply target: referenced message is not part of this chat")
)

func (s *ChatService) requireMembership(ctx context.Context, chatID, userID uuid.UUID) error {
	isMember, err := s.members.IsMember(ctx, chatID, userID)
	if err != nil {
		return fmt.Errorf("failed to check membership: %w", err)
	}
	if !isMember {
		log.Warn().Str("chat_id", chatID.String()).Str("user_id", userID.String()).Msg("Permission denied: not a member")
		return ErrNotMember
	}
	return nil
}

// requireReadAccess gates reading chat content. A PUBLIC channel is readable by
// anyone (no subscription needed) — UNLESS it has a token-gating rule the viewer
// doesn't meet. Members always read. Every other chat type (incl. PRIVATE
// channels) still requires membership; gating never applies to private channels.
func (s *ChatService) requireReadAccess(ctx context.Context, chatID, userID uuid.UUID) error {
	c, err := s.chats.GetChatByID(ctx, chatID)
	if err != nil {
		return fmt.Errorf("failed to load chat: %w", err)
	}
	if c == nil {
		return ErrNotMember
	}
	// "PUBLIC" literal avoids importing the channel package into chat/service.
	if c.Type == chat.TypeChannel && c.Visibility == "PUBLIC" {
		// Members (already subscribed/passed gating) always read.
		if isMember, _ := s.members.IsMember(ctx, chatID, userID); isMember {
			return nil
		}
		// Non-member: a PAID public channel (gating rule with price > 0) is
		// readable only by subscribers — non-members must subscribe (pay) first.
		// A free public channel (no rule / price 0) stays openly readable.
		if s.txClient != nil {
			g, gerr := s.txClient.CheckTokenGating(ctx, chatID.String(), userID.String())
			if gerr != nil {
				// Fail-open: don't break public channels when tx-service is down.
				log.Warn().Err(gerr).Str("chat_id", chatID.String()).Msg("gating check failed on read, allowing (fail-open)")
				return nil
			}
			if g.Required && g.PriceWei != "" && g.PriceWei != "0" {
				return ErrGatingRequired
			}
		}
		return nil
	}
	return s.requireMembership(ctx, chatID, userID)
}

func (s *ChatService) requireRole(ctx context.Context, chatID, userID uuid.UUID, roles ...string) error {
	role, err := s.members.GetMemberRole(ctx, chatID, userID)
	if err != nil {
		return ErrNotMember
	}

	for _, r := range roles {
		if role == r {
			return nil
		}
	}

	log.Warn().Str("chat_id", chatID.String()).Str("user_id", userID.String()).Msg("Permission denied: insufficient role")
	return ErrNotAdmin
}

func (s *ChatService) canSendMessage(ctx context.Context, chatID, userID uuid.UUID) error {
	chatType, err := s.getChatType(ctx, chatID)
	if err != nil {
		return err
	}

	switch chatType {
	case chat.TypeChannel:
        role, err := s.members.GetMemberRole(ctx, chatID, userID)
        if err != nil {
            return err
        }
        if role != chat.RoleOwner && role != chat.RoleAdmin {
            return fmt.Errorf("forbidden: only admins can post in this channel")
        }
        return nil
	case chat.TypeDirect:
		if err := s.requireMembership(ctx, chatID, userID); err != nil {
			return err
		}
		// Проверяем блокировку в direct
		members, err := s.getChatMembersCached(ctx, chatID)
		if err != nil {
			return err
		}
		for _, memberID := range members {
			if memberID != userID {
				blocked, err := s.isBlockedEither(ctx, userID, memberID)
				if err != nil {
					return err
				}
				if blocked {
					return ErrBlocked
				}
			}
		}
		return nil
	default:
		return s.requireMembership(ctx, chatID, userID)
	}
}

func (s *ChatService) canEditMessage(ctx context.Context, msg *chat.Message, userID uuid.UUID) error {
	if msg.SenderID == nil || *msg.SenderID != userID {
		return ErrCannotEditOther
	}
	if msg.Type == chat.MsgTypeSystem {
		return errors.New("forbidden: cannot edit system messages")
	}
	return nil
}

func (s *ChatService) canDeleteMessage(ctx context.Context, msg *chat.Message, userID uuid.UUID, forEveryone bool) error {
	// Удалить у себя может кто угодно
	if !forEveryone {
		return nil
	}

	// Удалить у всех — автор или админ
	if msg.SenderID != nil && *msg.SenderID == userID {
		return nil
	}

	chatType, err := s.getChatType(ctx, msg.ChatID)
	if err != nil {
		return err
	}

	if chatType == chat.TypeGroup || chatType == chat.TypeChannel {
		if err := s.requireRole(ctx, msg.ChatID, userID, chat.RoleOwner, chat.RoleAdmin); err == nil {
			return nil
		}
	}

	return ErrCannotDeleteOther
}

func (s *ChatService) canManageMembers(ctx context.Context, chatID, userID uuid.UUID) error {
	chatType, err := s.getChatType(ctx, chatID)
	if err != nil {
		return err
	}

	if chatType == chat.TypeDirect || chatType == chat.TypeSaved {
		return errors.New("forbidden: cannot manage members in this chat type")
	}

	if chatType == chat.TypeChannel { 
		return errors.New("forbidden: use subscribe/unsubscribe for channels")
	}

	return s.requireRole(ctx, chatID, userID, chat.RoleOwner, chat.RoleAdmin)
}

func (s *ChatService) canPinMessage(ctx context.Context, chatID, userID uuid.UUID) error {
	chatType, err := s.getChatType(ctx, chatID)
	if err != nil {
		return err
	}

	switch chatType {
	case chat.TypeDirect, chat.TypeSaved:
		return s.requireMembership(ctx, chatID, userID)
	default:
		return s.requireRole(ctx, chatID, userID, chat.RoleOwner, chat.RoleAdmin)
	}
}

func (s *ChatService) canUpdateChat(ctx context.Context, chatID, userID uuid.UUID) error {
	chatType, err := s.getChatType(ctx, chatID)
	if err != nil {
		return err
	}

	if chatType == chat.TypeDirect || chatType == chat.TypeSaved {
		return errors.New("forbidden: cannot update this chat type")
	}

	return s.requireRole(ctx, chatID, userID, chat.RoleOwner, chat.RoleAdmin)
}

func (s *ChatService) canDeleteChat(ctx context.Context, chatID, userID uuid.UUID) error {
	chatType, err := s.getChatType(ctx, chatID)
	if err != nil {
		return err
	}

	switch chatType {
	case chat.TypeDirect, chat.TypeSaved:
		// DIRECT/SAVED have no OWNER role — any member may delete the whole chat.
		return s.requireMembership(ctx, chatID, userID)
	default:
		return s.requireRole(ctx, chatID, userID, chat.RoleOwner)
	}
}