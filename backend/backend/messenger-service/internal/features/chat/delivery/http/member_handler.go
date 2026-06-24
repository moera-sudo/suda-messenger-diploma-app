package http

import (
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/service"
	httpUtils "github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/http"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/utils"
)

type MemberHandler struct {
	svc *service.ChatService
}

func NewMemberHandler(svc *service.ChatService) *MemberHandler {
	return &MemberHandler{svc: svc}
}

func (h *MemberHandler) RegisterRoutes(e *echo.Echo) {
	chats := e.Group("/chats")

	// Members
	chats.POST("/:id/members", h.AddMember)
	chats.DELETE("/:id/members/:uid", h.RemoveMember)
	chats.PUT("/:id/members/:uid/role", h.UpdateMemberRole)
	chats.POST("/:id/leave", h.LeaveChat)

	// Mute
	chats.POST("/:id/mute", h.MuteChat)
	chats.POST("/:id/unmute", h.UnmuteChat)

	// Pin
	chats.POST("/:id/pin", h.PinMessage)
	chats.DELETE("/:id/pin/:mid", h.UnpinMessage)

	// Block (user-level, не chat-level)
	users := e.Group("/users")
	users.POST("/block", h.BlockUser)
	users.POST("/unblock", h.UnblockUser)
	users.GET("/blocked", h.GetBlockedUsers)

	// Contacts
	users.POST("/contacts", h.SetContactName)
	users.DELETE("/contacts/:id", h.RemoveContactName)
	users.GET("/contacts", h.GetContacts)
}

// ── Members ─────────────────────────────────────────────────

// AddMember godoc
// @Summary Add a member to a group chat
// @Tags Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Param request body chat.AddMemberReq true "User to add"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/members [post]
func (h *MemberHandler) AddMember(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	var req chat.AddMemberReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}

	if err := h.svc.AddMember(c.Request().Context(), chatID, userID, req.UserID); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "MEMBER_ADDED", "Member added successfully")
}

// RemoveMember godoc
// @Summary Remove a member from a group chat
// @Tags Members
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Param uid path string true "User ID to remove (UUID)"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/members/{uid} [delete]
func (h *MemberHandler) RemoveMember(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	targetID, err := uuid.Parse(c.Param("uid"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid user ID")
	}

	if err := h.svc.RemoveMember(c.Request().Context(), chatID, userID, targetID); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "MEMBER_REMOVED", "Member removed successfully")
}

// UpdateMemberRole godoc
// @Summary Update a member's role in chat
// @Tags Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Param uid path string true "Target user ID (UUID)"
// @Param request body chat.UpdateMemberRoleReq true "New role"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/members/{uid}/role [put]
func (h *MemberHandler) UpdateMemberRole(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	targetID, err := uuid.Parse(c.Param("uid"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid user ID")
	}

	var req chat.UpdateMemberRoleReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}

	if err := h.svc.UpdateMemberRole(c.Request().Context(), chatID, userID, targetID, req.Role); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "ROLE_UPDATED", "Member role updated")
}

// LeaveChat godoc
// @Summary Leave a group chat
// @Tags Members
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/leave [post]
func (h *MemberHandler) LeaveChat(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	if err := h.svc.LeaveChat(c.Request().Context(), chatID, userID); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "LEFT_CHAT", "Left chat successfully")
}

// ── Mute ────────────────────────────────────────────────────

// MuteChat godoc
// @Summary Mute chat notifications
// @Tags Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Param request body chat.MuteChatReq true "Mute settings"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/mute [post]
func (h *MemberHandler) MuteChat(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	var req chat.MuteChatReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}

	if err := h.svc.MuteChat(c.Request().Context(), userID, chatID, req.MutedUntil); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "CHAT_MUTED", "Chat muted")
}

// UnmuteChat godoc
// @Summary Unmute chat notifications
// @Tags Members
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/unmute [post]
func (h *MemberHandler) UnmuteChat(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	if err := h.svc.UnmuteChat(c.Request().Context(), userID, chatID); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "CHAT_UNMUTED", "Chat unmuted")
}

// ── Pin ─────────────────────────────────────────────────────

// PinMessage godoc
// @Summary Pin a message in chat
// @Tags Members
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Param request body chat.PinMessageReq true "Message to pin"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/pin [post]
func (h *MemberHandler) PinMessage(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	var req chat.PinMessageReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}

	if err := h.svc.PinMessage(c.Request().Context(), chatID, userID, req.MessageID); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "MESSAGE_PINNED", "Message pinned")
}

// UnpinMessage godoc
// @Summary Unpin a message in chat
// @Tags Members
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Param mid path int true "Message ID"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/pin/{mid} [delete]
func (h *MemberHandler) UnpinMessage(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	mid, err := parseInt64(c.Param("mid"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid message ID")
	}

	if err := h.svc.UnpinMessage(c.Request().Context(), chatID, userID, mid); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "MESSAGE_UNPINNED", "Message unpinned")
}

// ── Block ───────────────────────────────────────────────────

// BlockUser godoc
// @Summary Block a user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body chat.BlockUserReq true "User to block"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /users/block [post]
func (h *MemberHandler) BlockUser(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req chat.BlockUserReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}

	if userID == req.UserID {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Cannot block yourself")
	}

	if err := h.svc.BlockUser(c.Request().Context(), userID, req.UserID); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "USER_BLOCKED", "User blocked")
}

// UnblockUser godoc
// @Summary Unblock a user
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body chat.BlockUserReq true "User to unblock"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /users/unblock [post]
func (h *MemberHandler) UnblockUser(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req chat.BlockUserReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}

	if err := h.svc.UnblockUser(c.Request().Context(), userID, req.UserID); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "USER_UNBLOCKED", "User unblocked")
}

// GetBlockedUsers godoc
// @Summary Get list of blocked users (with public profile info)
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {array} chat.BlockedUserInfo
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /users/blocked [get]
func (h *MemberHandler) GetBlockedUsers(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	users, err := h.svc.GetBlockedUsersDetailed(c.Request().Context(), userID)
	if err != nil {
		return httpUtils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return c.JSON(http.StatusOK, users)
}

// ── Contacts ────────────────────────────────────────────────

// SetContactName godoc
// @Summary Set a custom name for a contact
// @Tags Contacts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body chat.SetContactNameReq true "Contact data"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /users/contacts [post]
func (h *MemberHandler) SetContactName(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req chat.SetContactNameReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}

	if err := h.svc.SetContactName(c.Request().Context(), userID, req.UserID, req.CustomName); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "CONTACT_SET", "Contact name set")
}

// RemoveContactName godoc
// @Summary Remove a custom contact name
// @Tags Contacts
// @Produce json
// @Security BearerAuth
// @Param id path string true "Contact user ID (UUID)"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /users/contacts/{id} [delete]
func (h *MemberHandler) RemoveContactName(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	contactID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid user ID")
	}

	if err := h.svc.RemoveContactName(c.Request().Context(), userID, contactID); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "CONTACT_REMOVED", "Contact name removed")
}

// GetContacts godoc
// @Summary Get all contacts with custom names
// @Tags Contacts
// @Produce json
// @Security BearerAuth
// @Success 200 {array} chat.Contact
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /users/contacts [get]
func (h *MemberHandler) GetContacts(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	contacts, err := h.svc.GetContacts(c.Request().Context(), userID)
	if err != nil {
		return httpUtils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return c.JSON(http.StatusOK, contacts)
}

// ── Helpers ─────────────────────────────────────────────────

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}
