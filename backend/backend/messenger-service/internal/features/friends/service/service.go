// Package service — business logic for the Friends feature.
//
// Sends WebSocket events FRIEND_REQUEST_RECEIVED / FRIEND_REQUEST_ACCEPTED to
// affected users so their UIs update in real time.
package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	chatRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/repository"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/friends"
	friendsRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/friends/repository"
	userRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/user/repository"
	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket/hub"
)

type Service struct {
	repo      friendsRepo.Repository
	userRepo  userRepo.Repository
	blocklist chatRepo.ModerationRepo
	hub       *hub.Hub
}

func New(repo friendsRepo.Repository, userRepo userRepo.Repository, mods chatRepo.ModerationRepo, hub *hub.Hub) *Service {
	return &Service{repo: repo, userRepo: userRepo, blocklist: mods, hub: hub}
}

// SendRequest — current user wants to add target as friend.
// Domain errors: ErrSelfRequest, ErrBlocked, ErrAlreadyFriends, ErrAlreadyPending.
// If a prior REJECTED / CANCELLED request exists, it is re-opened (status -> PENDING).
func (s *Service) SendRequest(ctx context.Context, requesterID, targetID uuid.UUID) (*friends.FriendRequest, error) {
	if requesterID == targetID {
		return nil, friends.ErrSelfRequest
	}
	// Block-list check (both directions).
	if blocked, err := s.blocklist.IsBlocked(ctx, requesterID, targetID); err == nil && blocked {
		return nil, friends.ErrBlocked
	}
	if blocked, err := s.blocklist.IsBlocked(ctx, targetID, requesterID); err == nil && blocked {
		return nil, friends.ErrBlocked
	}

	if existing, err := s.repo.FindBetween(ctx, requesterID, targetID); err == nil && existing != nil {
		switch existing.Status {
		case friends.StatusAccepted:
			return nil, friends.ErrAlreadyFriends
		case friends.StatusPending:
			return nil, friends.ErrAlreadyPending
			// REJECTED / CANCELLED — fall through to Insert (ON CONFLICT) which resets row to PENDING.
		}
	}

	fr, err := s.repo.Insert(ctx, requesterID, targetID)
	if err != nil {
		return nil, err
	}

	// Notify target via WS.
	go s.notify(targetID, socketmodel.EventFriendRequestReceived, fr.ID, requesterID)

	log.Info().Str("req_id", fr.ID.String()).Str("from", requesterID.String()).Str("to", targetID.String()).Msg("friend request sent")
	return fr, nil
}

// AcceptRequest — current user accepts a PENDING request (must be target).
func (s *Service) AcceptRequest(ctx context.Context, userID, requestID uuid.UUID) error {
	fr, err := s.repo.GetByID(ctx, requestID)
	if err != nil {
		return err
	}
	if fr.TargetID != userID {
		return friends.ErrForbidden
	}
	if fr.Status != friends.StatusPending {
		return friends.ErrInvalidTransition
	}
	if err := s.repo.UpdateStatus(ctx, requestID, friends.StatusAccepted); err != nil {
		return err
	}
	// Tell the requester that we accepted.
	go s.notify(fr.RequesterID, socketmodel.EventFriendRequestAccepted, fr.ID, userID)
	log.Info().Str("req_id", fr.ID.String()).Msg("friend request accepted")
	return nil
}

// RejectRequest — current user (must be target) declines a PENDING request.
func (s *Service) RejectRequest(ctx context.Context, userID, requestID uuid.UUID) error {
	fr, err := s.repo.GetByID(ctx, requestID)
	if err != nil {
		return err
	}
	if fr.TargetID != userID {
		return friends.ErrForbidden
	}
	if fr.Status != friends.StatusPending {
		return friends.ErrInvalidTransition
	}
	if err := s.repo.UpdateStatus(ctx, requestID, friends.StatusRejected); err != nil {
		return err
	}
	go s.notify(fr.RequesterID, socketmodel.EventFriendRequestRejected, fr.ID, userID)
	return nil
}

// CancelRequest — current user (must be requester) withdraws their own PENDING request.
func (s *Service) CancelRequest(ctx context.Context, userID, requestID uuid.UUID) error {
	fr, err := s.repo.GetByID(ctx, requestID)
	if err != nil {
		return err
	}
	if fr.RequesterID != userID {
		return friends.ErrForbidden
	}
	if fr.Status != friends.StatusPending {
		return friends.ErrInvalidTransition
	}
	if err := s.repo.UpdateStatus(ctx, requestID, friends.StatusCancelled); err != nil {
		return err
	}
	go s.notify(fr.TargetID, socketmodel.EventFriendRequestCancelled, fr.ID, userID)
	return nil
}

