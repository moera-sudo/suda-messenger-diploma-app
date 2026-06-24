package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"

	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences/repository"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket/hub"
)

const prefsCacheTTL = time.Hour

type Service struct {
	prefs  repository.PreferencesRepo
	hidden repository.LastSeenHiddenRepo
	redis  *redis.Client
	hub    *hub.Hub
}

func NewService(
	prefs repository.PreferencesRepo,
	hidden repository.LastSeenHiddenRepo,
	redis *redis.Client,
	hub *hub.Hub,
) *Service {
	return &Service{prefs: prefs, hidden: hidden, redis: redis, hub: hub}
}

// ── Preferences ─────────────────────────────────────────────

func (s *Service) GetPreferences(ctx context.Context, userID uuid.UUID) (*preferences.UserPreferences, error) {
	// 1. Redis
	if cached, err := s.getCached(ctx, userID); err == nil && cached != nil {
		return cached, nil
	}

	// 2. DB
	p, err := s.prefs.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 3. Create defaults if missing
	if p == nil {
		p = preferences.Default(userID)
		if err := s.prefs.Create(ctx, p); err != nil {
			return nil, fmt.Errorf("create default preferences: %w", err)
		}
	}

	s.setCached(ctx, userID, p)
	return p, nil
}

func (s *Service) UpdatePreferences(ctx context.Context, userID uuid.UUID, req *preferences.UpdatePreferencesReq) (*preferences.UserPreferences, error) {
	// Убеждаемся что preferences существуют
	if _, err := s.GetPreferences(ctx, userID); err != nil {
		return nil, err
	}

	if err := s.prefs.Update(ctx, userID, req); err != nil {
		return nil, fmt.Errorf("update preferences: %w", err)
	}

	s.invalidateCache(ctx, userID)

	updated, err := s.GetPreferences(ctx, userID)
	if err != nil {
		return nil, err
	}

	// WS event для мульти-девайс синхронизации
	go s.hub.SendToUser(userID, socketmodel.WSMessage{
		Type:    socketmodel.EventPreferencesUpdated,
		Payload: updated,
	})

	return updated, nil
}

// ── Last Seen Hidden ────────────────────────────────────────

func (s *Service) HideLastSeen(ctx context.Context, ownerID, hiddenID uuid.UUID) error {
	if ownerID == hiddenID {
		return fmt.Errorf("cannot hide last seen from yourself")
	}
	return s.hidden.Add(ctx, ownerID, hiddenID)
}

func (s *Service) UnhideLastSeen(ctx context.Context, ownerID, hiddenID uuid.UUID) error {
	return s.hidden.Remove(ctx, ownerID, hiddenID)
}

func (s *Service) GetHiddenList(ctx context.Context, ownerID uuid.UUID) ([]uuid.UUID, error) {
	return s.hidden.GetAll(ctx, ownerID)
}

func (s *Service) IsLastSeenHiddenFrom(ctx context.Context, ownerID, viewerID uuid.UUID) (bool, error) {
	return s.hidden.IsHidden(ctx, ownerID, viewerID)
}

// ── Cache ───────────────────────────────────────────────────

func (s *Service) getCached(ctx context.Context, userID uuid.UUID) (*preferences.UserPreferences, error) {
	key := fmt.Sprintf("prefs:%s", userID)
	data, err := s.redis.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}
	var p preferences.UserPreferences
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	return &p, nil
}

func (s *Service) setCached(ctx context.Context, userID uuid.UUID, p *preferences.UserPreferences) {
	data, err := json.Marshal(p)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to marshal preferences for cache")
		return
	}
	if err := s.redis.Set(ctx, fmt.Sprintf("prefs:%s", userID), data, prefsCacheTTL).Err(); err != nil {
		log.Warn().Err(err).Str("user_id", userID.String()).Msg("Failed to cache preferences")
	}
}

func (s *Service) invalidateCache(ctx context.Context, userID uuid.UUID) {
	if err := s.redis.Del(ctx, fmt.Sprintf("prefs:%s", userID)).Err(); err != nil {
		log.Warn().Err(err).Str("user_id", userID.String()).Msg("Failed to invalidate preferences cache")
	}
}