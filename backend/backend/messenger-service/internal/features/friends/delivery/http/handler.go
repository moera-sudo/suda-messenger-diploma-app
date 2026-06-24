// Package http — Friends feature HTTP routes (mounted at /users/friends/*).
package http

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/friends"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/friends/service"
	httpUtils "github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/http"
	utils "github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/utils"
)

type Handler struct {
	svc *service.Service
}

func NewHandler(svc *service.Service) *Handler { return &Handler{svc: svc} }

func (h *Handler) RegisterRoutes(e *echo.Echo) {
	g := e.Group("/users/friends")
	g.GET("", h.ListFriends)
	g.DELETE("/:user_id", h.Unfriend)
	g.GET("/:user_id/status", h.RelationStatus)
	g.POST("/requests", h.SendRequest)
	g.GET("/requests", h.ListRequests)
	g.POST("/requests/:id/accept", h.Accept)
	g.POST("/requests/:id/reject", h.Reject)
	g.DELETE("/requests/:id", h.Cancel)
}

// SendRequest godoc
// @Summary Send friend request
// @Description Sends a PENDING request to target user. Errors if already friends / pending / self / blocked.
// @Tags Friends
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body friends.SendRequestReq true "Target user id"
// @Success 201 {object} friends.FriendRequest
// @Failure 400 {object} httpUtils.ErrorResponse "VALIDATION_ERROR | SELF_REQUEST"
// @Failure 403 {object} httpUtils.ErrorResponse "BLOCKED"
// @Failure 409 {object} httpUtils.ErrorResponse "ALREADY_FRIENDS | ALREADY_PENDING"
// @Router /users/friends/requests [post]
func (h *Handler) SendRequest(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid user")
	}
	var req friends.SendRequestReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "invalid body")
	}
	if err := c.Validate(&req); err != nil {
		return httpUtils.RespondError(c, 400, "VALIDATION_ERROR", err.Error())
	}
	fr, err := h.svc.SendRequest(c.Request().Context(), userID, req.UserID)
	if err != nil {
		return mapErr(c, err)
	}
	return c.JSON(http.StatusCreated, fr)
}

// ListRequests godoc
// @Summary List pending friend requests
// @Tags Friends
// @Produce json
// @Security BearerAuth
// @Param type query string false "incoming | outgoing (default: both)"
// @Param limit query int false "Page size (1..100)" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} friends.FriendRequestInfo
// @Router /users/friends/requests [get]
func (h *Handler) ListRequests(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid user")
	}
	dir := c.QueryParam("type")
	limit, offset := pageParams(c, 50)
	out, err := h.svc.ListRequests(c.Request().Context(), userID, dir, limit, offset)
	if err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	return c.JSON(http.StatusOK, out)
}

// Accept godoc
// @Summary Accept friend request
// @Description Target user accepts. After this both users are FRIENDS.
// @Tags Friends
// @Security BearerAuth
// @Param id path string true "Request UUID"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 403 {object} httpUtils.ErrorResponse "FORBIDDEN"
// @Failure 404 {object} httpUtils.ErrorResponse
// @Failure 409 {object} httpUtils.ErrorResponse "INVALID_TRANSITION"
// @Router /users/friends/requests/{id}/accept [post]
func (h *Handler) Accept(c echo.Context) error { return h.transition(c, h.svc.AcceptRequest, "FRIENDS_NOW") }

// Reject godoc
// @Summary Reject friend request
// @Tags Friends
// @Security BearerAuth
// @Param id path string true "Request UUID"
// @Success 200 {object} httpUtils.MessageRespond
// @Router /users/friends/requests/{id}/reject [post]
func (h *Handler) Reject(c echo.Context) error { return h.transition(c, h.svc.RejectRequest, "REJECTED") }

