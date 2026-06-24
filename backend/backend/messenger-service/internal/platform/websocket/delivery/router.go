package delivery

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/service"
	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
)

type Router struct {
	chatService *service.ChatService
}

func NewRouter(chatService *service.ChatService) *Router {
	return &Router{chatService: chatService}
}

func (r *Router) HandleIncomingMessage(userID uuid.UUID, rawMsg []byte) {
	var event socketmodel.WSEvent
	if err := json.Unmarshal(rawMsg, &event); err != nil {
		log.Error().Err(err).Msg("Invalid WS event format")
		return
	}

	switch event.Type {
	case socketmodel.EventSendMessage:
		r.handleSendMessage(userID, event.Payload)
	case socketmodel.EventEditMessage:
		r.handleEditMessage(userID, event.Payload)
	case socketmodel.EventDeleteMessage:
		r.handleDeleteMessage(userID, event.Payload)
	case socketmodel.EventReadMessage:
		r.handleRead(userID, event.Payload)
	case socketmodel.EventTyping:
		r.handleTyping(userID, event.Payload)
	case socketmodel.EventForward:
		r.handleForward(userID, event.Payload)
	case socketmodel.EventReaction:
		r.handleReaction(userID, event.Payload)
	default:
		log.Warn().Str("type", event.Type).Msg("Unknown WS event type")
	}
}

func (r *Router) handleSendMessage(userID uuid.UUID, payload []byte) {
	var req chat.SendMessageReq
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Error().Err(err).Msg("Bad payload for SEND_MESSAGE")
		return
	}

	_, err := r.chatService.SendMessage(context.Background(), userID, &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send message via WS")
	}
}

func (r *Router) handleEditMessage(userID uuid.UUID, payload []byte) {
	var req struct {
		MessageID int64  `json:"message_id"`
		Content   string `json:"content"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Error().Err(err).Msg("Bad payload for EDIT_MESSAGE")
		return
	}

	err := r.chatService.EditMessage(context.Background(), userID, req.MessageID, &chat.EditMessageReq{
		Content: req.Content,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to edit message via WS")
	}
}

func (r *Router) handleDeleteMessage(userID uuid.UUID, payload []byte) {
	var req struct {
		MessageID   int64 `json:"message_id"`
		ForEveryone bool  `json:"for_everyone"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Error().Err(err).Msg("Bad payload for DELETE_MESSAGE")
		return
	}

	err := r.chatService.DeleteMessage(context.Background(), userID, req.MessageID, &chat.DeleteMessageReq{
		ForEveryone: req.ForEveryone,
	})
	if err != nil {
		log.Error().Err(err).Msg("Failed to delete message via WS")
	}
}

func (r *Router) handleTyping(userID uuid.UUID, payload []byte) {
	var req struct {
		ChatID  uuid.UUID `json:"chat_id"`
		Stopped bool      `json:"stopped"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Error().Err(err).Msg("Bad payload for TYPING")
		return
	}

	r.chatService.BroadcastTyping(userID, req.ChatID, req.Stopped)
}

func (r *Router) handleRead(userID uuid.UUID, payload []byte) {
	var req struct {
		ChatID    uuid.UUID `json:"chat_id"`
		MessageID int64     `json:"message_id"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Error().Err(err).Msg("Bad payload for READ_MESSAGE")
		return
	}

	r.chatService.MarkMessageAsRead(context.Background(), userID, req.ChatID, req.MessageID)
}

func (r *Router) handleForward(userID uuid.UUID, payload []byte) {
	var req chat.ForwardMessageReq
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Error().Err(err).Msg("Bad payload for FORWARD_MESSAGE")
		return
	}

	_, err := r.chatService.ForwardMessage(context.Background(), userID, &req)
	if err != nil {
		log.Error().Err(err).Msg("Failed to forward message via WS")
	}
}

func (r *Router) handleReaction(userID uuid.UUID, payload []byte) {
	var req struct {
		MessageID int64  `json:"message_id"`
		Emoji     string `json:"emoji"`
	}
	if err := json.Unmarshal(payload, &req); err != nil {
		log.Error().Err(err).Msg("Bad payload for REACTION")
		return
	}

	// Заготовка — реализация в следующих этапах
	log.Info().
		Str("user_id", userID.String()).
		Int64("message_id", req.MessageID).
		Str("emoji", req.Emoji).
		Msg("Reaction received (not implemented yet)")
}
