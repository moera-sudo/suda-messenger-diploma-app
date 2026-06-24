package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)

func (s *Service) AddComment(ctx context.Context, channelID uuid.UUID, postID int64, senderID uuid.UUID, req *channel.AddCommentReq) (*channel.ChannelComment, error) {
	// Канал должен иметь включённые комментарии
	c, err := s.chatRepo.GetChatByID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if !c.CommentsEnabled {
		return nil, fmt.Errorf("forbidden: comments are disabled in this channel")
	}

	// Подписчик может комментировать
	isMember, err := s.memberRepo.IsMember(ctx, channelID, senderID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, fmt.Errorf("forbidden: subscribe to comment")
	}

	comment := &channel.ChannelComment{
		ChannelID:        channelID,
		PostID:           postID,
		SenderID:         senderID,
		Content:          req.Content,
		ReplyToCommentID: req.ReplyToCommentID,
	}

	if err := s.comments.Create(ctx, comment); err != nil {
		return nil, fmt.Errorf("create comment: %w", err)
	}

	// WS broadcast — всем подписчикам канала
	go s.broadcastToChannel(ctx, channelID, socketmodel.WSMessage{
		Type:    socketmodel.EventChannelNewComment,
		Payload: comment,
	})

	return comment, nil
}

func (s *Service) EditComment(ctx context.Context, commentID int64, userID uuid.UUID, content string) error {
	c, err := s.comments.GetByID(ctx, commentID)
	if err != nil {
		return err
	}
	if c == nil {
		return fmt.Errorf("comment not found")
	}
	if c.SenderID != userID {
		return fmt.Errorf("forbidden: can only edit own comments")
	}

	return s.comments.Update(ctx, commentID, content)
}

func (s *Service) DeleteComment(ctx context.Context, commentID int64, userID uuid.UUID) error {
	c, err := s.comments.GetByID(ctx, commentID)
	if err != nil {
		return err
	}
	if c == nil {
		return fmt.Errorf("comment not found")
	}

	// Автор может удалять свой комментарий в любом случае
	if c.SenderID != userID {
		// Иначе — только OWNER/ADMIN того же канала
		isMember, err := s.memberRepo.IsMember(ctx, c.ChannelID, userID)
		if err != nil {
			return err
		}
		if !isMember {
			return fmt.Errorf("forbidden: not a member of this channel")
		}
		role, err := s.memberRepo.GetMemberRole(ctx, c.ChannelID, userID)
		if err != nil {
			return err
		}
		if role != chat.RoleOwner && role != chat.RoleAdmin {
			return fmt.Errorf("forbidden: only author or admin can delete")
		}
	}

	return s.comments.SoftDelete(ctx, commentID)
}

func (s *Service) GetComments(ctx context.Context, channelID, requesterID uuid.UUID, postID int64, limit, offset int) ([]channel.ChannelComment, int, error) {
	// Приватные каналы — только подписчикам
	c, err := s.chatRepo.GetChatByID(ctx, channelID)
	if err != nil {
		return nil, 0, err
	}
	if c == nil {
		return nil, 0, fmt.Errorf("channel not found")
	}
	if c.Visibility != channel.VisibilityPublic {
		// Private channel — members only.
		isMember, err := s.memberRepo.IsMember(ctx, channelID, requesterID)
		if err != nil {
			return nil, 0, err
		}
		if !isMember {
			return nil, 0, fmt.Errorf("forbidden: subscribe to view comments")
		}
	} else {
		// Public channel — members always; non-members must pass token gating
		// (if any). Keeps a gated public channel's content hidden until qualified.
		isMember, err := s.memberRepo.IsMember(ctx, channelID, requesterID)
		if err != nil {
			return nil, 0, err
		}
		if !isMember {
			if err := s.checkGating(ctx, channelID, requesterID); err != nil {
				return nil, 0, err
			}
		}
	}

	comments, err := s.comments.GetByPost(ctx, postID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.comments.CountByPost(ctx, postID)
	if err != nil {
		return comments, 0, err
	}

	return comments, total, nil
}

// ── Helpers ────────────────────────────────────────────────

func (s *Service) broadcastToChannel(ctx context.Context, channelID uuid.UUID, msg socketmodel.WSMessage) {
	members, err := s.memberRepo.GetChatMembers(ctx, channelID)
	if err != nil {
		log.Error().Err(err).Str("channel_id", channelID.String()).Msg("Failed to load channel members for broadcast")
		return
	}
	for _, uid := range members {
		s.hub.SendToUser(uid, msg)
	}
}