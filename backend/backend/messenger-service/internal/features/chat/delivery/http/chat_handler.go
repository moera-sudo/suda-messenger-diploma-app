package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/chat/service"
	httpUtils "github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/http"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/utils"
)

type ChatHandler struct {
	svc *service.ChatService
}

func NewChatHandler(svc *service.ChatService) *ChatHandler {
	return &ChatHandler{svc: svc}
}

func (h *ChatHandler) RegisterRoutes(e *echo.Echo) {
	chats := e.Group("/chats")

	// Chat CRUD
	chats.POST("", h.CreateChat)
	chats.GET("", h.GetUserChats)
	chats.GET("/:id/info", h.GetChatInfo)
	chats.PUT("/:id", h.UpdateChat)
	chats.DELETE("/:id", h.DeleteChat)

	// Media in chat
	chats.GET("/:id/media", h.GetChatMedia)
}

// CreateChat godoc
// @Summary Create a new chat
// @Tags Chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body chat.CreateChatReq true "Chat creation data"
// @Success 201 {object} chat.Chat
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /chats [post]
func (h *ChatHandler) CreateChat(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req chat.CreateChatReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}

	resp, err := h.svc.CreateChat(c.Request().Context(), userID, &req)
	if err != nil {
		return classifyAndRespond(c, err)
	}

	return c.JSON(http.StatusCreated, resp)
}

// GetUserChats godoc
// @Summary Get all chats for current user
// @Tags Chats
// @Produce json
// @Security BearerAuth
// @Success 200 {array} chat.ChatListItemResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /chats [get]
func (h *ChatHandler) GetUserChats(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	resp, err := h.svc.GetUserChats(c.Request().Context(), userID)
	if err != nil {
		return httpUtils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

// GetChatInfo godoc
// @Summary Get detailed chat info with members
// @Tags Chats
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Success 200 {object} chat.ChatInfoResponse
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Failure 404 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/info [get]
func (h *ChatHandler) GetChatInfo(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	resp, err := h.svc.GetChatInfo(c.Request().Context(), chatID, userID)
	if err != nil {
		return classifyAndRespond(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}

// UpdateChat godoc
// @Summary Update chat name, description or avatar
// @Tags Chats
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Param request body chat.UpdateChatReq true "Update data"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /chats/{id} [put]
func (h *ChatHandler) UpdateChat(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	var req chat.UpdateChatReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}

	if err := h.svc.UpdateChat(c.Request().Context(), chatID, userID, &req); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "CHAT_UPDATED", "Chat updated successfully")
}

// DeleteChat godoc
// @Summary Delete a chat
// @Description scope=everyone (default) deletes the chat for all members (DIRECT: any member, GROUP/CHANNEL: owner). scope=me clears history and hides the chat only for the caller.
// @Tags Chats
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Param scope query string false "me | everyone" Enums(me, everyone)
// @Success 204 "No Content"
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /chats/{id} [delete]
func (h *ChatHandler) DeleteChat(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	// scope=me -> per-user clear+hide; anything else (default) -> delete for everyone.
	if c.QueryParam("scope") == "me" {
		if err := h.svc.DeleteChatForMe(c.Request().Context(), chatID, userID); err != nil {
			return classifyAndRespond(c, err)
		}
		return c.NoContent(http.StatusNoContent)
	}

	if err := h.svc.DeleteChat(c.Request().Context(), chatID, userID); err != nil {
		return classifyAndRespond(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetChatMedia godoc
// @Summary Get media files shared in chat
// @Tags Chats
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Success 200 {object} chat.ChatMediaResponse
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/media [get]
func (h *ChatHandler) GetChatMedia(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	resp, err := h.svc.GetChatMedia(c.Request().Context(), chatID, userID)
	if err != nil {
		return classifyAndRespond(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}