package http

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/channel/service"
	httpUtils "github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/http"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/utils"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/channels")

	g.POST("/:id/subscribe", h.Subscribe)
	g.POST("/:id/unsubscribe", h.Unsubscribe)
	g.GET("/:id/subscribers", h.GetSubscribers)
	g.GET("/by-username/:username", h.GetByUsername)
	g.GET("/:id/view", h.GetView)
	g.GET("/:id/settings", h.GetSettings)
	g.PUT("/:id/settings", h.UpdateSettings)

	// Private channel join requests (user side) + moderation (admin side)
	g.POST("/:id/join-request", h.RequestJoin)
	g.DELETE("/:id/join-request", h.CancelJoinRequest)
	g.GET("/:id/join-requests", h.ListJoinRequests)
	g.POST("/:id/join-requests/:uid/approve", h.ApproveJoin)
	g.POST("/:id/join-requests/:uid/reject", h.RejectJoin)

	// Invites: admin invites a user; user accepts/declines
	g.GET("/invites", h.ListMyInvites)
	g.POST("/:id/invites", h.InviteUser)
	g.DELETE("/:id/invites/:uid", h.RevokeInvite)
	g.POST("/:id/invites/accept", h.AcceptInvite)
	g.POST("/:id/invites/decline", h.DeclineInvite)

	g.GET("/:id/posts/:msg_id/comments", h.GetComments)
	g.POST("/:id/posts/:msg_id/comments", h.AddComment)
	g.PUT("/comments/:comment_id", h.EditComment)
	g.DELETE("/comments/:comment_id", h.DeleteComment)

	g.GET("/:id/apps", h.GetLinkedApps)
	g.POST("/:id/apps", h.LinkApp)
	g.DELETE("/:id/apps/:app_id", h.UnlinkApp)
}

// Subscribe godoc
// @Summary Subscribe to a channel
// @Description Adds the authenticated user as a subscriber. Token-gating is enforced — returns 402/403 if the user fails the rule.
// @Tags Channels
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel UUID"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 402 {object} httpUtils.ErrorResponse "INSUFFICIENT_BALANCE"
// @Failure 403 {object} httpUtils.ErrorResponse "NO_WALLET | NFT_REQUIRED | GATING_REQUIREMENT_NOT_MET | FORBIDDEN"
// @Router /channels/{id}/subscribe [post]
func (h *Handler) Subscribe(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}

	if err := h.svc.Subscribe(c.Request().Context(), channelID, userID); err != nil {
		return classify(c, err)
	}
	return httpUtils.RespondByMessage(c, 200, "SUBSCRIBED", "Subscribed to channel")
}

// Unsubscribe godoc
// @Summary Unsubscribe from a channel
// @Tags Channels
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel UUID"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /channels/{id}/unsubscribe [post]
func (h *Handler) Unsubscribe(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}

	if err := h.svc.Unsubscribe(c.Request().Context(), channelID, userID); err != nil {
		return classify(c, err)
	}
	return httpUtils.RespondByMessage(c, 200, "UNSUBSCRIBED", "Unsubscribed from channel")
}

// GetSubscribers godoc
// @Summary Get channel subscribers
// @Tags Channels
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel UUID"
// @Param limit query int false "Page size (1..100)" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} channel.SubscriberInfo
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /channels/{id}/subscribers [get]
func (h *Handler) GetSubscribers(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	subs, err := h.svc.GetSubscribers(c.Request().Context(), channelID, userID, limit, offset)
	if err != nil {
		return classify(c, err)
	}
	return c.JSON(http.StatusOK, subs)
}

// GetByUsername godoc
// @Summary Resolve channel by username
// @Description Returns channel_id by channel @username. The "@" prefix is trimmed automatically.
// @Tags Channels
// @Produce json
// @Security BearerAuth
// @Param username path string true "Channel username (with or without @)"
// @Success 200 {object} map[string]string "channel_id"
// @Failure 404 {object} httpUtils.ErrorResponse
// @Router /channels/by-username/{username} [get]
func (h *Handler) GetByUsername(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	username := strings.TrimPrefix(c.Param("username"), "@")
	view, err := h.svc.GetByUsername(c.Request().Context(), username, userID)
	if err != nil {
		return classify(c, err)
	}
	return c.JSON(http.StatusOK, view)
}

