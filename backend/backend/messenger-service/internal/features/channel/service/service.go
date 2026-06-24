package service

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	chatRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/repository"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel/repository"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	userRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/user/repository"
	grpcClient "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/grpc"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket/hub"
)

type Service struct {
	channels repository.ChannelRepo
	comments repository.CommentsRepo
	apps     repository.AppsRepo
	joinReq  repository.JoinRequestRepo

	chatRepo    chatRepo.ChatRepo
	memberRepo  chatRepo.MemberRepo
	userRepo    userRepo.Repository

	redis *redis.Client
	hub   *hub.Hub

	txClient *grpcClient.TransactionClient
}

func NewService(
	channels repository.ChannelRepo,
	comments repository.CommentsRepo,
	apps repository.AppsRepo,
	joinReq repository.JoinRequestRepo,
	chatRepo chatRepo.ChatRepo,
	memberRepo chatRepo.MemberRepo,
	userRepo userRepo.Repository,
	redis *redis.Client,
	hub *hub.Hub,
	txClient *grpcClient.TransactionClient,
) *Service {
	return &Service{
		channels: channels, comments: comments, apps: apps, joinReq: joinReq,
		chatRepo: chatRepo, memberRepo: memberRepo, userRepo: userRepo,
		redis: redis, hub: hub,
		txClient: txClient,
	}
}

// ResolveUserID resolves a username (with or without leading @) to a user id.
func (s *Service) ResolveUserID(ctx context.Context, username string) (uuid.UUID, error) {
	handle := strings.TrimPrefix(strings.TrimSpace(username), "@")
	if handle == "" {
		return uuid.Nil, fmt.Errorf("user not found")
	}
	u, err := s.userRepo.GetUserByExactUsername(ctx, handle)
	if err != nil || u == nil {
		return uuid.Nil, fmt.Errorf("user not found")
	}
	return u.ID, nil
}

// GetSettings returns current editable settings for the admin Settings form.
func (s *Service) GetSettings(ctx context.Context, channelID, userID uuid.UUID) (*channel.ChannelSettingsResponse, error) {
	role, err := s.memberRepo.GetMemberRole(ctx, channelID, userID)
	if err != nil {
		return nil, fmt.Errorf("forbidden: not a member")
	}
	if role != chat.RoleOwner && role != chat.RoleAdmin {
		return nil, fmt.Errorf("forbidden: only admins can view settings")
	}
	c, err := s.chatRepo.GetChatByID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if c == nil || c.Type != chat.TypeChannel {
		return nil, fmt.Errorf("channel not found")
	}
	resp := &channel.ChannelSettingsResponse{
		ID:              c.ID,
		Name:            c.Name,
		Username:        c.Username,
		AvatarMediaID:   c.AvatarMediaID,
		Visibility:      c.Visibility,
		CommentsEnabled: c.CommentsEnabled,
	}
	if c.Description != nil {
		resp.Description = *c.Description
	}
	return resp, nil
}

// ── Subscriptions ──────────────────────────────────────────

func (s *Service) Subscribe(ctx context.Context, channelID, userID uuid.UUID) error {
	chatType, err := s.chatRepo.GetChatType(ctx, channelID)
	if err != nil {
		return fmt.Errorf("get chat type: %w", err)
	}
	if chatType != chat.TypeChannel {
		return fmt.Errorf("not a channel")
	}

	// Проверяем visibility
	c, err := s.chatRepo.GetChatByID(ctx, channelID)
	if err != nil {
		return err
	}
	if c == nil {
		return fmt.Errorf("channel not found")
	}
	if c.Visibility == channel.VisibilityPrivate {
		return fmt.Errorf("forbidden: private channel, send a join request instead")
	}

	// Платная подписка (token gating, только PUBLIC-каналы). Если у канала есть
	// правило с ценой > 0 — это платный канал: списываем цену с юзера в treasury
	// канала и только при успехе добавляем в подписчики. Бесплатного пути для
	// платного канала нет, повторно не списываем (идемпотентно для участника).
	// fail-open — если transaction-service недоступен, бесплатные каналы не ломаем.
	if s.txClient != nil {
		gating, err := s.txClient.CheckTokenGating(ctx, channelID.String(), userID.String())
		if err != nil {
			log.Warn().Err(err).
				Str("channel_id", channelID.String()).
				Msg("token gating check failed, allowing free subscribe (fail-open)")
		} else if gating.Required && isPaidPrice(gating.PriceWei) {
			// Уже подписан → не списываем повторно.
			if member, _ := s.memberRepo.IsMember(ctx, channelID, userID); member {
				return nil
			}
			if gating.Reason == "wallet_not_found" {
				return &channel.GatingError{Reason: "wallet_not_found"}
			}
			// Не хватает средств на цену подписки — понятная 402 ещё до списания.
			if weiLess(gating.UserBalanceWei, gating.PriceWei) {
				return &channel.GatingError{
					Reason:         "insufficient_balance",
					MinBalanceWei:  gating.PriceWei,
					UserBalanceWei: gating.UserBalanceWei,
				}
			}
			// Списываем цену → treasury. Только при успехе подписываем.
			if _, err := s.txClient.ChargeChannelSubscription(ctx, userID.String(), channelID.String()); err != nil {
				log.Error().Err(err).
					Str("channel_id", channelID.String()).
					Str("user_id", userID.String()).
					Msg("paid subscription charge failed")
				return fmt.Errorf("subscription payment failed: %w", err)
			}
		}
	}

	// AddMember использует ON CONFLICT DO NOTHING — идемпотентно
	if err := s.memberRepo.AddMember(ctx, channelID, userID, channel.RoleSubscriber); err != nil {
		return fmt.Errorf("add subscriber: %w", err)
	}

	return nil
}

