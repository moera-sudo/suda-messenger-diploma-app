package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/repository"
	grpcClient "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/grpc"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/notification"
	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket/hub"
	preferencesService 	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences/service"
	userPinsRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/userpins/repository"
)

type ChatService struct {
	chats   repository.ChatRepo
	msgs    repository.MessageRepo
	members repository.MemberRepo
	mods    repository.ModerationRepo
	prefs *preferencesService.Service
	userPinsRepo userPinsRepo.Repository
	hub      *hub.Hub
	redis    *redis.Client
	notifier notification.Service
	media    *grpcClient.MediaClient
	txClient *grpcClient.TransactionClient
}

func NewChatService(
	chats repository.ChatRepo,
	msgs repository.MessageRepo,
	members repository.MemberRepo,
	mods repository.ModerationRepo,
	prefs *preferencesService.Service,
	userPinsRepo userPinsRepo.Repository,
	hub *hub.Hub,
	redis *redis.Client,
	notifier notification.Service,
	media *grpcClient.MediaClient,
	txClient *grpcClient.TransactionClient,
) *ChatService {
	return &ChatService{
		chats:    chats,
		msgs:     msgs,
		members:  members,
		mods:     mods,
		prefs: prefs,
		userPinsRepo: userPinsRepo,
		hub:      hub,
		redis:    redis,
		notifier: notifier,
		media:    media,
		txClient: txClient,
	}
}

// * Redis Cache Helpers
func (s *ChatService) getChatType(ctx context.Context, chatID uuid.UUID) (string, error) {
	key := fmt.Sprintf("chat:%s:type", chatID)

	if cached, err := s.redis.Get(ctx, key).Result(); err == nil {
		return cached, nil
	}

	chatType, err := s.chats.GetChatType(ctx, chatID)
	if err != nil {
		return "", err
	}

	s.redis.Set(ctx, key, chatType, 0) // 0 = без TTL
	return chatType, nil
}

func (s *ChatService) getChatMembersCached(ctx context.Context, chatID uuid.UUID) ([]uuid.UUID, error) {
	key := fmt.Sprintf("chat:%s:members", chatID)

	memberStrs, err := s.redis.SMembers(ctx, key).Result()
	if err == nil && len(memberStrs) > 0 {
		ids := make([]uuid.UUID, 0, len(memberStrs))
		for _, str := range memberStrs {
			if uid, err := uuid.Parse(str); err == nil {
				ids = append(ids, uid)
			}
		}
		return ids, nil
	}

	ids, err := s.members.GetChatMembers(ctx, chatID)
	if err != nil {
		return nil, err
	}

	if len(ids) > 0 {
		pipe := s.redis.Pipeline()
		for _, uid := range ids {
			pipe.SAdd(ctx, key, uid.String())
		}
		pipe.Expire(ctx, key, time.Hour)
		pipe.Exec(ctx)
	}

	return ids, nil
}

func (s *ChatService) invalidateMembersCache(ctx context.Context, chatID uuid.UUID) {
	s.redis.Del(ctx, fmt.Sprintf("chat:%s:members", chatID))
}

func (s *ChatService) isUserOnline(ctx context.Context, userID uuid.UUID) bool {
	val, err := s.redis.Exists(ctx, fmt.Sprintf("user:%s:online", userID)).Result()
	return err == nil && val > 0
}

func (s *ChatService) isBlockedEither(ctx context.Context, userA, userB uuid.UUID) (bool, error) {
	key := fmt.Sprintf("block:%s:%s", userA, userB)

	if cached, err := s.redis.Get(ctx, key).Result(); err == nil {
		return cached == "1", nil
	}

	blockedAB, err := s.mods.IsBlocked(ctx, userA, userB)
	if err != nil {
		return false, err
	}

	blockedBA, err := s.mods.IsBlocked(ctx, userB, userA)
	if err != nil {
		return false, err
	}

	blocked := blockedAB || blockedBA

	val := "0"
	if blocked {
		val = "1"
	}
	s.redis.Set(ctx, key, val, 30*time.Minute)

	return blocked, nil
}

func (s *ChatService) invalidateBlockCache(ctx context.Context, userA, userB uuid.UUID) {
	s.redis.Del(ctx, fmt.Sprintf("block:%s:%s", userA, userB))
	s.redis.Del(ctx, fmt.Sprintf("block:%s:%s", userB, userA))
}

func (s *ChatService) sendSystemMessage(ctx context.Context, chatID uuid.UUID, action string) {
	msg := &chat.Message{
		ChatID:  chatID,
		Type:    chat.MsgTypeSystem,
		Content: action,
		Status:  chat.StatusSent,
	}

	if err := s.msgs.SaveMessage(ctx, msg); err != nil {
		log.Warn().Err(err).Str("action", action).Msg("Failed to save system message")
		return
	}

	s.broadcastToChat(ctx, chatID, socketmodel.EventNewMessage, msg)
}
