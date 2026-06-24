package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	authRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/auth/repository"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/repository"
	preferencesService "github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences/service"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/user"
	userRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/user/repository"
	grpcClient "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/grpc"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/notification"
	socketmodel "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/websocket/hub"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences"


	cfg "github.com/moera-sudo/backend/backend/messenger-service/config"
)

// ErrInvalidOldPassword is returned by ChangePassword when the supplied old password does not match.
var ErrInvalidOldPassword = errors.New("invalid old password")

type Service struct {
	repo     userRepo.Repository
	auth     authRepo.Repository
	hub      *hub.Hub
	redis    *redis.Client
	notifier notification.Service
	media    *grpcClient.MediaClient
	prefs    *preferencesService.Service
	mods     repository.ModerationRepo
	cfg      *cfg.Config
}

func NewService(
	repo userRepo.Repository,
	auth authRepo.Repository,
	hub *hub.Hub,
	redis *redis.Client,
	notifier notification.Service,
	media *grpcClient.MediaClient,
	prefs *preferencesService.Service,
	mods repository.ModerationRepo,
	cfg *cfg.Config,
) *Service {
	return &Service{
		repo:     repo,
		auth:     auth,
		hub:      hub,
		redis:    redis,
		notifier: notifier,
		media:    media,
		mods:     mods,
		prefs:    prefs,
		cfg:      cfg,
	}
}

// ListSessions returns active refresh sessions of userID for the Active Sessions UI.
// currentSessionID is the "sid" of the access token used in this request (uuid.Nil
// when unknown) and marks the matching session as IsCurrent.
func (s *Service) ListSessions(ctx context.Context, userID uuid.UUID, currentSessionID uuid.UUID) ([]user.SessionInfo, error) {
	rows, err := s.auth.ListActiveSessions(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]user.SessionInfo, 0, len(rows))
	for _, r := range rows {
		out = append(out, user.SessionInfo{
			ID:         r.ID,
			UserAgent:  r.UserAgent,
			ClientIP:   r.ClientIP,
			CreatedAt:  r.CreatedAt,
			LastUsedAt: r.LastUsedAt,
			DeviceName: r.DeviceName,
			IsCurrent:  currentSessionID != uuid.Nil && r.ID == currentSessionID,
		})
	}
	return out, nil
}

// RevokeSession revokes a single session by id (ownership check inside repo).
func (s *Service) RevokeSession(ctx context.Context, userID, sessionID uuid.UUID) error {
	return s.auth.RevokeSessionByID(ctx, userID, sessionID)
}

// RevokeOtherSessions revokes every non-revoked session except the one matching the
// supplied refresh_token. Returns count revoked.
func (s *Service) RevokeOtherSessions(ctx context.Context, userID uuid.UUID, keepRefreshToken string) (int64, error) {
	return s.auth.RevokeSessionsExceptToken(ctx, userID, keepRefreshToken)
}

// ChangePassword — verify current password via bcrypt, store new hash, revoke all other refresh sessions.
// Returns ErrInvalidOldPassword if old password does not match.
func (s *Service) ChangePassword(ctx context.Context, userID uuid.UUID, oldPassword, newPassword string) error {
	u, err := s.auth.GetUserByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("change password: lookup user: %w", err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(oldPassword)); err != nil {
		return ErrInvalidOldPassword
	}
	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("change password: hash new: %w", err)
	}
	if err := s.auth.UpdatePasswordByUserID(ctx, userID, string(newHash)); err != nil {
		return err
	}
	// Revoke every refresh session — user has to re-login everywhere with the new password.
	if err := s.auth.RevokeAllUserSessions(ctx, userID); err != nil {
		log.Warn().Err(err).Str("user_id", userID.String()).Msg("change password: failed to revoke sessions (password updated anyway)")
	}
	log.Info().Str("user_id", userID.String()).Msg("password changed")
	return nil
}

func (s *Service) GetUser(ctx context.Context, userID uuid.UUID) (*user.UserResponse, error) {
	return s.repo.GetUserByID(ctx, userID)
}

