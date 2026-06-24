package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)

func (s *ChatService) SendMessage(ctx context.Context, userID uuid.UUID, req *chat.SendMessageReq) (*chat.Message, error) {
	if err := s.canSendMessage(ctx, req.ChatID, userID); err != nil {
		return nil, err
	}

	// Validate the reply target up front so a stale/foreign reply_to_message_id
	// surfaces as a clean error instead of an FK-violation 500 on INSERT.
	if req.ReplyToMessageID != nil {
		target, err := s.msgs.GetMessageByID(ctx, *req.ReplyToMessageID)
		if err != nil || target.ChatID != req.ChatID || target.DeletedAt != nil {
			log.Warn().Str("chat_id", req.ChatID.String()).Int64("reply_to", *req.ReplyToMessageID).Msg("Rejected message with invalid reply target")
			return nil, ErrInvalidReplyTarget
		}
	}

	msg := &chat.Message{
		ChatID:            req.ChatID,
		SenderID:          &userID,
		Content:           req.Content,
		Type:              req.Type,
		AttachmentMediaID: req.AttachmentMediaID,
		ReplyToMessageID:  req.ReplyToMessageID,
		ClientSideID:      req.ClientSideID,
	}

	// Tell the repo which chat type this is so it can decide whether to encrypt
	// the content at rest (DIRECT/GROUP only).
	if ct, err := s.getChatType(ctx, req.ChatID); err == nil {
		msg.ChatType = ct
	}

	if err := s.msgs.SaveMessage(ctx, msg); err != nil {
		log.Error().Err(err).Str("chat_id", req.ChatID.String()).Str("sender_id", userID.String()).Msg("Failed to save message")
		return nil, fmt.Errorf("save message: %w", err)
	}

	// Hydrate the reply preview so the realtime NEW_MESSAGE payload renders the
	// reply block without the client needing the target in its local cache.
	if msg.ReplyToMessageID != nil {
		if preview, err := s.msgs.GetMessagePreview(ctx, *msg.ReplyToMessageID); err == nil {
			msg.ReplyPreview = preview
		}
	}

	log.Debug().Str("chat_id", req.ChatID.String()).Str("sender_id", userID.String()).Int64("message_id", msg.ID).Msg("Message sent")

	// Привязываем медиа к сообщению для ACL
	if msg.AttachmentMediaID != nil && s.media != nil {
		go func() {
			err := s.media.LinkMediaToEntity(
				context.Background(),
				msg.AttachmentMediaID.String(),
				"MESSAGE",
				fmt.Sprintf("%d", msg.ID),
			)
			if err != nil {
				log.Error().Err(err).Msg("Failed to link media to message")
			}
		}()
	}

	go s.broadcastMessageWithPush(context.Background(), msg)

	return msg, nil
}

func (s *ChatService) EditMessage(ctx context.Context, userID uuid.UUID, messageID int64, req *chat.EditMessageReq) error {
	msg, err := s.msgs.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("message not found: %w", err)
	}

	if err := s.requireMembership(ctx, msg.ChatID, userID); err != nil {
		return err
	}

	if err := s.canEditMessage(ctx, msg, userID); err != nil {
		return err
	}

	if err := s.msgs.EditMessage(ctx, messageID, req.Content); err != nil {
		log.Error().Err(err).Int64("message_id", messageID).Msg("Failed to edit message")
		return fmt.Errorf("edit message: %w", err)
	}

	log.Debug().Int64("message_id", messageID).Str("user_id", userID.String()).Msg("Message edited")

	go s.broadcastToChat(context.Background(), msg.ChatID, socketmodel.EventMessageEdited, map[string]interface{}{
		"message_id": messageID,
		"chat_id":    msg.ChatID,
		"content":    req.Content,
	})

	return nil
}

func (s *ChatService) DeleteMessage(ctx context.Context, userID uuid.UUID, messageID int64, req *chat.DeleteMessageReq) error {
	msg, err := s.msgs.GetMessageByID(ctx, messageID)
	if err != nil {
		return fmt.Errorf("message not found: %w", err)
	}

	if err := s.requireMembership(ctx, msg.ChatID, userID); err != nil {
		return err
	}

	if err := s.canDeleteMessage(ctx, msg, userID, req.ForEveryone); err != nil {
		return err
	}

	if req.ForEveryone {
		if err := s.msgs.SoftDeleteMessage(ctx, messageID); err != nil {
			log.Error().Err(err).Int64("message_id", messageID).Msg("Failed to soft-delete message")
			return fmt.Errorf("delete message: %w", err)
		}

		log.Debug().Int64("message_id", messageID).Str("user_id", userID.String()).Msg("Message deleted for everyone")

		go s.broadcastToChat(context.Background(), msg.ChatID, socketmodel.EventMessageDeleted, map[string]interface{}{
			"message_id": messageID,
			"chat_id":    msg.ChatID,
		})
	} else {
		if err := s.msgs.HideMessageForUser(ctx, userID, messageID); err != nil {
			return fmt.Errorf("hide message: %w", err)
		}

		log.Debug().Int64("message_id", messageID).Str("user_id", userID.String()).Msg("Message hidden for user")
	}

	return nil
}

