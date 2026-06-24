package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/auth"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/auth/config"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/auth/repository"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	chatSrv "github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/service"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/jwt"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/platform/email/sender"

	grpcclient "github.com/moera-sudo/backend/backend/messenger-service/internal/platform/grpc"
	userRepo "github.com/moera-sudo/backend/backend/messenger-service/internal/features/user/repository"
)

// ErrInvalidInitData is returned when the init_data signature does not match.
var ErrInvalidInitData = errors.New("invalid init_data signature")

// ErrInitDataExpired is returned when auth_date is older than initDataMaxAge.
var ErrInitDataExpired = errors.New("init_data expired")

// initDataMaxAge — how long an initData string stays valid for /webapp-session
// exchange. 30 minutes gives the WebView a comfortable buffer for the user to
// browse around before the underlying signature goes stale.
const initDataMaxAge = 30 * time.Minute

type Service struct {
	repo repository.Repository
	tokenManager *jwt.TokenManager
	chatSrv *chatSrv.ChatService
	emailSender *sender.Sender
	cfg *config.AuthConfig
	userRepo userRepo.Repository
	txClient *grpcclient.TransactionClient
	// initDataSignToken is the shared secret used for HMAC validation of
	// Mini App init_data (same value used to sign it in user.Service.GenerateInitData).
	initDataSignToken string
}

func NewService(repo repository.Repository, email *sender.Sender, chatSrv *chatSrv.ChatService, cfg *config.AuthConfig, userRepo userRepo.Repository, txClient *grpcclient.TransactionClient, initDataSignToken string) *Service {
	return &Service{
		repo: repo,
		chatSrv: chatSrv,
		tokenManager: jwt.NewManager(cfg.JWTSecret),
		emailSender: email,
		cfg: cfg,
		userRepo: userRepo,
		txClient: txClient,
		initDataSignToken: initDataSignToken,
	}
}


func (s *Service) Register(ctx context.Context, req *auth.RegisterReq) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &auth.User{
		Username: req.Username,
		Email : req.Email,
		DisplayName: req.DisplayName,
		PasswordHash: string(passHash),
	}

	_, err = s.repo.CreateUser(ctx, user)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to create user")
		return err
	}

	code := generateRandomCode()
	log.Debug().Str("code", code).Msg("Code created")
	// Хэширование кода верификации
	codeHash, _ := bcrypt.GenerateFromPassword([]byte(code), bcrypt.MinCost)

	if err := s.repo.CreateVerification(ctx, req.Email, string(codeHash), "REGISTRATION"); err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to create verification")
		return err
	}

	log.Info().Str("email", req.Email).Str("username", req.Username).Msg("User registered")

	// Отправка письма(тестово в горутине)
	return s.emailSender.SendVerificationCode(req.Email, code, "VERIFY_EMAIL")
}

func (s *Service) Verify(ctx context.Context, req *auth.VerifyReq) error {
    hash, userID, err := s.repo.GetVerification(ctx, req.Email, "REGISTRATION")
    if err != nil {
        return err
    }
    if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Code)); err != nil {
        log.Warn().Str("email", req.Email).Msg("Invalid verification code")
        return errors.New("Invalid code")
    }

    // 1. Создаём SavedMessages чат (best-effort)
    c := &chat.CreateChatReq{
        Type: "SAVED",
    }
    if _, err := s.chatSrv.CreateChat(ctx, userID, c); err != nil {
        log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to create saved chat on verification")
        // Не возвращаем ошибку — savedchat не критичен для верификации.
    }

    // 2. Помечаем юзера как verified в БД.
    if err := s.repo.MarkUserVerified(ctx, userID, req.Email); err != nil {
        log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to mark user verified")
        return err
    }

    log.Info().Str("email", req.Email).Str("user_id", userID.String()).Msg("Email verified")

    // 3. Создаём кошелёк в transaction-service + автоматически получаем welcome bonus 100 SUDA.
    //
    // Best-effort (Вариант Б): если transaction-service лёг или ошибся —
    // verify всё равно прошёл, юзер верифицирован. В UI он увидит «Кошелёк не активирован»
    // и сможет повторить вручную (endpoint, который дёрнет CreateWalletForUser ещё раз —
    // он идемпотентен).
    if s.txClient != nil {
        s.createWalletForVerifiedUser(ctx, userID)
    } else {
        log.Warn().Str("user_id", userID.String()).Msg("Transaction client not configured, skipping wallet creation")
    }

    return nil
}


