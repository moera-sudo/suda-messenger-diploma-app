package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
)

// ── User side: join requests ───────────────────────────────

// RequestJoin — a user asks to join a PRIVATE channel. Public channels must use
// Subscribe. Token-gating (if any) is checked fail-open. Notifies channel admins.
func (s *Service) RequestJoin(ctx context.Context, channelID, userID uuid.UUID) error {
	c, err := s.requireChannel(ctx, channelID)
	if err != nil {
		return err
	}
	if c.Visibility == channel.VisibilityPublic {
		return fmt.Errorf("public channel: use subscribe")
	}
	if isMember, _ := s.memberRepo.IsMember(ctx, channelID, userID); isMember {
		return fmt.Errorf("already a member")
	}
	// No token gating for private channels — entry is approval/invite based only.
	if err := s.joinReq.Upsert(ctx, channelID, userID, channel.JoinKindRequest, userID); err != nil {
		return fmt.Errorf("create join request: %w", err)
	}
	s.notifyAdmins(ctx, channelID, socketmodel.EventChannelJoinRequest, map[string]any{
		"channel_id": channelID, "user_id": userID,
	})
	return nil
}

// CancelJoinRequest — user withdraws their own pending request.
func (s *Service) CancelJoinRequest(ctx context.Context, channelID, userID uuid.UUID) error {
	jr, err := s.joinReq.Get(ctx, channelID, userID)
	if err != nil {
		return err
	}
	if jr == nil || jr.Kind != channel.JoinKindRequest {
		return fmt.Errorf("join request not found")
	}
	return s.joinReq.Delete(ctx, channelID, userID)
}

// ── Admin side: moderate requests ──────────────────────────

func (s *Service) ListJoinRequests(ctx context.Context, channelID, adminID uuid.UUID, limit, offset int) ([]channel.PendingRequestInfo, error) {
	if err := s.requireChannelAdmin(ctx, channelID, adminID); err != nil {
		return nil, err
	}
	return s.joinReq.ListPending(ctx, channelID, channel.JoinKindRequest, limit, offset)
}

func (s *Service) ApproveJoin(ctx context.Context, channelID, adminID, userID uuid.UUID) error {
	if err := s.requireChannelAdmin(ctx, channelID, adminID); err != nil {
		return err
	}
	jr, err := s.joinReq.Get(ctx, channelID, userID)
	if err != nil {
		return err
	}
	if jr == nil || jr.Kind != channel.JoinKindRequest || jr.Status != channel.JoinStatusPending {
		return fmt.Errorf("join request not found")
	}
	if err := s.memberRepo.AddMember(ctx, channelID, userID, channel.RoleSubscriber); err != nil {
		return fmt.Errorf("add subscriber: %w", err)
	}
	if err := s.joinReq.SetStatus(ctx, channelID, userID, channel.JoinStatusApproved, adminID); err != nil {
		return err
	}
	s.notifyUser(userID, socketmodel.EventChannelJoinApproved, map[string]any{"channel_id": channelID})
	return nil
}

func (s *Service) RejectJoin(ctx context.Context, channelID, adminID, userID uuid.UUID) error {
	if err := s.requireChannelAdmin(ctx, channelID, adminID); err != nil {
		return err
	}
	jr, err := s.joinReq.Get(ctx, channelID, userID)
	if err != nil {
		return err
	}
	if jr == nil || jr.Kind != channel.JoinKindRequest || jr.Status != channel.JoinStatusPending {
		return fmt.Errorf("join request not found")
	}
	if err := s.joinReq.SetStatus(ctx, channelID, userID, channel.JoinStatusRejected, adminID); err != nil {
		return err
	}
	s.notifyUser(userID, socketmodel.EventChannelJoinRejected, map[string]any{"channel_id": channelID})
	return nil
}

// ── Admin side: invites ────────────────────────────────────

func (s *Service) InviteUser(ctx context.Context, channelID, adminID, targetID uuid.UUID) error {
	c, err := s.requireChannel(ctx, channelID)
	if err != nil {
		return err
	}
	if err := s.requireChannelAdmin(ctx, channelID, adminID); err != nil {
		return err
	}
	if c.Visibility == channel.VisibilityPublic {
		return fmt.Errorf("public channel: users subscribe themselves")
	}
	if isMember, _ := s.memberRepo.IsMember(ctx, channelID, targetID); isMember {
		return fmt.Errorf("already a member")
	}
	if err := s.joinReq.Upsert(ctx, channelID, targetID, channel.JoinKindInvite, adminID); err != nil {
		return fmt.Errorf("create invite: %w", err)
	}
	s.notifyUser(targetID, socketmodel.EventChannelInvite, map[string]any{
		"channel_id": channelID, "invited_by": adminID,
	})
	return nil
}

