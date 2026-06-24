package http

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/user"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/user/service"
	httpUtils "github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/http"
	utils "github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/utils"
)

type UserHandler struct {
	service *service.Service
}

func NewUserHandler(s *service.Service) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/user")
	g.GET("/me", h.GetMe)
	g.GET("/id/:id", h.GetUserByID)
	g.POST("/logout", h.Logout)
	g.PUT("/profile", h.UpdateProfile)
	g.PUT("/avatar", h.UpdateAvatar) // Было POST с multipart, теперь PUT с media_id
	g.POST("/device", h.RegisterDevice)
	g.POST("/init-data", h.GetInitData)
	g.GET("/id/:id/status", h.GetUserStatus)

	// Account self-management
	g.PUT("/password", h.ChangePassword)
	g.GET("/sessions", h.ListSessions)
	g.DELETE("/sessions/:id", h.RevokeSession)
	g.POST("/sessions/terminate-others", h.TerminateOtherSessions)
}

// GetUserByID godoc
// @Summary Get user profile by ID
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} user.UserResponse
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Failure 404 {object} httpUtils.ErrorResponse
// @Router /user/id/{id} [get]
func (h *UserHandler) GetUserByID(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	targetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid UUID format")
	}

	// Свой профиль — без проверки блока
	if userID == targetID {
		resp, err := h.service.GetUser(c.Request().Context(), userID)
		if err != nil {
			return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
		}
		if resp == nil {
			return httpUtils.RespondError(c, 404, "NOT_FOUND", "User not found")
		}
		return c.JSON(200, resp)
	}

	resp, err := h.service.GetUserWithBlockCheck(c.Request().Context(), userID, targetID)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return httpUtils.RespondError(c, 403, "FORBIDDEN", "User is blocked")
		}
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	if resp == nil {
		return httpUtils.RespondError(c, 404, "NOT_FOUND", "User not found")
	}

	return c.JSON(200, resp)
}

// GetMe godoc
// @Summary Get current user profile
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} user.UserResponse
// @Failure 404 {object} httpUtils.ErrorResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /user/me [get]
func (h *UserHandler) GetMe(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	resp, err := h.service.GetUser(c.Request().Context(), userID)
	if err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	if resp == nil {
		return httpUtils.RespondError(c, 404, "NOT_FOUND", "User not found")
	}

	return c.JSON(200, resp)
}

// GetUserStatus godoc
// @Summary Get user online status
// @Tags User
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} user.UserStatusResponse
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /user/id/{id}/status [get]
func (h *UserHandler) GetUserStatus(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	targetID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid UUID format")
	}

	resp, err := h.service.GetUserStatus(c.Request().Context(), userID, targetID)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return httpUtils.RespondError(c, 403, "FORBIDDEN", "User is blocked")
		}
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}

	return c.JSON(200, resp)
}

// Logout godoc
// @Summary Log out and revoke session
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user.LogoutReq true "Logout data"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /user/logout [post]
func (h *UserHandler) Logout(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	var req user.LogoutReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}

	if err := h.service.Logout(c.Request().Context(), userID, &req); err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}

	return httpUtils.RespondByMessage(c, 200, "USER_LOGGED_OUT", "Logged out successfully")
}

// UpdateProfile godoc
// @Summary Update user profile
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user.UpdateProfileReq true "Profile data"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /user/profile [put]
func (h *UserHandler) UpdateProfile(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	var req user.UpdateProfileReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}

	if err := c.Validate(&req); err != nil {
		return httpUtils.RespondError(c, 400, "VALIDATION_ERROR", err.Error())
	}

	if err := h.service.UpdateProfile(c.Request().Context(), userID, req.DisplayName, req.Bio, req.FirstName, req.LastName); err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}

	return httpUtils.RespondByMessage(c, 200, "PROFILE_UPDATED", "Profile successfully updated")
}

// UpdateAvatar godoc
// @Summary Update user avatar by media ID
// @Description Client uploads avatar via media-service first, then passes media_id here
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user.UpdateAvatarReq true "Avatar media ID"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /user/avatar [put]
func (h *UserHandler) UpdateAvatar(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	var req user.UpdateAvatarReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}

	if err := h.service.UpdateAvatar(c.Request().Context(), userID, req.MediaID); err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}

	return httpUtils.RespondByMessage(c, 200, "AVATAR_UPDATED", "Avatar updated successfully")
}