func (s *Service) UpdateProfile(ctx context.Context, userID uuid.UUID, displayName, bio, firstName, lastName string) error {
	if err := s.repo.UpdateProfile(ctx, userID, displayName, bio, firstName, lastName); err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to update profile")
		return err
	}

	log.Debug().Str("user_id", userID.String()).Msg("Profile updated")
	return nil
}

// UpdateAvatar — клиент уже загрузил аватарку через media-service,
// теперь передаёт media_id для привязки к профилю
func (s *Service) UpdateAvatar(ctx context.Context, userID uuid.UUID, mediaID uuid.UUID) error {
	if err := s.repo.UpdateAvatar(ctx, userID, mediaID); err != nil {
		return fmt.Errorf("failed to update avatar: %w", err)
	}

	// Привязываем медиа к юзеру для ACL
	if s.media != nil {
		go func() {
			err := s.media.LinkMediaToEntity(
				context.Background(),
				mediaID.String(),
				"USER_AVATAR",
				userID.String(),
			)
			if err != nil {
				log.Error().Err(err).
					Str("media_id", mediaID.String()).
					Str("user_id", userID.String()).
					Msg("Failed to link avatar media to user")
			}
		}()
	}

	return nil
}

func (s *Service) Logout(ctx context.Context, userID uuid.UUID, req *user.LogoutReq) error {
	if err := s.repo.UserLogout(ctx, userID, req.RefreshToken); err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to logout user")
		return err
	}

	if req.FCMToken != "" {
		_ = s.repo.RemoveDeviceToken(ctx, userID, req.FCMToken)
	}

	_ = s.redis.Del(ctx, fmt.Sprintf("user:%s:online", userID))

	s.hub.DisconnectUser(userID)

	log.Debug().Str("user_id", userID.String()).Msg("User logged out")

	return nil
}

func (s *Service) RegisterDevice(ctx context.Context, userID uuid.UUID, token, deviceName string) error {
	return s.notifier.SaveToken(ctx, userID, token, deviceName)
}

func (s *Service) GenerateInitData(ctx context.Context, userID uuid.UUID) (string, error) {
	u, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return "", err
	}

	initDataUser := map[string]any{
		"id":         u.ID.String(),
		"username":   strings.TrimPrefix(u.Username, "@"),
		"first_name": u.DisplayName,
		"last_name":  "",
	}

	userJSONBytes, err := json.Marshal(initDataUser)
	if err != nil {
		return "", err
	}
	userJSON := string(userJSONBytes)

	authDate := time.Now().Unix()

	params := map[string]string{
		"auth_date": strconv.FormatInt(authDate, 10),
		"user":      userJSON,
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b strings.Builder
	for i, k := range keys {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(params[k])
	}
	dataCheckString := b.String()

	secretKey := hmacSha256([]byte("WebAppData"), []byte(s.cfg.InitDataSignToken))
	hash := hmacSha256(secretKey, []byte(dataCheckString))
	hashHex := hex.EncodeToString(hash)

	final := fmt.Sprintf("auth_date=%s&user=%s&hash=%s",
		url.QueryEscape(params["auth_date"]),
		url.QueryEscape(params["user"]),
		hashHex,
	)

	return final, nil
}

func (s *Service) OnUserConnected(userID uuid.UUID) {
	log.Info().Str("user_id", userID.String()).Msg("User online")

	// Рассылаем USER_ONLINE всем собеседникам
	go s.broadcastOnlineStatus(userID, true)
}

func (s *Service) OnUserDisconnected(userID uuid.UUID) {
	if err := s.repo.UpdateLastSeen(context.Background(), userID); err != nil {
		log.Warn().Err(err).Msg("Failed to update last_seen")
	}

	log.Info().Str("user_id", userID.String()).Msg("User offline")

	// Рассылаем USER_OFFLINE всем собеседникам
	go s.broadcastOnlineStatus(userID, false)
}

func (s *Service) broadcastOnlineStatus(userID uuid.UUID, online bool) {
	interlocutors, err := s.repo.GetAllInterlocutorIDs(context.Background(), userID)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get interlocutors for online broadcast")
		return
	}

	eventType := socketmodel.EventUserOffline
	if online {
		eventType = socketmodel.EventUserOnline
	}

	event := socketmodel.WSMessage{
		Type: eventType,
		Payload: map[string]interface{}{
			"user_id": userID,
		},
	}

	for _, uid := range interlocutors {
		s.hub.SendToUser(uid, event)
	}
}

