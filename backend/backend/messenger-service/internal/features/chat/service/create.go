package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel"
)

func (s *ChatService) CreateChat(ctx context.Context, userID uuid.UUID, req *chat.CreateChatReq) (*chat.Chat, error) {
	switch req.Type {
	case chat.TypeSaved:
		return s.createSavedChat(ctx, userID)
	case chat.TypeDirect:
		return s.createDirectChat(ctx, userID, req.TargetID)
	case chat.TypeGroup:
		return s.createGroupChat(ctx, userID, req.Name, req.MemberIDs)
	case chat.TypeChannel:
		return s.createChannelChat(ctx, userID, req)
	default:
		return nil, errors.New("Invalid chat type")
	}
}


func (s *ChatService) createChannelChat(ctx context.Context, userID uuid.UUID, req *chat.CreateChatReq) (*chat.Chat, error) {
	if req.Name == "" {
		return nil, fmt.Errorf("channel name is required")
	}

	visibility := channel.VisibilityPrivate
	if req.Visibility != "" {
		visibility = req.Visibility
	}

	if req.Username != "" {
		taken, err := s.chats.IsUsernameTaken(ctx, req.Username)
		if err != nil {
			return nil, fmt.Errorf("check username: %w", err)
		}
		if taken {
			return nil, fmt.Errorf("username already taken")
		}
	}

	c := &chat.Chat{
		ID:              uuid.New(),
		Type:            chat.TypeChannel,
		Name:            req.Name,
		Description:     nilIfEmpty(req.Description),
		CreatorID:       &userID,
		Visibility:      visibility,
		CommentsEnabled: true,
		Username:        nilIfEmpty(req.Username),
	}

	// Только создатель — он же OWNER, без member_ids
	if err := s.chats.CreateChat(ctx, c, []uuid.UUID{userID}, chat.RoleOwner); err != nil {
		return nil, fmt.Errorf("create channel: %w", err)
	}

	// Create channel treasury wallet in background — best-effort.
	// If transaction-service is down, channel is still created; treasury can be
	// retried via admin endpoint later (see Stage 4 plan).
	if s.txClient != nil {
		go func(channelID uuid.UUID) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			wallet, err := s.txClient.CreateWalletForChannel(ctx, channelID.String())
			if err != nil {
				log.Error().Err(err).
					Str("channel_id", channelID.String()).
					Msg("failed to create channel treasury wallet (channel created without wallet)")
				return
			}
			log.Info().
				Str("channel_id", channelID.String()).
				Str("address", wallet.Address).
				Bool("existed", wallet.Existed).
				Msg("channel treasury wallet created")
		}(c.ID)
	}

	return c, nil
}

func (s *ChatService) createSavedChat(ctx context.Context, userID uuid.UUID) (*chat.Chat, error) {
	existing, err := s.chats.GetSavedChat(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	c := &chat.Chat{
		ID:   uuid.New(),
		Type: chat.TypeSaved,
		Name: "Saved Messages",
	}

	if err := s.chats.CreateChat(ctx, c, []uuid.UUID{userID}, chat.RoleMember); err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to create saved chat")
		return nil, fmt.Errorf("create saved chat: %w", err)
	}

	log.Debug().Str("user_id", userID.String()).Str("chat_id", c.ID.String()).Msg("Saved chat created")

	return c, nil
}

func (s *ChatService) createDirectChat(ctx context.Context, userID, targetID uuid.UUID) (*chat.Chat, error) {
	// Проверяем блокировку
	blocked, err := s.isBlockedEither(ctx, userID, targetID)
	if err != nil {
		return nil, err
	}
	if blocked {
		return nil, ErrBlocked
	}

	if s.prefs != nil {
		targetPrefs, err := s.prefs.GetPreferences(ctx, targetID)
		if err == nil && targetPrefs.WhoCanMessage == preferences.VisibilityContacts {
			_, err := s.mods.GetContactName(ctx, targetID, userID)
			if err != nil {
				return nil, fmt.Errorf("forbidden: this user only accepts messages from contacts")
			}
		}
	}

	existing, err := s.chats.GetDirectChat(ctx, userID, targetID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}

	c := &chat.Chat{
		ID:   uuid.New(),
		Type: chat.TypeDirect,
	}

	if err := s.chats.CreateChat(ctx, c, []uuid.UUID{userID, targetID}, ""); err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Str("target_id", targetID.String()).Msg("Failed to create direct chat")
		return nil, fmt.Errorf("create direct chat: %w", err)
	}

	log.Debug().Str("chat_id", c.ID.String()).Str("user_id", userID.String()).Str("target_id", targetID.String()).Msg("Direct chat created")

	return c, nil
}

func (s *ChatService) createGroupChat(ctx context.Context, creatorID uuid.UUID, name string, memberIDs []uuid.UUID) (*chat.Chat, error) {
	if name == "" {
		return nil, errors.New("group name is required")
	}

	// Создатель всегда первый в списке — получит роль OWNER
	allMembers := []uuid.UUID{creatorID}
	seen := map[uuid.UUID]bool{creatorID: true}

	for _, uid := range memberIDs {
		if !seen[uid] {
			allMembers = append(allMembers, uid)
			seen[uid] = true
		}
	}

	c := &chat.Chat{
		ID:        uuid.New(),
		Type:      chat.TypeGroup,
		Name:      name,
		CreatorID: &creatorID,
	}

	if err := s.chats.CreateChat(ctx, c, allMembers, chat.RoleOwner); err != nil {
		log.Error().Err(err).Str("creator_id", creatorID.String()).Str("name", name).Msg("Failed to create group chat")
		return nil, fmt.Errorf("create group chat: %w", err)
	}

	log.Info().Str("chat_id", c.ID.String()).Str("creator_id", creatorID.String()).Str("name", name).Int("members", len(allMembers)).Msg("Group chat created")

	// Системное сообщение
	go s.sendSystemMessage(context.Background(), c.ID, chat.SystemGroupCreated)

	return c, nil
}