// isPaidPrice — true, если wei-строка задаёт цену подписки > 0 (платный канал).
func isPaidPrice(wei string) bool {
	if wei == "" {
		return false
	}
	n, ok := new(big.Int).SetString(wei, 10)
	return ok && n.Sign() > 0
}

// weiLess — a < b для десятичных wei-строк (пустые трактуются как 0).
func weiLess(a, b string) bool {
	x, y := big.NewInt(0), big.NewInt(0)
	if a != "" {
		if v, ok := new(big.Int).SetString(a, 10); ok {
			x = v
		}
	}
	if b != "" {
		if v, ok := new(big.Int).SetString(b, 10); ok {
			y = v
		}
	}
	return x.Cmp(y) < 0
}

func (s *Service) Unsubscribe(ctx context.Context, channelID, userID uuid.UUID) error {
	chatType, err := s.chatRepo.GetChatType(ctx, channelID)
	if err != nil {
		return err
	}
	if chatType != chat.TypeChannel {
		return fmt.Errorf("not a channel")
	}

	role, err := s.memberRepo.GetMemberRole(ctx, channelID, userID)
	if err != nil {
		return fmt.Errorf("not subscribed")
	}
	if role == chat.RoleOwner {
		return fmt.Errorf("forbidden: owner cannot unsubscribe, transfer ownership or delete channel")
	}

	return s.memberRepo.RemoveMember(ctx, channelID, userID)
}

// GetByUsername returns a ChannelView for the channel with the given username.
// Private channels are NOT hidden anymore — a non-member gets a preview (card +
// join_status) so they can request to join; posts stay gated by CanRead.
func (s *Service) GetByUsername(ctx context.Context, username string, requesterID uuid.UUID) (*channel.ChannelView, error) {
	id, err := s.channels.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	return s.GetView(ctx, id, requesterID)
}

// GetView builds the channel card + the viewer's relationship to it.
func (s *Service) GetView(ctx context.Context, channelID, requesterID uuid.UUID) (*channel.ChannelView, error) {
	c, err := s.chatRepo.GetChatByID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if c == nil || c.Type != chat.TypeChannel {
		return nil, fmt.Errorf("channel not found")
	}

	view := &channel.ChannelView{
		ID:              c.ID,
		Name:            c.Name,
		Username:        c.Username,
		AvatarMediaID:   c.AvatarMediaID,
		Visibility:      c.Visibility,
		SubscriberCount: c.SubscriberCount,
		CommentsEnabled: c.CommentsEnabled,
		CreatedAt:       c.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
		JoinStatus:      "none",
	}
	if c.Description != nil {
		view.Description = *c.Description
	}

	role, err := s.memberRepo.GetMemberRole(ctx, channelID, requesterID)
	if err == nil && role != "" {
		view.IsMember = true
		view.MyRole = role
		view.JoinStatus = "member"
	} else if jr, _ := s.joinReq.Get(ctx, channelID, requesterID); jr != nil && jr.Status == channel.JoinStatusPending {
		if jr.Kind == channel.JoinKindInvite {
			view.JoinStatus = "pending_invite"
		} else {
			view.JoinStatus = "pending_request"
		}
	}

	// Posts are readable for public channels or members.
	view.CanRead = c.Visibility == channel.VisibilityPublic || view.IsMember

	return view, nil
}

func (s *Service) GetSubscribers(ctx context.Context, channelID, requesterID uuid.UUID, limit, offset int) ([]channel.SubscriberInfo, error) {
	// Только админы могут видеть полный список
	role, err := s.memberRepo.GetMemberRole(ctx, channelID, requesterID)
	if err != nil {
		return nil, fmt.Errorf("not a member")
	}
	if role != chat.RoleOwner && role != chat.RoleAdmin {
		return nil, fmt.Errorf("forbidden: only admins can view subscribers")
	}

	return s.channels.GetSubscribers(ctx, channelID, limit, offset)
}

// ── Settings ───────────────────────────────────────────────

func (s *Service) UpdateSettings(ctx context.Context, channelID, userID uuid.UUID, req *channel.UpdateChannelSettingsReq) error {
	role, err := s.memberRepo.GetMemberRole(ctx, channelID, userID)
	if err != nil {
		return fmt.Errorf("not a member")
	}
	if role != chat.RoleOwner && role != chat.RoleAdmin {
		return fmt.Errorf("forbidden: only admins can update settings")
	}

	// Проверка username на уникальность
	if req.Username != nil {
		taken, err := s.channels.IsUsernameTaken(ctx, *req.Username, &channelID)
		if err != nil {
			return err
		}
		if taken {
			return fmt.Errorf("username already taken")
		}
	}

	return s.channels.UpdateSettings(ctx, channelID, req)
}