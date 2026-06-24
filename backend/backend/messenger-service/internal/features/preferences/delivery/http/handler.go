package http

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/preferences/service"
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
	g := e.Group("/user")
	g.GET("/preferences", h.GetPreferences)
	g.PUT("/preferences", h.UpdatePreferences)

	g.POST("/last-seen-hidden", h.HideLastSeen)
	g.DELETE("/last-seen-hidden/:id", h.UnhideLastSeen)
	g.GET("/last-seen-hidden", h.GetHiddenList)
}

// GetPreferences godoc
// @Summary Get user preferences
// @Description Returns the authenticated user's preferences (notifications, privacy, language, etc.).
// @Tags Preferences
// @Produce json
// @Security BearerAuth
// @Success 200 {object} preferences.UserPreferences
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /user/preferences [get]
func (h *Handler) GetPreferences(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	p, err := h.svc.GetPreferences(c.Request().Context(), userID)
	if err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	return c.JSON(http.StatusOK, p)
}

// UpdatePreferences godoc
// @Summary Update user preferences
// @Tags Preferences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body preferences.UpdatePreferencesReq true "Fields to update"
// @Success 200 {object} preferences.UserPreferences
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /user/preferences [put]
func (h *Handler) UpdatePreferences(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req preferences.UpdatePreferencesReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}
	if err := c.Validate(&req); err != nil {
		return httpUtils.RespondError(c, 400, "VALIDATION_ERROR", err.Error())
	}

	updated, err := h.svc.UpdatePreferences(c.Request().Context(), userID, &req)
	if err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	return c.JSON(http.StatusOK, updated)
}

// HideLastSeen godoc
// @Summary Hide my last-seen from a specific user
// @Tags Preferences
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body preferences.HideLastSeenReq true "User ID to hide last-seen from"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /user/last-seen-hidden [post]
func (h *Handler) HideLastSeen(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req preferences.HideLastSeenReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}

	if err := h.svc.HideLastSeen(c.Request().Context(), userID, req.UserID); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", err.Error())
	}
	return httpUtils.RespondByMessage(c, 200, "LAST_SEEN_HIDDEN", "User added to hidden list")
}

// UnhideLastSeen godoc
// @Summary Unhide my last-seen for a specific user
// @Tags Preferences
// @Produce json
// @Security BearerAuth
// @Param id path string true "Hidden user UUID"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /user/last-seen-hidden/{id} [delete]
func (h *Handler) UnhideLastSeen(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	hiddenID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid user ID")
	}

	if err := h.svc.UnhideLastSeen(c.Request().Context(), userID, hiddenID); err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	return httpUtils.RespondByMessage(c, 200, "LAST_SEEN_UNHIDDEN", "User removed from hidden list")
}

// GetHiddenList godoc
// @Summary Get list of users from whom my last-seen is hidden
// @Tags Preferences
// @Produce json
// @Security BearerAuth
// @Success 200 {array} string "User UUIDs"
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /user/last-seen-hidden [get]
func (h *Handler) GetHiddenList(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	ids, err := h.svc.GetHiddenList(c.Request().Context(), userID)
	if err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	return c.JSON(http.StatusOK, ids)
}