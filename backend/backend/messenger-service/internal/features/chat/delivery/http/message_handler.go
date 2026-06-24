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

type MessageHandler struct {
	svc *service.ChatService
}

func NewMessageHandler(svc *service.ChatService) *MessageHandler {
	return &MessageHandler{svc: svc}
}

func (h *MessageHandler) RegisterRoutes(e *echo.Echo) {
	chats := e.Group("/chats")

	// Messages
	chats.POST("/:id/messages", h.SendMessage)
	chats.GET("/:id/messages", h.GetMessages)
	chats.POST("/forward", h.ForwardMessage)

	// Single message operations
	msgs := e.Group("/messages")
	msgs.PUT("/:id", h.EditMessage)
	msgs.DELETE("/:id", h.DeleteMessage)
	msgs.GET("/:id/readers", h.GetMessageReaders)
}

// SendMessage godoc
// @Summary Send a message to a chat
// @Tags Messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Param request body chat.SendMessageReq true "Message data"
// @Success 201 {object} chat.Message
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/messages [post]
func (h *MessageHandler) SendMessage(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	var req chat.SendMessageReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}
	req.ChatID = chatID

	resp, err := h.svc.SendMessage(c.Request().Context(), userID, &req)
	if err != nil {
		return classifyAndRespond(c, err)
	}

	return c.JSON(http.StatusCreated, resp)
}

// GetMessages godoc
// @Summary Get messages from a chat
// @Tags Messages
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Param limit query int false "Limit (default 50, max 100)"
// @Param offset query int false "Offset (default 0)"
// @Success 200 {array} chat.Message
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/messages [get]
func (h *MessageHandler) GetMessages(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid chat ID")
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset, _ := strconv.Atoi(c.QueryParam("offset"))

	resp, err := h.svc.GetChatMessages(c.Request().Context(), userID, chatID, limit, offset)
	if err != nil {
		return classifyAndRespond(c, err)
	}

	return c.JSON(http.StatusOK, resp)
}

// ForwardMessage godoc
// @Summary Forward a message to another chat
// @Tags Messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body chat.ForwardMessageReq true "Forward data"
// @Success 201 {object} chat.Message
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /chats/forward [post]
func (h *MessageHandler) ForwardMessage(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req chat.ForwardMessageReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}

	resp, err := h.svc.ForwardMessage(c.Request().Context(), userID, &req)
	if err != nil {
		return classifyAndRespond(c, err)
	}

	return c.JSON(http.StatusCreated, resp)
}

// EditMessage godoc
// @Summary Edit a message
// @Tags Messages
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Message ID"
// @Param request body chat.EditMessageReq true "New content"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /messages/{id} [put]
func (h *MessageHandler) EditMessage(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	messageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid message ID")
	}

	var req chat.EditMessageReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid input")
	}

	if err := h.svc.EditMessage(c.Request().Context(), userID, messageID, &req); err != nil {
		return classifyAndRespond(c, err)
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "MESSAGE_EDITED", "Message edited successfully")
}

// DeleteMessage godoc
// @Summary Delete a message
// @Tags Messages
// @Security BearerAuth
// @Param id path int true "Message ID"
// @Param for_everyone query bool false "Delete for all members"
// @Success 204 "No Content"
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	messageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid message ID")
	}

	forEveryone := c.QueryParam("for_everyone") == "true"

	req := &chat.DeleteMessageReq{ForEveryone: forEveryone}

	if err := h.svc.DeleteMessage(c.Request().Context(), userID, messageID, req); err != nil {
		return classifyAndRespond(c, err)
	}

	return c.NoContent(http.StatusNoContent)
}

// GetMessageReaders godoc
// @Summary Get users who read a message
// @Tags Messages
// @Produce json
// @Security BearerAuth
// @Param id path int true "Message ID"
// @Success 200 {array} chat.MessageReaderInfoResponse
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /messages/{id}/readers [get]
func (h *MessageHandler) GetMessageReaders(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	messageID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid message ID")
	}

	readers, err := h.svc.GetMessageReaders(c.Request().Context(), userID, messageID)
	if err != nil {
		return classifyAndRespond(c, err)
	}

	return c.JSON(http.StatusOK, readers)
}