// Cancel godoc
// @Summary Cancel sent friend request
// @Description Requester withdraws their own PENDING request.
// @Tags Friends
// @Security BearerAuth
// @Param id path string true "Request UUID"
// @Success 200 {object} httpUtils.MessageRespond
// @Router /users/friends/requests/{id} [delete]
func (h *Handler) Cancel(c echo.Context) error { return h.transition(c, h.svc.CancelRequest, "CANCELLED") }

// transition is a shared helper for Accept / Reject / Cancel — they all take (userID, requestID)
// and the only difference is the service method + the action label returned to the client.
func (h *Handler) transition(c echo.Context, fn func(ctx context.Context, userID, requestID uuid.UUID) error, action string) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid user")
	}
	rid, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "invalid request id")
	}
	if err := fn(c.Request().Context(), userID, rid); err != nil {
		return mapErr(c, err)
	}
	return httpUtils.RespondByMessage(c, http.StatusOK, action, "OK")
}

// ListFriends godoc
// @Summary List my friends
// @Tags Friends
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Page size (1..100)" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} friends.FriendInfo
// @Router /users/friends [get]
func (h *Handler) ListFriends(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid user")
	}
	limit, offset := pageParams(c, 50)
	out, err := h.svc.ListFriends(c.Request().Context(), userID, limit, offset)
	if err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	return c.JSON(http.StatusOK, out)
}

// Unfriend godoc
// @Summary Remove friend
// @Tags Friends
// @Security BearerAuth
// @Param user_id path string true "Friend user UUID"
// @Success 204 "No Content"
// @Router /users/friends/{user_id} [delete]
func (h *Handler) Unfriend(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid user")
	}
	otherID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "invalid user_id")
	}
	if err := h.svc.Unfriend(c.Request().Context(), userID, otherID); err != nil {
		return mapErr(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// RelationStatus godoc
// @Summary Get relation status with a user
// @Description Returns NONE / PENDING_SENT / PENDING_RECEIVED / FRIENDS / BLOCKED.
// @Tags Friends
// @Produce json
// @Security BearerAuth
// @Param user_id path string true "Other user UUID"
// @Success 200 {object} friends.RelationStatusResponse
// @Router /users/friends/{user_id}/status [get]
func (h *Handler) RelationStatus(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid user")
	}
	otherID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "invalid user_id")
	}
	rel, rid, err := h.svc.GetRelation(c.Request().Context(), userID, otherID)
	if err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	return c.JSON(http.StatusOK, friends.RelationStatusResponse{UserID: otherID, Relation: rel, RequestID: rid})
}

func pageParams(c echo.Context, defLimit int) (int, int) {
	limit := defLimit
	if v, err := strconv.Atoi(c.QueryParam("limit")); err == nil && v > 0 && v <= 100 {
		limit = v
	}
	offset := 0
	if v, err := strconv.Atoi(c.QueryParam("offset")); err == nil && v >= 0 {
		offset = v
	}
	return limit, offset
}

func mapErr(c echo.Context, err error) error {
	switch {
	case errors.Is(err, friends.ErrSelfRequest):
		return httpUtils.RespondError(c, 400, "SELF_REQUEST", err.Error())
	case errors.Is(err, friends.ErrBlocked):
		return httpUtils.RespondError(c, 403, "BLOCKED", err.Error())
	case errors.Is(err, friends.ErrAlreadyFriends):
		return httpUtils.RespondError(c, 409, "ALREADY_FRIENDS", err.Error())
	case errors.Is(err, friends.ErrAlreadyPending):
		return httpUtils.RespondError(c, 409, "ALREADY_PENDING", err.Error())
	case errors.Is(err, friends.ErrRequestNotFound):
		return httpUtils.RespondError(c, 404, "NOT_FOUND", err.Error())
	case errors.Is(err, friends.ErrForbidden):
		return httpUtils.RespondError(c, 403, "FORBIDDEN", err.Error())
	case errors.Is(err, friends.ErrInvalidTransition):
		return httpUtils.RespondError(c, 409, "INVALID_TRANSITION", err.Error())
	default:
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
}