func (s *Service) Login(ctx context.Context, req *auth.LoginReq) (*auth.TokenResponse, error) {
	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		log.Warn().Str("email", req.Email).Msg("Login attempt: user not found")
		return nil, errors.New("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		log.Warn().Str("email", req.Email).Msg("Login attempt: invalid password")
		return nil, errors.New("Invalid credentials")
	}

	if !user.IsVerified {
		log.Warn().Str("email", req.Email).Msg("Login attempt: email not verified")
		return nil, errors.New("Email not verified")
	}

	tokens, err := s.createSession(ctx, user.ID, req.UserAgent, req.ClientIP)
	if err != nil {
		log.Error().Err(err).Str("user_id", user.ID.String()).Msg("Failed to create session")
		return nil, err
	}

	log.Info().Str("user_id", user.ID.String()).Str("email", req.Email).Msg("User logged in")

	return tokens, nil
}


func (s *Service) RefreshSessionTokens(ctx context.Context, oldToken string, userAgent string, clientIP string) (*auth.TokenResponse, error) {

	userID, err := s.tokenManager.ParseRefreshToken(oldToken)
	if err != nil {
		log.Warn().Msg("Refresh attempt: invalid refresh token")
		return nil, errors.New("Invalid refresh token")
	}

	// New session id for the rotated session row (RefreshSession revokes the old
	// row and inserts a fresh one — see authRepo.RefreshSession).
	sessionID := uuid.New()
	tokens, err := s.generateTokens(userID, sessionID)

	if err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to generate tokens on refresh")
		return nil, err
	}

	if err := s.repo.RefreshSession(ctx, sessionID, oldToken, tokens.RefreshToken, time.Now().Add(tokens.RefreshExpiresIn), userAgent, clientIP); err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("Failed to refresh session in DB")
		return nil, err
	}

	log.Debug().Str("user_id", userID.String()).Msg("Session tokens refreshed")

	return &auth.TokenResponse{
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		AccessExpiresIn: tokens.AccessExpiresIn,
		RefreshExpiresIn: tokens.RefreshExpiresIn,
	}, nil
}



func (s *Service) createSession(ctx context.Context, userID uuid.UUID, userAgent string, clientIP string) (*auth.TokenResponse, error) {

	sessionID := uuid.New()
	tokens, err := s.generateTokens(userID, sessionID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SaveSession(ctx, sessionID, userID, tokens.RefreshToken, time.Now().Add(tokens.RefreshExpiresIn), userAgent, clientIP); err != nil {
		return nil, err
	}

	return &auth.TokenResponse{  
		AccessToken: tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		AccessExpiresIn: tokens.AccessExpiresIn,
		RefreshExpiresIn: tokens.RefreshExpiresIn,
	}, nil
}

func (s *Service) ForgotPassword(ctx context.Context, req *auth.ForgotPasswordReq) error {

	user, err := s.repo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to get user for password reset")
		return err
	}
	if user == nil {
		return errors.New("User not found")
	}

	code := generateRandomCode()
	codeHash, err := bcrypt.GenerateFromPassword([]byte(code), bcrypt.MinCost)

	if err := s.repo.CreateVerification(ctx, req.Email, string(codeHash), "PASSWORD_RESET"); err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to create password reset verification")
		return err
	}

	log.Debug().Str("email", req.Email).Msg("Password reset requested")

	return s.emailSender.SendVerificationCode(req.Email, code, "CHANGE_PASSWORD")
}

func (s *Service) ResetPassword(ctx context.Context, req *auth.ResetPasswordReq) error {
	hash, userID, err := s.repo.GetVerification(ctx, req.Email, "PASSWORD_RESET")
	if err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to get verification for password reset")
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Code)); err != nil {
		log.Warn().Str("email", req.Email).Msg("Invalid password reset code")
		return errors.New("Invalid code")
	}

	newPassHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.repo.UpdatePassword(ctx, req.Email, string(newPassHash)); err != nil {
		log.Error().Err(err).Str("email", req.Email).Msg("Failed to update password")
		return err
	}

	if err := s.repo.UseVerificationCode(ctx, req.Email, hash); err != nil {
		return nil
	}

	log.Info().Str("user_id", userID.String()).Str("email", req.Email).Msg("Password reset successfully")

	return s.repo.RevokeAllUserSessions(ctx, userID)

}



// generateTokens mints an access/refresh pair bound to sessionID via the "sid"
// claim, so the same id can be persisted as the session row id and later matched
// for is_current in the Active Sessions list.
func (s *Service) generateTokens(userID uuid.UUID, sessionID uuid.UUID) (*auth.TokenResponse, error) {
	accessTTL := s.cfg.AccessTokenTTL
	refreshTTL := s.cfg.RefreshTokenTTL

	accessToken, err := s.tokenManager.NewJWTWithSession(userID, "user", "access", sessionID, accessTTL)

	if err != nil {
		return nil, err
	}

	refreshToken, err := s.tokenManager.NewJWTWithSession(userID, "user", "refresh", sessionID, refreshTTL)
	if err != nil {
		return nil, err
	}

	return &auth.TokenResponse{  
		AccessToken: accessToken,
		RefreshToken: refreshToken,
		AccessExpiresIn: int64(accessTTL.Seconds()),
		RefreshExpiresIn: refreshTTL,
	}, nil
}

