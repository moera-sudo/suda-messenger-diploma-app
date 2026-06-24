package hub

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/olahol/melody"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// onlineTTL — how long the `user:<id>:online` flag survives without a refresh.
// It is refreshed on every inbound WS message and on every pong, so it only
// expires after the session is actually gone.
const onlineTTL = 10 * time.Minute

type Hub struct {
	m            *melody.Melody
	userSessions map[uuid.UUID][]*melody.Session
	mu           sync.RWMutex
	onMessage    func(userID uuid.UUID, msg []byte)
	onConnect    func(userID uuid.UUID)
	onDisconnect func(userID uuid.UUID)
	redis        *redis.Client
	jwtSecret    string
}

func (h *Hub) SetRouter(handler func(userID uuid.UUID, msg []byte)) {
	h.onMessage = handler
}

func (h *Hub) SetOnConnect(handler func(userID uuid.UUID)) {
	h.onConnect = handler
}

func (h *Hub) SetOnDisconnect(handler func(userID uuid.UUID)) {
	h.onDisconnect = handler
}

func NewHub(r *redis.Client, jwtSecret string) *Hub {
	m := melody.New()
	m.Config.MaxMessageSize = 1024 * 1024

	hub := &Hub{
		m:            m,
		userSessions: make(map[uuid.UUID][]*melody.Session),
		redis:        r,
		jwtSecret:    jwtSecret,
	}

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		if uid, ok := s.Get("user_id"); ok {
			// Refresh the online flag on any inbound activity so it never
			// expires while a live session exists.
			hub.touchOnline(uid.(uuid.UUID))
			if hub.onMessage != nil {
				hub.onMessage(uid.(uuid.UUID), msg)
			}
		}
	})

	// melody pings clients periodically (~54s); the pong response keeps the
	// online flag fresh even when the user is idle (not sending messages).
	m.HandlePong(func(s *melody.Session) {
		if uid, ok := s.Get("user_id"); ok {
			hub.touchOnline(uid.(uuid.UUID))
		}
	})

	m.HandleConnect(func(s *melody.Session) {
		// Auth source priority:
		//  1) X-User-ID header — set by api-gateway after JWT validation (server-to-server flow).
		//  2) `?token=<jwt>` query — for browser / WebView clients that cannot send custom headers
		//     on the WebSocket handshake. We validate the JWT locally with our shared secret.
		userIDStr := s.Request.Header.Get("X-User-ID")

		if userIDStr == "" {
			tokenStr := s.Request.URL.Query().Get("token")
			if tokenStr == "" {
				s.CloseWithMsg(melody.FormatCloseMessage(401, "Unauthorized"))
				return
			}
			parsed, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return []byte(hub.jwtSecret), nil
			})
			if err != nil || !parsed.Valid {
				log.Warn().Err(err).Msg("WS: invalid JWT in query")
				s.CloseWithMsg(melody.FormatCloseMessage(401, "Unauthorized"))
				return
			}
			claims, ok := parsed.Claims.(jwt.MapClaims)
			if !ok {
				s.CloseWithMsg(melody.FormatCloseMessage(401, "Unauthorized"))
				return
			}
			if tt, _ := claims["type"].(string); tt != "" && tt != "access" {
				s.CloseWithMsg(melody.FormatCloseMessage(401, "Wrong token type"))
				return
			}
			if v, ok := claims["user_id"].(string); ok {
				userIDStr = v
			} else if v, ok := claims["sub"].(string); ok {
				userIDStr = v
			}
			if userIDStr == "" {
				s.CloseWithMsg(melody.FormatCloseMessage(401, "Missing user_id in token"))
				return
			}
		}

		uid, err := uuid.Parse(userIDStr)
		if err != nil {
			s.Close()
			return
		}

		s.Set("user_id", uid)

		hub.registerSession(uid, s)

		// Ставим онлайн в Redis
		hub.touchOnline(uid)

		log.Info().Str("user_id", uid.String()).Msg("WS connected")

		// Callback для user feature (USER_ONLINE event, etc.)
		if hub.onConnect != nil {
			go hub.onConnect(uid)
		}
	})

	m.HandleDisconnect(func(s *melody.Session) {
		if uid, ok := s.Get("user_id"); ok {
			userID := uid.(uuid.UUID)
			hub.unregisterSession(userID, s)

			// Проверяем остались ли сессии у этого юзера
			hub.mu.RLock()
			remaining := len(hub.userSessions[userID])
			hub.mu.RUnlock()

			if remaining == 0 {
				// Последняя сессия закрылась — юзер оффлайн
				r.Del(context.Background(), fmt.Sprintf("user:%s:online", userID))

				log.Info().Str("user_id", userID.String()).Msg("WS disconnected (last session)")

				// Callback для user feature (UpdateLastSeen, USER_OFFLINE event)
				if hub.onDisconnect != nil {
					go hub.onDisconnect(userID)
				}
			} else {
				log.Info().
					Str("user_id", userID.String()).
					Int("remaining_sessions", remaining).
					Msg("WS disconnected (other sessions active)")
			}
		}
	})

	return hub
}

func (h *Hub) HandleRequest(c echo.Context) error {
	return h.m.HandleRequest(c.Response(), c.Request())
}

// touchOnline (re)sets the online flag for a user with a fresh TTL.
func (h *Hub) touchOnline(uid uuid.UUID) {
	h.redis.Set(context.Background(), fmt.Sprintf("user:%s:online", uid), "1", onlineTTL)
}

func (h *Hub) registerSession(uid uuid.UUID, s *melody.Session) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.userSessions[uid] = append(h.userSessions[uid], s)
}

func (h *Hub) unregisterSession(uid uuid.UUID, s *melody.Session) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sessions := h.userSessions[uid]
	for i, session := range sessions {
		if session == s {
			h.userSessions[uid] = append(sessions[:i], sessions[i+1:]...)
			break
		}
	}

	if len(h.userSessions[uid]) == 0 {
		delete(h.userSessions, uid)
	}
}

func (h *Hub) SendToUser(userID uuid.UUID, msg interface{}) error {
	h.mu.RLock()
	sessions, ok := h.userSessions[userID]
	h.mu.RUnlock()

	if !ok {
		return errors.New("user not connected")
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	for _, s := range sessions {
		s.Write(jsonMsg)
	}
	return nil
}

func (h *Hub) Broadcast(msg interface{}) {
	jsonMsg, _ := json.Marshal(msg)
	h.m.Broadcast(jsonMsg)
}

func (h *Hub) DisconnectUser(userID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	sessions, ok := h.userSessions[userID]
	if !ok {
		return
	}

	for _, s := range sessions {
		s.CloseWithMsg(melody.FormatCloseMessage(1000, "Logged out"))
	}

	delete(h.userSessions, userID)
}