package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
)

func (s *Service) LinkApp(ctx context.Context, channelID, userID, appID uuid.UUID) error {
	role, err := s.memberRepo.GetMemberRole(ctx, channelID, userID)
	if err != nil {
		return fmt.Errorf("not a member")
	}
	if role != chat.RoleOwner && role != chat.RoleAdmin {
		return fmt.Errorf("forbidden: only admins can link apps")
	}

	return s.apps.Link(ctx, &channel.ChannelAppLink{
		ChannelID: channelID,
		AppID:     appID,
		AddedBy:   userID,
	})
}

func (s *Service) UnlinkApp(ctx context.Context, channelID, userID, appID uuid.UUID) error {
	role, err := s.memberRepo.GetMemberRole(ctx, channelID, userID)
	if err != nil {
		return fmt.Errorf("not a member")
	}
	if role != chat.RoleOwner && role != chat.RoleAdmin {
		return fmt.Errorf("forbidden: only admins can unlink apps")
	}

	return s.apps.Unlink(ctx, channelID, appID)
}

func (s *Service) GetLinkedApps(ctx context.Context, channelID, requesterID uuid.UUID) ([]channel.ChannelAppLink, error) {
	c, err := s.chatRepo.GetChatByID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, fmt.Errorf("channel not found")
	}
	if c.Visibility != channel.VisibilityPublic {
		isMember, err := s.memberRepo.IsMember(ctx, channelID, requesterID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, fmt.Errorf("channel not found")
		}
	}
	return s.apps.GetByChannel(ctx, channelID)
}