func (s *ChatService) ForwardMessage(ctx context.Context, userID uuid.UUID, req *chat.ForwardMessageReq) (*chat.Message, error) {
	// Проверяем membership в обоих чатах
	if err := s.requireMembership(ctx, req.FromChatID, userID); err != nil {
		return nil, fmt.Errorf("source chat: %w", err)
	}

	if err := s.canSendMessage(ctx, req.ToChatID, userID); err != nil {
		return nil, fmt.Errorf("target chat: %w", err)
	}

	targetType, err := s.getChatType(ctx, req.ToChatID)
	if err != nil {
		return nil, err
	}
	if targetType == chat.TypeDirect {
		members, err := s.getChatMembersCached(ctx, req.ToChatID)
		if err != nil {
			return nil, err
		}
		for _, memberID := range members {
			if memberID != userID {
				blocked, err := s.isBlockedEither(ctx, userID, memberID)
				if err != nil {
					return nil, err
				}
				if blocked {
					return nil, ErrBlocked
				}
			}
		}
	}

	original, err := s.msgs.GetMessageByID(ctx, req.MessageID)
	if err != nil {
		return nil, fmt.Errorf("original message not found: %w", err)
	}

	if original.ChatID != req.FromChatID {
		return nil, fmt.Errorf("message does not belong to source chat")
	}

	forwarded := &chat.Message{
		ChatID:            req.ToChatID,
		SenderID:          &userID,
		Content:           original.Content, // already decrypted on read; re-encrypted per target type
		Type:              original.Type,
		AttachmentMediaID: original.AttachmentMediaID,
		ForwardedFromChat: &req.FromChatID,
		ForwardedFromMsg:  &req.MessageID,
		ChatType:          targetType,
	}

	if err := s.msgs.SaveMessage(ctx, forwarded); err != nil {
		log.Error().Err(err).Str("to_chat", req.ToChatID.String()).Str("user_id", userID.String()).Msg("Failed to save forwarded message")
		return nil, fmt.Errorf("save forwarded message: %w", err)
	}

	// Hydrate the "forwarded from" preview (original sender + snippet) so the
	// realtime payload and UI can render the attribution badge.
	if preview, err := s.msgs.GetMessagePreview(ctx, req.MessageID); err == nil {
		forwarded.ForwardedFrom = preview
	}

	log.Debug().Int64("message_id", forwarded.ID).Str("from_chat", req.FromChatID.String()).Str("to_chat", req.ToChatID.String()).Msg("Message forwarded")

	// Привязываем медиа к пересланному сообщению
	if forwarded.AttachmentMediaID != nil && s.media != nil {
		go func() {
			_ = s.media.LinkMediaToEntity(
				context.Background(),
				forwarded.AttachmentMediaID.String(),
				"MESSAGE",
				fmt.Sprintf("%d", forwarded.ID),
			)
		}()
	}

	go s.broadcastMessageWithPush(context.Background(), forwarded)

	return forwarded, nil
}

func (s *ChatService) GetChatMessages(ctx context.Context, userID uuid.UUID, chatID uuid.UUID, limit, offset int) ([]chat.Message, error) {
	// PUBLIC channels are readable without a subscription; everything else needs membership.
	if err := s.requireReadAccess(ctx, chatID, userID); err != nil {
		return nil, err
	}

	return s.msgs.GetMessages(ctx, chatID, userID, limit, offset)
}

func (s *ChatService) MarkMessageAsRead(ctx context.Context, userID uuid.UUID, chatID uuid.UUID, messageID int64) {
	if err := s.members.UpdateReadCursor(ctx, userID, chatID, messageID); err != nil {
		log.Warn().Err(err).Msg("Failed to update read cursor")
		return
	}

	// Уведомляем собеседников о прочтении
	go s.broadcastToChatExcept(context.Background(), chatID, userID, socketmodel.EventMessagesRead, map[string]interface{}{
		"chat_id":      chatID,
		"last_read_id": messageID,
		"by_user":      userID,
	})
}

func (s *ChatService) GetMessageReaders(ctx context.Context, userID uuid.UUID, messageID int64) ([]chat.MessageReaderInfoResponse, error) {
	msg, err := s.msgs.GetMessageByID(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("message not found: %w", err)
	}

	if err := s.requireMembership(ctx, msg.ChatID, userID); err != nil {
		return nil, err
	}

	// Проверяем тип чата — readers только для групп
	chatType, err := s.getChatType(ctx, msg.ChatID)
	if err != nil {
		return nil, err
	}
	if chatType == chat.TypeDirect || chatType == chat.TypeSaved {
		return nil, fmt.Errorf("readers are only available for groups and channels")
	}

	readers, err := s.members.GetMessageReaders(ctx, msg.ChatID, messageID)
	if err != nil {
		return nil, fmt.Errorf("get readers: %w", err)
	}

	// Исключаем автора из списка прочитавших
	if msg.SenderID != nil {
		filtered := make([]chat.MessageReaderInfoResponse, 0, len(readers))
		for _, r := range readers {
			if r.UserID != *msg.SenderID {
				filtered = append(filtered, r)
			}
		}
		readers = filtered
	}

	return readers, nil
}