// GetView godoc
// @Summary Get channel card + viewer relationship
// @Description Returns channel metadata and the caller's join_status / can_read.
// @Tags Channels
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel UUID"
// @Success 200 {object} channel.ChannelView
// @Router /channels/{id}/view [get]
func (h *Handler) GetView(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	view, err := h.svc.GetView(c.Request().Context(), channelID, userID)
	if err != nil {
		return classify(c, err)
	}
	return c.JSON(http.StatusOK, view)
}

// ── Join requests & invites ────────────────────────────────

// RequestJoin — POST /channels/{id}/join-request (user asks to join a private channel).
func (h *Handler) RequestJoin(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	if err := h.svc.RequestJoin(c.Request().Context(), channelID, userID); err != nil {
		return classify(c, err)
	}
	return httpUtils.RespondByMessage(c, 200, "JOIN_REQUESTED", "Join request sent")
}

// CancelJoinRequest — DELETE /channels/{id}/join-request.
func (h *Handler) CancelJoinRequest(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	if err := h.svc.CancelJoinRequest(c.Request().Context(), channelID, userID); err != nil {
		return classify(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ListJoinRequests — GET /channels/{id}/join-requests (admin).
func (h *Handler) ListJoinRequests(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	list, err := h.svc.ListJoinRequests(c.Request().Context(), channelID, userID, limit, offset)
	if err != nil {
		return classify(c, err)
	}
	return c.JSON(http.StatusOK, list)
}

// ApproveJoin — POST /channels/{id}/join-requests/{uid}/approve (admin).
func (h *Handler) ApproveJoin(c echo.Context) error {
	return h.decideJoin(c, true)
}

// RejectJoin — POST /channels/{id}/join-requests/{uid}/reject (admin).
func (h *Handler) RejectJoin(c echo.Context) error {
	return h.decideJoin(c, false)
}

func (h *Handler) decideJoin(c echo.Context, approve bool) error {
	adminID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	targetID, err := uuid.Parse(c.Param("uid"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid user ID")
	}
	if approve {
		err = h.svc.ApproveJoin(c.Request().Context(), channelID, adminID, targetID)
	} else {
		err = h.svc.RejectJoin(c.Request().Context(), channelID, adminID, targetID)
	}
	if err != nil {
		return classify(c, err)
	}
	action := "JOIN_REJECTED"
	if approve {
		action = "JOIN_APPROVED"
	}
	return httpUtils.RespondByMessage(c, 200, action, "OK")
}

// InviteUser — POST /channels/{id}/invites (admin invites a user).
func (h *Handler) InviteUser(c echo.Context) error {
	adminID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	var req channel.InviteUserReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "invalid body")
	}

	// Accept either user_id or username (frontend sends username).
	var targetID uuid.UUID
	switch {
	case req.UserID != nil:
		targetID = *req.UserID
	case req.Username != "":
		targetID, err = h.svc.ResolveUserID(c.Request().Context(), req.Username)
		if err != nil {
			return classify(c, err)
		}
	default:
		return httpUtils.RespondError(c, 400, "VALIDATION_ERROR", "user_id or username required")
	}

	if err := h.svc.InviteUser(c.Request().Context(), channelID, adminID, targetID); err != nil {
		return classify(c, err)
	}
	return httpUtils.RespondByMessage(c, 200, "INVITED", "User invited")
}

// RevokeInvite — DELETE /channels/{id}/invites/{uid} (admin).
func (h *Handler) RevokeInvite(c echo.Context) error {
	adminID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	targetID, err := uuid.Parse(c.Param("uid"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid user ID")
	}
	if err := h.svc.RevokeInvite(c.Request().Context(), channelID, adminID, targetID); err != nil {
		return classify(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ListMyInvites — GET /channels/invites (current user's pending invites).
func (h *Handler) ListMyInvites(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	list, err := h.svc.ListMyInvites(c.Request().Context(), userID)
	if err != nil {
		return classify(c, err)
	}
	return c.JSON(http.StatusOK, list)
}

// AcceptInvite — POST /channels/{id}/invites/accept.
func (h *Handler) AcceptInvite(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	if err := h.svc.AcceptInvite(c.Request().Context(), channelID, userID); err != nil {
		return classify(c, err)
	}
	return httpUtils.RespondByMessage(c, 200, "INVITE_ACCEPTED", "Joined channel")
}

// DeclineInvite — POST /channels/{id}/invites/decline.
func (h *Handler) DeclineInvite(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	if err := h.svc.DeclineInvite(c.Request().Context(), channelID, userID); err != nil {
		return classify(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GetSettings godoc
// @Summary Get current channel settings (admin only)
// @Description Returns current editable settings to prefill the Settings form. OWNER/ADMIN only.
// @Tags Channels
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel UUID"
// @Success 200 {object} channel.ChannelSettingsResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /channels/{id}/settings [get]
func (h *Handler) GetSettings(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	settings, err := h.svc.GetSettings(c.Request().Context(), channelID, userID)
	if err != nil {
		return classify(c, err)
	}
	return c.JSON(http.StatusOK, settings)
}

// UpdateSettings godoc
// @Summary Update channel settings
// @Description Channel OWNER/ADMIN only.
// @Tags Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel UUID"
// @Param request body channel.UpdateChannelSettingsReq true "Settings to update"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /channels/{id}/settings [put]
func (h *Handler) UpdateSettings(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}

	var req channel.UpdateChannelSettingsReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}
	if err := c.Validate(&req); err != nil {
		return httpUtils.RespondError(c, 400, "VALIDATION_ERROR", err.Error())
	}

	if err := h.svc.UpdateSettings(c.Request().Context(), channelID, userID, &req); err != nil {
		return classify(c, err)
	}
	return httpUtils.RespondByMessage(c, 200, "SETTINGS_UPDATED", "Channel settings updated")
}

// ── Comments ───────────────────────────────────────────────

// GetComments godoc
// @Summary Get comments for a channel post
// @Tags Channels
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel UUID"
// @Param msg_id path int true "Post (message) ID"
// @Param limit query int false "Page size (1..100)" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} channel.CommentsListResponse
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /channels/{id}/posts/{msg_id}/comments [get]
func (h *Handler) GetComments(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	postID, err := strconv.ParseInt(c.Param("msg_id"), 10, 64)
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid post ID")
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	comments, total, err := h.svc.GetComments(c.Request().Context(), channelID, userID, postID, limit, offset)
	if err != nil {
		return classify(c, err)
	}

	return c.JSON(http.StatusOK, channel.CommentsListResponse{
		Comments: comments,
		Total:    total,
	})
}

// AddComment godoc
// @Summary Add a comment to a channel post
// @Tags Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel UUID"
// @Param msg_id path int true "Post (message) ID"
// @Param request body channel.AddCommentReq true "Comment content"
// @Success 201 {object} channel.ChannelComment
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /channels/{id}/posts/{msg_id}/comments [post]
func (h *Handler) AddComment(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	postID, err := strconv.ParseInt(c.Param("msg_id"), 10, 64)
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid post ID")
	}

	var req channel.AddCommentReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}
	if err := c.Validate(&req); err != nil {
		return httpUtils.RespondError(c, 400, "VALIDATION_ERROR", err.Error())
	}

	comment, err := h.svc.AddComment(c.Request().Context(), channelID, postID, userID, &req)
	if err != nil {
		return classify(c, err)
	}
	return c.JSON(http.StatusCreated, comment)
}

// EditComment godoc
// @Summary Edit a comment
// @Description Only the comment author can edit their comment.
// @Tags Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param comment_id path int true "Comment ID"
// @Param request body channel.EditCommentReq true "New content"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /channels/comments/{comment_id} [put]
func (h *Handler) EditComment(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	commentID, err := strconv.ParseInt(c.Param("comment_id"), 10, 64)
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid comment ID")
	}

	var req channel.EditCommentReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}
	if err := c.Validate(&req); err != nil {
		return httpUtils.RespondError(c, 400, "VALIDATION_ERROR", err.Error())
	}

	if err := h.svc.EditComment(c.Request().Context(), commentID, userID, req.Content); err != nil {
		return classify(c, err)
	}
	return httpUtils.RespondByMessage(c, 200, "COMMENT_EDITED", "Comment edited")
}

// DeleteComment godoc
// @Summary Delete a comment
// @Description Comment author or channel OWNER/ADMIN can delete.
// @Tags Channels
// @Security BearerAuth
// @Param comment_id path int true "Comment ID"
// @Success 204 "No Content"
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /channels/comments/{comment_id} [delete]
func (h *Handler) DeleteComment(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	commentID, err := strconv.ParseInt(c.Param("comment_id"), 10, 64)
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid comment ID")
	}

	if err := h.svc.DeleteComment(c.Request().Context(), commentID, userID); err != nil {
		return classify(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ── Apps ───────────────────────────────────────────────────

// GetLinkedApps godoc
// @Summary Get apps linked to a channel
// @Tags Channels
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel UUID"
// @Success 200 {array} channel.ChannelAppLink
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /channels/{id}/apps [get]
func (h *Handler) GetLinkedApps(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}

	apps, err := h.svc.GetLinkedApps(c.Request().Context(), channelID, userID)
	if err != nil {
		return classify(c, err)
	}
	return c.JSON(http.StatusOK, apps)
}

// LinkApp godoc
// @Summary Link an app to a channel
// @Description Channel OWNER/ADMIN only.
// @Tags Channels
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Channel UUID"
// @Param request body channel.LinkAppReq true "App ID to link"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /channels/{id}/apps [post]
func (h *Handler) LinkApp(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}

	var req channel.LinkAppReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}
	if err := c.Validate(&req); err != nil {
		return httpUtils.RespondError(c, 400, "VALIDATION_ERROR", err.Error())
	}

	if err := h.svc.LinkApp(c.Request().Context(), channelID, userID, req.AppID); err != nil {
		return classify(c, err)
	}
	return httpUtils.RespondByMessage(c, 200, "APP_LINKED", "App linked to channel")
}

// UnlinkApp godoc
// @Summary Unlink an app from a channel
// @Description Channel OWNER/ADMIN only.
// @Tags Channels
// @Security BearerAuth
// @Param id path string true "Channel UUID"
// @Param app_id path string true "App UUID"
// @Success 204 "No Content"
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /channels/{id}/apps/{app_id} [delete]
func (h *Handler) UnlinkApp(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}
	channelID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid channel ID")
	}
	appID, err := uuid.Parse(c.Param("app_id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid app ID")
	}

	if err := h.svc.UnlinkApp(c.Request().Context(), channelID, userID, appID); err != nil {
		return classify(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// ── Helpers ────────────────────────────────────────────────

func classify(c echo.Context, err error) error {
	// Token-gating: типизированная ошибка несёт reason + числа баланса.
	var ge *channel.GatingError
	if errors.As(err, &ge) {
		switch ge.Reason {
		case "insufficient_balance":
			return httpUtils.RespondError(c, 402, "INSUFFICIENT_BALANCE",
				fmt.Sprintf("channel requires %s SUDA wei, you have %s", ge.MinBalanceWei, ge.UserBalanceWei))
		case "wallet_not_found":
			return httpUtils.RespondError(c, 403, "NO_WALLET", "you need a wallet to join this channel")
		case "nft_required":
			return httpUtils.RespondError(c, 403, "NFT_REQUIRED", "this channel requires an NFT from a specific collection")
		default:
			return httpUtils.RespondError(c, 403, "GATING_REQUIREMENT_NOT_MET", ge.Error())
		}
	}

	msg := err.Error()
	switch {
	case strings.Contains(msg, "forbidden"):
		return httpUtils.RespondError(c, 403, "FORBIDDEN", msg)
	case strings.Contains(msg, "not found"):
		return httpUtils.RespondError(c, 404, "NOT_FOUND", msg)
	case strings.Contains(msg, "already taken"):
		return httpUtils.RespondError(c, 409, "CONFLICT", msg)
	default:
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", msg)
	}
}