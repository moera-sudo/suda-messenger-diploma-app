package http

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	_ "github.com/moera-sudo/backend/backend/messenger-service/internal/features/search" // swagger models
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/search/service"
	httpUtils "github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/http"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/utils"
)

type SearchHandler struct {
	service *service.Service
}

func NewSearchHandler(s *service.Service) *SearchHandler {
	return &SearchHandler{service: s}
}

func (h *SearchHandler) RegisterRoutes(e *echo.Echo) {
	e.GET("/search", h.GlobalSearch)
	e.GET("/chats/:id/search", h.ChatSearch)
}

// GlobalSearch godoc
// @Summary Search users, chats and messages globally
// @Tags Search
// @Produce json
// @Security BearerAuth
// @Param q query string true "Search query (min 2 chars)"
// @Success 200 {object} search.GlobalSearchResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /search [get]
func (h *SearchHandler) GlobalSearch(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	q := c.QueryParam("q")
	if len(q) < 2 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"users": []interface{}{}, "chats": []interface{}{}, "messages": []interface{}{},
		})
	}

	resp, err := h.service.GlobalSearch(c.Request().Context(), userID, q)
	if err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

// ChatSearch godoc
// @Summary Search messages within a specific chat
// @Tags Search
// @Produce json
// @Security BearerAuth
// @Param id path string true "Chat ID (UUID)"
// @Param q query string true "Search query (min 2 chars)"
// @Success 200 {object} search.ChatSearchResponse
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /chats/{id}/search [get]
func (h *SearchHandler) ChatSearch(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	chatID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid chat ID")
	}

	q := c.QueryParam("q")
	if len(q) < 2 {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"messages": []interface{}{}, "total": 0,
		})
	}

	resp, err := h.service.SearchInChat(c.Request().Context(), userID, chatID, q)
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return httpUtils.RespondError(c, 403, "FORBIDDEN", err.Error())
		}
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}