// RegisterDevice godoc
// @Summary Register a device for push notifications
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user.RegisterDeviceReq true "Device data"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /user/device [post]
func (h *UserHandler) RegisterDevice(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	var req user.RegisterDeviceReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}

	if err := c.Validate(&req); err != nil {
		return httpUtils.RespondError(c, 400, "VALIDATION_ERROR", err.Error())
	}

	if err := h.service.RegisterDevice(c.Request().Context(), userID, req.FCMToken, req.DeviceName); err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", "Failed to save device token")
	}

	return httpUtils.RespondByMessage(c, 200, "NEW_DEVICE", "Device registered")
}

// GetInitData godoc
// @Summary Generate init data for mini-app integration
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string "initData string"
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /user/init-data [post]
func (h *UserHandler) GetInitData(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	initDataStr, err := h.service.GenerateInitData(c.Request().Context(), userID)
	if err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}

	return c.JSON(200, map[string]string{"initData": initDataStr})
}

// ChangePassword godoc
// @Summary Change password (authenticated)
// @Description Verifies old password, sets new one, revokes all refresh sessions (user re-logs in everywhere).
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user.ChangePasswordReq true "Old and new password"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse "INVALID_OLD_PASSWORD | VALIDATION_ERROR"
// @Failure 401 {object} httpUtils.ErrorResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /user/password [put]
func (h *UserHandler) ChangePassword(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}
	var req user.ChangePasswordReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return httpUtils.RespondError(c, 400, "VALIDATION_ERROR", err.Error())
	}
	if err := h.service.ChangePassword(c.Request().Context(), userID, req.OldPassword, req.NewPassword); err != nil {
		if err == service.ErrInvalidOldPassword {
			return httpUtils.RespondError(c, 400, "INVALID_OLD_PASSWORD", "Old password is incorrect")
		}
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	return httpUtils.RespondByMessage(c, 200, "PASSWORD_CHANGED", "Password changed; please re-login")
}

// ListSessions godoc
// @Summary List active refresh sessions
// @Tags User
// @Produce json
// @Security BearerAuth
// @Success 200 {array} user.SessionInfo
// @Failure 401 {object} httpUtils.ErrorResponse
// @Router /user/sessions [get]
func (h *UserHandler) ListSessions(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}
	// X-Session-ID is forwarded by the gateway from the access token "sid" claim.
	// Empty/invalid → uuid.Nil → no session is marked current.
	currentSessionID, _ := uuid.Parse(c.Request().Header.Get("X-Session-ID"))
	sessions, err := h.service.ListSessions(c.Request().Context(), userID, currentSessionID)
	if err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	return c.JSON(200, sessions)
}

// RevokeSession godoc
// @Summary Revoke a specific refresh session
// @Tags User
// @Security BearerAuth
// @Param id path string true "Session UUID"
// @Success 204 "No Content"
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 401 {object} httpUtils.ErrorResponse
// @Failure 404 {object} httpUtils.ErrorResponse
// @Router /user/sessions/{id} [delete]
func (h *UserHandler) RevokeSession(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}
	sid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid session id")
	}
	if err := h.service.RevokeSession(c.Request().Context(), userID, sid); err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	return c.NoContent(204)
}

// TerminateOtherSessions godoc
// @Summary Revoke all other refresh sessions (keep current)
// @Description Client supplies its own refresh token in body so the backend knows which session to keep alive.
// @Tags User
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body user.TerminateOtherSessionsReq true "Current refresh token to preserve"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 401 {object} httpUtils.ErrorResponse
// @Router /user/sessions/terminate-others [post]
func (h *UserHandler) TerminateOtherSessions(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}
	var req user.TerminateOtherSessionsReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid request body")
	}
	if err := c.Validate(&req); err != nil {
		return httpUtils.RespondError(c, 400, "VALIDATION_ERROR", err.Error())
	}
	n, err := h.service.RevokeOtherSessions(c.Request().Context(), userID, req.CurrentRefreshToken)
	if err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	return httpUtils.RespondByMessage(c, 200, "SESSIONS_REVOKED", fmt.Sprintf("%d sessions revoked", n))
}