// Unfriend — drop ACCEPTED relationship from either side.
func (s *Service) Unfriend(ctx context.Context, userID, otherID uuid.UUID) error {
	if err := s.repo.DeleteBetween(ctx, userID, otherID); err != nil {
		return err
	}
	go s.notify(otherID, socketmodel.EventFriendRemoved, uuid.Nil, userID)
	return nil
}

// GetRelation derives the relation status between current user and target.
func (s *Service) GetRelation(ctx context.Context, userID, otherID uuid.UUID) (string, *uuid.UUID, error) {
	if blocked, _ := s.blocklist.IsBlocked(ctx, userID, otherID); blocked {
		return friends.RelationBlocked, nil, nil
	}
	if blocked, _ := s.blocklist.IsBlocked(ctx, otherID, userID); blocked {
		return friends.RelationBlocked, nil, nil
	}
	fr, err := s.repo.FindBetween(ctx, userID, otherID)
	if err != nil {
		return "", nil, err
	}
	if fr == nil {
		return friends.RelationNone, nil, nil
	}
	switch fr.Status {
	case friends.StatusAccepted:
		return friends.RelationFriends, &fr.ID, nil
	case friends.StatusPending:
		if fr.RequesterID == userID {
			return friends.RelationPendingSent, &fr.ID, nil
		}
		return friends.RelationPendingReceived, &fr.ID, nil
	}
	return friends.RelationNone, nil, nil
}

// ListFriends — paginated friend list with minimal user info enrichment.
func (s *Service) ListFriends(ctx context.Context, userID uuid.UUID, limit, offset int) ([]friends.FriendInfo, error) {
	rows, err := s.repo.ListFriends(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]friends.FriendInfo, 0, len(rows))
	for _, r := range rows {
		out = append(out, friends.FriendInfo{
			UserID:        r.UserID,
			Username:      r.Username,
			DisplayName:   r.DisplayName,
			AvatarMediaID: r.AvatarMediaID,
			BecameAt:      r.BecameAt,
		})
	}
	return out, nil
}

// ListRequests — direction is "incoming" / "outgoing" (default both).
func (s *Service) ListRequests(ctx context.Context, userID uuid.UUID, direction string, limit, offset int) ([]friends.FriendRequestInfo, error) {
	rows, err := s.repo.ListRequests(ctx, userID, direction, limit, offset)
	if err != nil {
		return nil, err
	}
	out := make([]friends.FriendRequestInfo, 0, len(rows))
	for _, r := range rows {
		// Use the real per-row direction computed in SQL ("incoming"/"outgoing").
		// Previously every row was labeled "either", so the client couldn't tell an
		// incoming request (acceptable) from its own outgoing one → accept → 403.
		out = append(out, friends.FriendRequestInfo{
			RequestID:     r.RequestID,
			Direction:     r.Direction,
			Status:        r.Status,
			OtherUserID:   r.OtherUserID,
			OtherUsername: r.OtherUsername,
			OtherDisplay:  r.OtherDisplay,
			OtherAvatarID: r.OtherAvatarID,
			CreatedAt:     r.CreatedAt,
		})
	}
	return out, nil
}

// notify — best-effort WebSocket push to target user. Enriches payload with
// the actor's username / avatar (the user that performed the action).
func (s *Service) notify(toUserID uuid.UUID, event string, requestID, actorID uuid.UUID) {
	if s.hub == nil {
		return
	}
	actor, err := s.userRepo.GetUserByID(context.Background(), actorID)
	if err != nil || actor == nil {
		return
	}
	payload := friends.WSPayload{
		RequestID:    requestID,
		FromUserID:   actorID,
		FromUsername: actor.Username,
		FromDisplay:  actor.DisplayName,
		FromAvatarID: actor.AvatarMediaID,
	}
	if err := s.hub.SendToUser(toUserID, socketmodel.WSMessage{Type: event, Payload: payload}); err != nil {
		log.Debug().Err(err).Str("to", toUserID.String()).Str("event", event).Msg("ws send failed (ok if user offline)")
	}
}
