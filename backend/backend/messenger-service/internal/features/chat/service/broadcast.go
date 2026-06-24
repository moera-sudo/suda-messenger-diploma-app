package service

import (
	"context"


	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)

// broadcastToChat — отправляет WS событие всем участникам чата
func (s *ChatService) broadcastToChat(ctx context.Context, chatID uuid.UUID, eventType string, payload interface{}) {
	memberIDs, err := s.getChatMembersCached(ctx, chatID)
	if err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Msg("Failed to get members for broadcast")
		return
	}

	event := socketmodel.WSMessage{
		Type:    eventType,
		Payload: payload,
	}

	for _, uid := range memberIDs {
		if err := s.hub.SendToUser(uid, event); err != nil {
			log.Debug().Err(err).Str("chat_id", chatID.String()).Str("user_id", uid.String()).Str("event", eventType).Msg("WS delivery skipped (user offline)")
		}
	}
}

// broadcastToChatExcept — отправляет всем кроме указанного пользователя
func (s *ChatService) broadcastToChatExcept(ctx context.Context, chatID uuid.UUID, exceptUID uuid.UUID, eventType string, payload interface{}) {
	memberIDs, err := s.getChatMembersCached(ctx, chatID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get members for broadcast")
		return
	}

	event := socketmodel.WSMessage{
		Type:    eventType,
		Payload: payload,
	}

	for _, uid := range memberIDs {
		if uid != exceptUID {
			if err := s.hub.SendToUser(uid, event); err != nil {
				log.Debug().Err(err).Str("chat_id", chatID.String()).Str("user_id", uid.String()).Str("event", eventType).Msg("WS delivery skipped (user offline)")
			}
		}
	}
}

// broadcastMessageWithPush — отправляет сообщение онлайн по WS, оффлайн по push
func (s *ChatService) broadcastMessageWithPush(ctx context.Context, msg *chat.Message) {
	memberIDs, err := s.getChatMembersCached(ctx, msg.ChatID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get members for broadcast")
		return
	}

	event := socketmodel.WSMessage{
		Type:    socketmodel.EventNewMessage,
		Payload: msg,
	}

	var offlineUserIDs []uuid.UUID

	for _, uid := range memberIDs {
		// The hub's in-memory session map is the source of truth for live
		// delivery. The Redis `user:<id>:online` flag has a TTL and may be
		// stale (it was the cause of "messages only appear after reload"),
		// so it must NOT gate WS delivery. We try to deliver over WS first
		// and only fall back to push when the user has no live session.
		if err := s.hub.SendToUser(uid, event); err == nil {
			continue
		}
		// Not connected — queue for push (never push to the sender).
		if msg.SenderID == nil || uid != *msg.SenderID {
			offlineUserIDs = append(offlineUserIDs, uid)
		}
	}

	if len(offlineUserIDs) > 0 && s.notifier != nil {
		go func() {
			var withPreview []uuid.UUID
			var withoutPreview []uuid.UUID

			for _, uid := range offlineUserIDs {
				// Проверяем мут
				muted, _ := s.mods.IsChatMuted(ctx, uid, msg.ChatID)
				if muted {
					continue
				}

				// Проверяем preview_in_notifications
				if s.prefs != nil {
					prefs, err := s.prefs.GetPreferences(ctx, uid)
					if err == nil && prefs != nil && !prefs.PreviewInNotifications {
						withoutPreview = append(withoutPreview, uid)
						continue
					}
				}

				withPreview = append(withPreview, uid)
			}

			pushContent := msg.Content
			if pushContent == "" && msg.AttachmentMediaID != nil {
				pushContent = "Sent an attachment"
			}
			if msg.Type == chat.MsgTypeSystem {
				pushContent = "Chat updated"
			}

			data := map[string]string{
				"type":    "NEW_MESSAGE",
				"chat_id": msg.ChatID.String(),
			}

			// Группа с превью — видят текст сообщения
			if len(withPreview) > 0 {
				if err := s.notifier.SendPush(context.Background(), withPreview, "New message", pushContent, data); err != nil {
					log.Error().Err(err).Msg("Failed to send push (with preview)")
				}
			}

			// Группа без превью — видят только "New message"
			if len(withoutPreview) > 0 {
				if err := s.notifier.SendPush(context.Background(), withoutPreview, "New message", "New message", data); err != nil {
					log.Error().Err(err).Msg("Failed to send push (no preview)")
				}
			}
		}()
	}
}

// notifyUser — отправляет WS событие конкретному пользователю
func (s *ChatService) notifyUser(userID uuid.UUID, eventType string, payload interface{}) {
	if s.hub == nil {
		return
	}
	event := socketmodel.WSMessage{
		Type:    eventType,
		Payload: payload,
	}
	s.hub.SendToUser(userID, event)
}