func generateRandomCode() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	return fmt.Sprintf("%06d", n.Int64()+100000)
}

func (s *Service) createWalletForVerifiedUser(ctx context.Context, userID uuid.UUID) {
    wallet, err := s.txClient.CreateWalletForUser(ctx, userID.String())
    if err != nil {
        log.Error().Err(err).Str("user_id", userID.String()).
            Msg("Failed to create wallet (user verified but wallet not initialized)")
        return
    }

    if err := s.userRepo.SetWalletAddress(ctx, userID, wallet.Address); err != nil {
        log.Error().Err(err).
            Str("user_id", userID.String()).
            Str("address", wallet.Address).
            Msg("Failed to save wallet_address (wallet exists in transaction-service but not in messenger_users)")
        return
    }

    log.Info().
        Str("user_id", userID.String()).
        Str("address", wallet.Address).
        Bool("existed", wallet.Existed).
        Msg("User wallet created and linked")
}

// ExchangeInitData accepts a Telegram-style init_data string produced by
// /user/init-data and returns a short-lived access JWT for the same user.
//
// Used by WebView mini-apps (Wallet, Marketplace): on load the SPA reads
// `?initData=...` from its URL and POSTs it here; the returned JWT is then
// used for every backend call from inside the WebView.
func (s *Service) ExchangeInitData(ctx context.Context, initData string) (*auth.WebappSessionResponse, error) {
	if s.initDataSignToken == "" {
		log.Error().Msg("ExchangeInitData called but INIT_DATA_SIGN_TOKEN is empty")
		return nil, fmt.Errorf("server misconfigured: init data secret missing")
	}

	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, ErrInvalidInitData
	}
	hashHex := values.Get("hash")
	authDateStr := values.Get("auth_date")
	if hashHex == "" || authDateStr == "" {
		return nil, ErrInvalidInitData
	}

	// 1. Age check.
	authDateUnix, err := strconv.ParseInt(authDateStr, 10, 64)
	if err != nil {
		return nil, ErrInvalidInitData
	}
	if time.Since(time.Unix(authDateUnix, 0)) > initDataMaxAge {
		return nil, ErrInitDataExpired
	}

	// 2. Rebuild dataCheckString — every key=value except `hash`, sorted alphabetically by key, joined by '\n'.
	keys := make([]string, 0, len(values))
	for k := range values {
		if k == "hash" {
			continue
		}
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
		b.WriteString(values.Get(k))
	}
	dataCheckString := b.String()

	// 3. HMAC verification (mirror of user.Service.GenerateInitData).
	secretKey := hmacSha256([]byte("WebAppData"), []byte(s.initDataSignToken))
	expected := hex.EncodeToString(hmacSha256(secretKey, []byte(dataCheckString)))
	if !hmac.Equal([]byte(expected), []byte(hashHex)) {
		return nil, ErrInvalidInitData
	}

	// 4. Extract user_id from the embedded user JSON.
	userJSON := values.Get("user")
	if userJSON == "" {
		return nil, ErrInvalidInitData
	}
	var userPayload struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
	if err := json.Unmarshal([]byte(userJSON), &userPayload); err != nil {
		return nil, ErrInvalidInitData
	}
	userID, err := uuid.Parse(userPayload.ID)
	if err != nil {
		return nil, ErrInvalidInitData
	}

	// 5. Sanity-check the user still exists. Repo returns (nil, nil) when missing
	// per its current contract (see authRepo.GetUserByID).
	u, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("webapp-session: lookup user: %w", err)
	}
	if u == nil {
		return nil, ErrInvalidInitData
	}

	// 6. Mint short-lived access JWT (same TTL as native login access token).
	role := u.Role
	if role == "" {
		role = "USER"
	}
	token, err := s.tokenManager.NewJWT(userID, role, "access", s.cfg.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("webapp-session: sign jwt: %w", err)
	}

	avatarMediaID := ""
	if u.AvatarMediaID != nil {
		avatarMediaID = u.AvatarMediaID.String()
	}

	log.Info().Str("user_id", userID.String()).Msg("webapp-session: token issued")
	return &auth.WebappSessionResponse{
		AccessToken:   token,
		ExpiresIn:     int64(s.cfg.AccessTokenTTL.Seconds()),
		UserID:        userID.String(),
		Username:      u.Username,
		AvatarMediaID: avatarMediaID,
	}, nil
}

// hmacSha256 — internal helper duplicated from user.Service (HMAC-SHA256 keyed by `key`).
func hmacSha256(key, msg []byte) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write(msg)
	return mac.Sum(nil)
}