// GetUserStatus — проверяет блокировку, возвращает статус
func (s *Service) GetUserStatus(ctx context.Context, requesterID, targetID uuid.UUID) (*user.UserStatusResponse, error) {
	// 1. Проверка блока (уже есть)
	blocked, err := s.isBlockedEither(ctx, requesterID, targetID)
	if err != nil {
		return nil, err
	}
	if blocked {
		return nil, fmt.Errorf("forbidden: user is blocked")
	}

	// 2. Получаем preferences target и requester
	targetPrefs, err := s.prefs.GetPreferences(ctx, targetID)
	if err != nil {
		return nil, err
	}
	requesterPrefs, err := s.prefs.GetPreferences(ctx, requesterID)
	if err != nil {
		return nil, err
	}

	// 3. Проверка NOBODY (взаимность: если любой из них NOBODY — никто никого не видит)
	if targetPrefs.LastSeenVisibility == preferences.VisibilityNobody ||
		requesterPrefs.LastSeenVisibility == preferences.VisibilityNobody {
		return &user.UserStatusResponse{
			UserID:   targetID,
			IsOnline: false, // Скрываем даже online
		}, nil
	}

	// 4. Проверка CONTACTS — requester должен быть в контактах у target
	if targetPrefs.LastSeenVisibility == preferences.VisibilityContacts {
		_, err := s.mods.GetContactName(ctx, targetID, requesterID)
		if err != nil {
			// Не в контактах — скрываем
			return &user.UserStatusResponse{UserID: targetID, IsOnline: false}, nil
		}
	}

	// 5. Проверка индивидуального скрытия
	hidden, err := s.prefs.IsLastSeenHiddenFrom(ctx, targetID, requesterID)
	if err != nil {
		return nil, err
	}
	if hidden {
		return &user.UserStatusResponse{UserID: targetID, IsOnline: false}, nil
	}

	// 6. Всё ок — возвращаем реальный статус
	isOnline := s.isUserOnline(ctx, targetID)
	resp := &user.UserStatusResponse{
		UserID:   targetID,
		IsOnline: isOnline,
	}
	if !isOnline {
		if lastSeen, err := s.repo.GetLastSeenAt(ctx, targetID); err == nil {
			resp.LastSeenAt = lastSeen
		}
	}
	return resp, nil
}

// * Updated, with check on block for other users
func (s *Service) GetUserWithBlockCheck(ctx context.Context, requesterID, targetID uuid.UUID) (*user.UserResponse, error) {
	// Блок
	blocked, err := s.isBlockedEither(ctx, requesterID, targetID)
	if err != nil {
		return nil, err
	}
	if blocked {
		return nil, fmt.Errorf("forbidden: user is blocked")
	}

	resp, err := s.repo.GetUserByID(ctx, targetID)
	if err != nil || resp == nil {
		return resp, err
	}

	// Privacy: avatar
	targetPrefs, err := s.prefs.GetPreferences(ctx, targetID)
	if err != nil {
		return resp, nil // preferences fail — возвращаем без фильтрации
	}

	if targetPrefs.AvatarVisibility == preferences.VisibilityNobody {
		resp.AvatarMediaID = nil
	} else if targetPrefs.AvatarVisibility == preferences.VisibilityContacts {
		_, err := s.mods.GetContactName(ctx, targetID, requesterID)
		if err != nil {
			resp.AvatarMediaID = nil
		}
	}

	return resp, nil
}

// * Helpers
func (s *Service) isUserOnline(ctx context.Context, userID uuid.UUID) bool {
	val, err := s.redis.Exists(ctx, fmt.Sprintf("user:%s:online", userID)).Result()
	return err == nil && val > 0
}

func (s *Service) isBlockedEither(ctx context.Context, userA, userB uuid.UUID) (bool, error) {
	blockedAB, err := s.mods.IsBlocked(ctx, userA, userB)
	if err != nil {
		return false, err
	}
	if blockedAB {
		return true, nil
	}

	blockedBA, err := s.mods.IsBlocked(ctx, userB, userA)
	if err != nil {
		return false, err
	}

	return blockedBA, nil
}

func hmacSha256(key, msg []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}