func (s *Service) RevokeInvite(ctx context.Context, channelID, adminID, targetID uuid.UUID) error {
	if err := s.requireChannelAdmin(ctx, channelID, adminID); err != nil {
		return err
	}
	jr, err := s.joinReq.Get(ctx, channelID, targetID)
	if err != nil {
		return err
	}
	if jr == nil || jr.Kind != channel.JoinKindInvite {
		return fmt.Errorf("invite not found")
	}
	return s.joinReq.Delete(ctx, channelID, targetID)
}

// ── User side: invites ─────────────────────────────────────

func (s *Service) ListMyInvites(ctx context.Context, userID uuid.UUID) ([]channel.MyInviteInfo, error) {
	return s.joinReq.ListMyInvites(ctx, userID)
}

func (s *Service) AcceptInvite(ctx context.Context, channelID, userID uuid.UUID) error {
	jr, err := s.joinReq.Get(ctx, channelID, userID)
	if err != nil {
		return err
	}
	if jr == nil || jr.Kind != channel.JoinKindInvite || jr.Status != channel.JoinStatusPending {
		return fmt.Errorf("invite not found")
	}
	if err := s.memberRepo.AddMember(ctx, channelID, userID, channel.RoleSubscriber); err != nil {
		return fmt.Errorf("add subscriber: %w", err)
	}
	if err := s.joinReq.SetStatus(ctx, channelID, userID, channel.JoinStatusApproved, userID); err != nil {
		return err
	}
	s.notifyAdmins(ctx, channelID, socketmodel.EventChannelInviteAccepted, map[string]any{
		"channel_id": channelID, "user_id": userID,
	})
	return nil
}

func (s *Service) DeclineInvite(ctx context.Context, channelID, userID uuid.UUID) error {
	jr, err := s.joinReq.Get(ctx, channelID, userID)
	if err != nil {
		return err
	}
	if jr == nil || jr.Kind != channel.JoinKindInvite {
		return fmt.Errorf("invite not found")
	}
	return s.joinReq.SetStatus(ctx, channelID, userID, channel.JoinStatusRejected, userID)
}

// ── Helpers ────────────────────────────────────────────────

func (s *Service) requireChannel(ctx context.Context, channelID uuid.UUID) (*chat.Chat, error) {
	c, err := s.chatRepo.GetChatByID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if c == nil || c.Type != chat.TypeChannel {
		return nil, fmt.Errorf("channel not found")
	}
	return c, nil
}

func (s *Service) requireChannelAdmin(ctx context.Context, channelID, userID uuid.UUID) error {
	role, err := s.memberRepo.GetMemberRole(ctx, channelID, userID)
	if err != nil {
		return fmt.Errorf("forbidden: not a member")
	}
	if role != chat.RoleOwner && role != chat.RoleAdmin {
		return fmt.Errorf("forbidden: only admins can manage join requests")
	}
	return nil
}

// checkGating mirrors Subscribe's token-gating (fail-open).
func (s *Service) checkGating(ctx context.Context, channelID, userID uuid.UUID) error {
	if s.txClient == nil {
		return nil
	}
	gating, err := s.txClient.CheckTokenGating(ctx, channelID.String(), userID.String())
	if err != nil {
		log.Warn().Err(err).Str("channel_id", channelID.String()).Msg("token gating check failed, allowing (fail-open)")
		return nil
	}
	if gating.Required && !gating.Passed {
		return &channel.GatingError{
			Reason:         gating.Reason,
			MinBalanceWei:  gating.MinBalanceWei,
			UserBalanceWei: gating.UserBalanceWei,
		}
	}
	return nil
}

func (s *Service) notifyUser(userID uuid.UUID, event string, payload any) {
	if s.hub == nil {
		return
	}
	_ = s.hub.SendToUser(userID, socketmodel.WSMessage{Type: event, Payload: payload})
}

func (s *Service) notifyAdmins(ctx context.Context, channelID uuid.UUID, event string, payload any) {
	if s.hub == nil {
		return
	}
	admins, err := s.joinReq.ListChannelAdminIDs(ctx, channelID)
	if err != nil {
		return
	}
	for _, id := range admins {
		_ = s.hub.SendToUser(id, socketmodel.WSMessage{Type: event, Payload: payload})
	}
}
