package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)


func (s *ChatService) MuteChat(ctx context.Context, userID, chatID uuid.UUID, mutedUntil *time.Time) error {
	if err := s.requireMembership(ctx, chatID, userID); err != nil {
		return err
	}

	return s.mods.MuteChat(ctx, userID, chatID, mutedUntil)
}

func (s *ChatService) UnmuteChat(ctx context.Context, userID, chatID uuid.UUID) error {
	if err := s.requireMembership(ctx, chatID, userID); err != nil {
		return err
	}

	return s.mods.UnmuteChat(ctx, userID, chatID)
}


func (s *ChatService) PinMessage(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, messageID int64) error {
	if err := s.canPinMessage(ctx, chatID, userID); err != nil {
		return err
	}

	// Проверяем что сообщение из этого чата
	msg, err := s.msgs.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("message not found: %w", err)
	}
	if msg.ChatID != chatID {
		return fmt.Errorf("message does not belong to this chat")
	}

	if err := s.mods.PinMessage(ctx, chatID, messageID, userID); err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Int64("message_id", messageID).Msg("Failed to pin message")
		return fmt.Errorf("pin message: %w", err)
	}

	log.Debug().Str("chat_id", chatID.String()).Int64("message_id", messageID).Str("user_id", userID.String()).Msg("Message pinned")

	go s.sendSystemMessage(context.Background(), chatID, chat.SystemMessagePinned)

	go s.broadcastToChat(context.Background(), chatID, socketmodel.EventMessagePinned, map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
		"pinned_by":  userID,
	})

	return nil
}

func (s *ChatService) UnpinMessage(ctx context.Context, chatID uuid.UUID, userID uuid.UUID, messageID int64) error {
	if err := s.canPinMessage(ctx, chatID, userID); err != nil {
		return err
	}

	if err := s.mods.UnpinMessage(ctx, chatID, messageID); err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Int64("message_id", messageID).Msg("Failed to unpin message")
		return fmt.Errorf("unpin message: %w", err)
	}

	log.Debug().Str("chat_id", chatID.String()).Int64("message_id", messageID).Msg("Message unpinned")

	go s.sendSystemMessage(context.Background(), chatID, chat.SystemMessageUnpinned)

	go s.broadcastToChat(context.Background(), chatID, socketmodel.EventMessageUnpinned, map[string]interface{}{
		"chat_id":    chatID,
		"message_id": messageID,
	})

	return nil
}


func (s *ChatService) GetChatInfo(ctx context.Context, chatID, userID uuid.UUID) (*chat.ChatInfoResponse, error) {
	if err := s.requireMembership(ctx, chatID, userID); err != nil {
		return nil, err
	}

	c, err := s.chats.GetChatByID(ctx, chatID)
	if err != nil || c == nil {
		return nil, fmt.Errorf("chat not found: %s", chatID)
	}

	members, err := s.members.GetChatMembersDetailed(ctx, chatID)
	if err != nil {
		return nil, fmt.Errorf("get members: %w", err)
	}

	// Обогащаем онлайн-статусом
	onlineCount := 0
	for i := range members {
		members[i].IsOnline = s.isUserOnline(ctx, members[i].UserID)
		if members[i].IsOnline {
			onlineCount++
		}
	}

	pins, err := s.mods.GetPinnedMessages(ctx, chatID)
	if err != nil {
		pins = nil // не критично
	}

	// Block status (DIRECT only) so the client can disable the composer before
	// attempting to send and getting a 403.
	var blockedByMe, blockedMe bool
	if c.Type == chat.TypeDirect {
		for i := range members {
			other := members[i].UserID
			if other == userID {
				continue
			}
			blockedByMe, _ = s.mods.IsBlocked(ctx, userID, other)
			blockedMe, _ = s.mods.IsBlocked(ctx, other, userID)
			break
		}
	}

	return &chat.ChatInfoResponse{
		Chat:           *c,
		Members:        members,
		MemberCount:    len(members),
		OnlineCount:    onlineCount,
		PinnedMessages: pins,
		BlockedByMe:    blockedByMe,
		BlockedMe:      blockedMe,
	}, nil
}

func (s *ChatService) GetChatMedia(ctx context.Context, chatID, userID uuid.UUID) (*chat.ChatMediaResponse, error) {
	if err := s.requireMembership(ctx, chatID, userID); err != nil {
		return nil, err
	}

	return s.msgs.GetChatMedia(ctx, chatID)
}

func (s *ChatService) BlockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	if err := s.mods.BlockUser(ctx, blockerID, blockedID); err != nil {
		log.Error().Err(err).Str("blocker_id", blockerID.String()).Str("blocked_id", blockedID.String()).Msg("Failed to block user")
		return fmt.Errorf("block user: %w", err)
	}

	log.Info().Str("blocker_id", blockerID.String()).Str("blocked_id", blockedID.String()).Msg("User blocked")

	s.invalidateBlockCache(ctx, blockerID, blockedID)
	// Notify the blocked user so their UI can lock the chat without waiting for
	// a 403 on the next send.
	s.notifyUser(blockedID, socketmodel.EventUserBlocked, map[string]interface{}{
		"by_user": blockerID,
	})
	return nil
}

func (s *ChatService) UnblockUser(ctx context.Context, blockerID, blockedID uuid.UUID) error {
	if err := s.mods.UnblockUser(ctx, blockerID, blockedID); err != nil {
		log.Error().Err(err).Str("blocker_id", blockerID.String()).Str("blocked_id", blockedID.String()).Msg("Failed to unblock user")
		return fmt.Errorf("unblock user: %w", err)
	}

	log.Info().Str("blocker_id", blockerID.String()).Str("blocked_id", blockedID.String()).Msg("User unblocked")

	s.invalidateBlockCache(ctx, blockerID, blockedID)
	s.notifyUser(blockedID, socketmodel.EventUserUnblocked, map[string]interface{}{
		"by_user": blockerID,
	})
	return nil
}

func (s *ChatService) GetBlockedUsers(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	return s.mods.GetBlockedUsers(ctx, userID)
}

// GetBlockedUsersDetailed returns blocked users enriched with public profile info
// for rendering the Blocked tab.
func (s *ChatService) GetBlockedUsersDetailed(ctx context.Context, userID uuid.UUID) ([]chat.BlockedUserInfo, error) {
	return s.mods.GetBlockedUsersDetailed(ctx, userID)
}