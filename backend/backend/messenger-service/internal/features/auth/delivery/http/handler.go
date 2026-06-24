package http


import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/auth"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/auth/service"
	httpUtils "github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/http"
	"github.com/rs/zerolog/log"
)

type AuthHandler struct {
	service *service.Service
}

func NewAuthHandler(s *service.Service) *AuthHandler{
	return &AuthHandler{service: s}
}

func (h *AuthHandler) RegisterAuthRoutes(e *echo.Echo) {
	g := e.Group("/auth")
	g.POST("/register", h.Register)
	g.POST("/verify", h.Verify)
	g.POST("/login", h.Login)
	g.POST("/refresh", h.Refresh)
	g.POST("/forgot-password", h.Forgot)
	g.POST("/reset-password", h.Reset)
	g.POST("/webapp-session", h.WebappSession)
}

// Register godoc
// @Summary Register a new user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.RegisterReq true "Registration data"
// @Success 201 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 409 {object} httpUtils.ErrorResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(c echo.Context) error {
	var req auth.RegisterReq

	if err := c.Bind(&req); err != nil {
		log.Warn().Msg("Failed to bind request")
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid Input")
	}
	if err := c.Validate(&req); err != nil {
		return err
	} 

	err := h.service.Register(c.Request().Context(), &req)
	if err != nil {
		if err.Error() == "User already exists" {
			return httpUtils.RespondError(c, 409, "USER_EXISTS", err.Error())
		}
		if err.Error() == "Username is already taken" {
			return httpUtils.RespondError(c, 409, "USERNAME_TAKEN", err.Error())
		}
		if err.Error() == "Email is already taken" {
			return httpUtils.RespondError(c, 409, "EMAIL_TAKEN", err.Error())
		}
		if err.Error() == "Email and username are already taken" {
			return httpUtils.RespondError(c, 409, "USER_CREDENTIALS_TAKEN", err.Error())
		}
		return httpUtils.RespondError(c, 500, "INTERNAL ERROR", err.Error())
 	}

	return httpUtils.RespondByMessage(c, 201, "USER_CREATED", "User was successfully created. Verification code sent to email")
}


// Verify godoc
// @Summary Verify email with code
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.VerifyReq true "Verification code"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /auth/verify [post]
func (h *AuthHandler) Verify(c echo.Context) error {
	var req auth.VerifyReq
	if err := c.Bind(&req); err != nil{
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	err := h.service.Verify(c.Request().Context(), &req)
	if err != nil {
		return httpUtils.RespondError(c, 400, "INVALID_CODE", "Invalid or expired verification code")
	}

	return httpUtils.RespondByMessage(c, 200, "EMAIL_VERIFIED", "Email was verified. Now, please login")
}

// Login godoc
// @Summary Log in with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.LoginReq true "Login credentials"
// @Success 200 {object} auth.TokenResponse
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 401 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c echo.Context) error {
	var req auth.LoginReq
	req.ClientIP = c.RealIP()
	req.UserAgent = c.Request().UserAgent()

	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	resp, err := h.service.Login(c.Request().Context(), &req)
	if err != nil {
		  if err.Error() == "Email not verified" {
			httpUtils.RespondError(c, 403, "NOT_VERIFIED", "Email not verified")
		}	
		log.Warn().Err(err).Msg("Login failed")
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid credentials")
	}

	return c.JSON(200, resp)

}

// Refresh godoc
// @Summary Refresh access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.RefreshSessionReq true "Refresh token"
// @Success 200 {object} auth.TokenResponse
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 401 {object} httpUtils.ErrorResponse
// @Router /auth/refresh [post]
func (h *AuthHandler) Refresh(c echo.Context) error {
	var req auth.RefreshSessionReq
	req.ClientIP = c.RealIP()
	req.UserAgent = c.Request().UserAgent()

	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	resp, err := h.service.RefreshSessionTokens(c.Request().Context(), req.RefreshToken, req.UserAgent, req.ClientIP)
	if err != nil {
		if err.Error() == "Invalid or revoked session" {
			return httpUtils.RespondError(c, 401, "SESSION_EXPIRED", "Current session expired")
		}
		log.Warn().Err(err).Msg("Session refresh failed")
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Cannot refresh session")
	}

	return c.JSON(200, resp)

}

// Forgot godoc
// @Summary Request password reset code
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.ForgotPasswordReq true "Email address"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /auth/forgot-password [post]
func (h *AuthHandler) Forgot(c echo.Context) error {
	var req auth.ForgotPasswordReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}

	if err := c.Validate(&req); err != nil {
		return err
	}

	err := h.service.ForgotPassword(c.Request().Context(), &req)

	if err != nil{
		log.Warn().Err(err).Msg("Forgot password request failed")
	}

	return httpUtils.RespondByMessage(c, 200, "REQUEST_ON_PASS_RESET", "If email exists, verification code for password reset sent")
}

// Reset godoc
// @Summary Reset password with verification code
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.ResetPasswordReq true "Reset password data"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /auth/reset-password [post]
func (h *AuthHandler) Reset(c echo.Context) error {
	var req auth.ResetPasswordReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}
	if err := c.Validate(&req); err != nil {
		return err
	}

	err := h.service.ResetPassword(c.Request().Context(), &req)

	if err != nil{
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", err.Error())
	}

	return httpUtils.RespondByMessage(c, 200, "RESET_PASSWORD", "Password changed successfully")
}

// WebappSession godoc
// @Summary Exchange Mini App init_data for a short-lived access JWT
// @Description Called by WebView mini-apps (Wallet / Marketplace) on load.
// @Description The SPA reads `?initData=...` from its URL and posts it here;
// @Description backend validates HMAC + auth_date and returns an access JWT
// @Description (no refresh token — when it expires the SPA re-runs this
// @Description endpoint with the same initData, or the host app reopens the
// @Description WebView with a fresh one).
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body auth.WebappSessionReq true "Init data string from URL"
// @Success 200 {object} auth.WebappSessionResponse
// @Failure 400 {object} httpUtils.ErrorResponse "INVALID_INIT_DATA | INIT_DATA_EXPIRED | VALIDATION_ERROR"
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /auth/webapp-session [post]
func (h *AuthHandler) WebappSession(c echo.Context) error {
	var req auth.WebappSessionReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}
	if err := c.Validate(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "VALIDATION_ERROR", err.Error())
	}

	resp, err := h.service.ExchangeInitData(c.Request().Context(), req.InitData)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidInitData):
			return httpUtils.RespondError(c, http.StatusBadRequest, "INVALID_INIT_DATA", "init_data signature is invalid")
		case errors.Is(err, service.ErrInitDataExpired):
			return httpUtils.RespondError(c, http.StatusBadRequest, "INIT_DATA_EXPIRED", "init_data is older than the allowed 30 minutes")
		default:
			log.Error().Err(err).Msg("webapp-session: unexpected error")
			return httpUtils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "failed to issue webapp session")
		}
	}

	return c.JSON(http.StatusOK, resp)
}