func (s *ChatService) UpdateChat(ctx context.Context, chatID, userID uuid.UUID, req *chat.UpdateChatReq) error {
	if err := s.canUpdateChat(ctx, chatID, userID); err != nil {
		return err
	}

	if err := s.chats.UpdateChat(ctx, chatID, req); err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Str("user_id", userID.String()).Msg("Failed to update chat")
		return fmt.Errorf("update chat: %w", err)
	}

	log.Debug().Str("chat_id", chatID.String()).Str("user_id", userID.String()).Msg("Chat updated")

	if req.Name != nil {
		go s.sendSystemMessage(context.Background(), chatID, chat.SystemGroupRenamed)
	}
	if req.AvatarMediaID != nil {
		go s.sendSystemMessage(context.Background(), chatID, chat.SystemAvatarChanged)

		// Привязываем аватарку через media-service
		if s.media != nil {
			go func() {
				_ = s.media.LinkMediaToEntity(
					context.Background(),
					req.AvatarMediaID.String(),
					"CHAT_AVATAR",
					chatID.String(),
				)
			}()
		}
	}

	go s.broadcastToChat(context.Background(), chatID, socketmodel.EventChatUpdated, map[string]interface{}{
		"chat_id": chatID,
		"update":  req,
	})

	return nil
}

func (s *ChatService) DeleteChat(ctx context.Context, chatID, userID uuid.UUID) error {
	if err := s.canDeleteChat(ctx, chatID, userID); err != nil {
		return err
	}

	// Capture members before the cascade delete so we can tell their clients to drop the chat.
	members, _ := s.getChatMembersCached(ctx, chatID)

	if err := s.chats.DeleteChat(ctx, chatID); err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Msg("Failed to delete chat")
		return fmt.Errorf("delete chat: %w", err)
	}

	// Очищаем кэш
	s.invalidateMembersCache(ctx, chatID)
	s.redis.Del(ctx, fmt.Sprintf("chat:%s:type", chatID))

	// Notify all (now former) members so their UI removes the chat in real time.
	for _, uid := range members {
		s.notifyUser(uid, socketmodel.EventChatDeleted, map[string]interface{}{
			"chat_id": chatID,
			"scope":   "everyone",
		})
	}

	log.Info().Str("chat_id", chatID.String()).Str("user_id", userID.String()).Msg("Chat deleted")

	return nil
}

// DeleteChatForMe implements "delete chat for me": clears the caller's visible
// history and hides the chat from their list. Other members keep the chat; it
// reappears for the caller (with only newer messages) once a new message arrives.
func (s *ChatService) DeleteChatForMe(ctx context.Context, chatID, userID uuid.UUID) error {
	if err := s.requireMembership(ctx, chatID, userID); err != nil {
		return err
	}

	if err := s.mods.ClearAndHideChat(ctx, userID, chatID); err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Str("user_id", userID.String()).Msg("Failed to clear/hide chat")
		return fmt.Errorf("clear and hide chat: %w", err)
	}

	// Sync the caller's other devices so they hide the chat too.
	s.notifyUser(userID, socketmodel.EventChatDeleted, map[string]interface{}{
		"chat_id": chatID,
		"scope":   "me",
	})

	log.Info().Str("chat_id", chatID.String()).Str("user_id", userID.String()).Msg("Chat deleted for user")

	return nil
}

func (s *ChatService) GetUserChats(ctx context.Context, userID uuid.UUID) ([]chat.ChatListItemResponse, error) {
	chats, err := s.chats.GetUserChats(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user chats: %w", err)
	}

	// Закреплённые чаты с датой закрепления (уже отсортированы: новые сверху).
	pins, err := s.userPinsRepo.GetPinnedChats(ctx, userID)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to load pinned chats")
		return chats, nil
	}
	if len(pins) == 0 {
		return chats, nil
	}

	pinAt := make(map[uuid.UUID]time.Time, len(pins))
	for _, p := range pins {
		pinAt[p.ChatID] = p.PinnedAt
	}

	chatMap := make(map[uuid.UUID]chat.ChatListItemResponse, len(chats))
	for _, c := range chats {
		chatMap[c.ID] = c
	}

	// Закреплённые — в порядке pins (created_at DESC, новые сверху); помечаем is_pinned + pinned_at.
	pinned := make([]chat.ChatListItemResponse, 0, len(pins))
	for _, p := range pins {
		if c, ok := chatMap[p.ChatID]; ok {
			c.IsPinned = true
			at := p.PinnedAt
			c.PinnedAt = &at
			pinned = append(pinned, c)
		}
	}

	// Остальные — в порядке БД (last_message_time DESC).
	rest := make([]chat.ChatListItemResponse, 0, len(chats))
	for _, c := range chats {
		if _, isPinned := pinAt[c.ID]; !isPinned {
			rest = append(rest, c)
		}
	}

	return append(pinned, rest...), nil
}

func nilIfEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}