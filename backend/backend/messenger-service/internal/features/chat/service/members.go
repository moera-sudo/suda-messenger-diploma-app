package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences"
)

func (s *ChatService) AddMember(ctx context.Context, chatID, requesterID, newUserID uuid.UUID) error {
	if err := s.canManageMembers(ctx, chatID, requesterID); err != nil {
		return err
	}

	// Проверяем что не добавляем заблокированного
	blocked, err := s.isBlockedEither(ctx, requesterID, newUserID)
	if err != nil {
		return err
	}
	if blocked {
		return ErrBlocked
	}
	if s.prefs != nil {
		targetPrefs, err := s.prefs.GetPreferences(ctx, newUserID)
		if err == nil {
			switch targetPrefs.WhoCanAddToGroups {
			case preferences.VisibilityNobody:
				return fmt.Errorf("forbidden: this user does not allow being added to groups")
			case preferences.VisibilityContacts:
				_, err := s.mods.GetContactName(ctx, newUserID, requesterID)
				if err != nil {
					return fmt.Errorf("forbidden: this user only accepts group invites from contacts")
				}
			}
		}
	}

	members, err := s.getChatMembersCached(ctx, chatID)
	if err != nil {
		return err
	}
	for _, memberID := range members {
		blocked, err := s.isBlockedEither(ctx, newUserID, memberID)
		if err != nil {
			return err
		}
		if blocked {
			return fmt.Errorf("forbidden: user has a block relationship with existing member")
		}
	}

	if err := s.members.AddMember(ctx, chatID, newUserID, chat.RoleMember); err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Str("new_user_id", newUserID.String()).Msg("Failed to add member")
		return fmt.Errorf("add member: %w", err)
	}

	log.Debug().Str("chat_id", chatID.String()).Str("new_user_id", newUserID.String()).Str("requester_id", requesterID.String()).Msg("Member added")

	s.invalidateMembersCache(ctx, chatID)

	go s.sendSystemMessage(context.Background(), chatID, chat.SystemUserJoined)

	go s.broadcastToChat(context.Background(), chatID, socketmodel.EventMemberAdded, map[string]interface{}{
		"chat_id": chatID,
		"user_id": newUserID,
	})

	return nil
}

func (s *ChatService) RemoveMember(ctx context.Context, chatID, requesterID, targetID uuid.UUID) error {
	if err := s.canManageMembers(ctx, chatID, requesterID); err != nil {
		return err
	}

	// Нельзя удалить владельца
	targetRole, err := s.members.GetMemberRole(ctx, chatID, targetID)
	if err != nil {
		return fmt.Errorf("target not a member: %w", err)
	}
	if targetRole == chat.RoleOwner {
		return errors.New("forbidden: cannot remove the chat owner")
	}

	// Админ не может удалить другого админа
	requesterRole, _ := s.members.GetMemberRole(ctx, chatID, requesterID)
	if requesterRole == chat.RoleAdmin && targetRole == chat.RoleAdmin {
		return errors.New("forbidden: admin cannot remove another admin")
	}

	if err := s.members.RemoveMember(ctx, chatID, targetID); err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Str("target_id", targetID.String()).Msg("Failed to remove member")
		return fmt.Errorf("remove member: %w", err)
	}

	log.Debug().Str("chat_id", chatID.String()).Str("target_id", targetID.String()).Str("requester_id", requesterID.String()).Msg("Member removed")

	s.invalidateMembersCache(ctx, chatID)

	go s.sendSystemMessage(context.Background(), chatID, chat.SystemUserRemoved)

	go s.broadcastToChat(context.Background(), chatID, socketmodel.EventMemberRemoved, map[string]interface{}{
		"chat_id": chatID,
		"user_id": targetID,
	})

	return nil
}

func (s *ChatService) LeaveChat(ctx context.Context, chatID, userID uuid.UUID) error {
	chatType, err := s.getChatType(ctx, chatID)
	if err != nil {
		return err
	}

	if chatType == chat.TypeDirect || chatType == chat.TypeSaved {
		return errors.New("forbidden: cannot leave this chat type")
	}

	if err := s.requireMembership(ctx, chatID, userID); err != nil {
		return err
	}

	// Если владелец уходит — проверяем есть ли другие участники
	role, _ := s.members.GetMemberRole(ctx, chatID, userID)
	if role == chat.RoleOwner {
		count, _ := s.members.GetMemberCount(ctx, chatID)
		if count > 1 {
			// Передаём владение первому админу, или первому участнику
			members, _ := s.members.GetChatMembers(ctx, chatID)
			for _, uid := range members {
				if uid != userID {
					if err := s.members.UpdateMemberRole(ctx, chatID, uid, chat.RoleOwner); err != nil {
						log.Warn().Err(err).Msg("Failed to transfer ownership")
					}
					break
				}
			}
		}
	}

	if err := s.members.RemoveMember(ctx, chatID, userID); err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Str("user_id", userID.String()).Msg("Failed to leave chat")
		return fmt.Errorf("leave chat: %w", err)
	}

	log.Debug().Str("chat_id", chatID.String()).Str("user_id", userID.String()).Msg("User left chat")

	s.invalidateMembersCache(ctx, chatID)

	go s.sendSystemMessage(context.Background(), chatID, chat.SystemUserLeft)

	go s.broadcastToChat(context.Background(), chatID, socketmodel.EventMemberLeft, map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
	})

	// Если в чате никого не осталось — удаляем чат
	count, _ := s.members.GetMemberCount(ctx, chatID)
	if count == 0 {
		_ = s.chats.DeleteChat(ctx, chatID)
		s.redis.Del(ctx, fmt.Sprintf("chat:%s:type", chatID))
	}

	return nil
}

func (s *ChatService) UpdateMemberRole(ctx context.Context, chatID, requesterID, targetID uuid.UUID, role string) error {
	// Только владелец может менять роли
	if err := s.requireRole(ctx, chatID, requesterID, chat.RoleOwner); err != nil {
		return ErrNotOwner
	}

	if requesterID == targetID {
		return errors.New("forbidden: cannot change own role")
	}

	if role == chat.RoleOwner {
		return errors.New("forbidden: use transfer ownership instead")
	}

	if err := s.members.UpdateMemberRole(ctx, chatID, targetID, role); err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Str("target_id", targetID.String()).Str("role", role).Msg("Failed to update member role")
		return fmt.Errorf("update role: %w", err)
	}

	log.Info().Str("chat_id", chatID.String()).Str("target_id", targetID.String()).Str("new_role", role).Str("requester_id", requesterID.String()).Msg("Member role updated")

	go s.broadcastToChat(context.Background(), chatID, socketmodel.EventChatUpdated, map[string]interface{}{
		"chat_id":  chatID,
		"user_id":  targetID,
		"new_role": role,
	})

	return nil
}

// BroadcastTyping — рассылает TYPING всем участникам кроме отправителя.
// stopped=true означает, что отправитель перестал печатать — получатель сразу
// скрывает индикатор, не дожидаясь debounce.
func (s *ChatService) BroadcastTyping(userID uuid.UUID, chatID uuid.UUID, stopped bool) {
	go s.broadcastToChatExcept(context.Background(), chatID, userID, socketmodel.EventTypingNotify, map[string]interface{}{
		"chat_id": chatID,
		"user_id": userID,
		"stopped": stopped,
	})
}