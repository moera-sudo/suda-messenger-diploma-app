package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"

	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/userpins"
	"github.com/moera-sudo/backend/backend/messenger-service/internal/features/userpins/service"
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
	g := e.Group("/user/pins")
	g.POST("", h.CreatePin)
	g.DELETE("/:id", h.DeletePin)
	g.PUT("/reorder", h.Reorder)
	g.GET("", h.GetPins)
}

// CreatePin godoc
// @Summary Create a user pin
// @Description Pins an entity (chat, channel) to the user's pin list (chatlist or other).
// @Tags UserPins
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userpins.CreatePinReq true "Pin parameters (type, target_id)"
// @Success 201 {object} userpins.UserPin
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Router /user/pins [post]
func (h *Handler) CreatePin(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req userpins.CreatePinReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}
	if err := c.Validate(&req); err != nil {
		return httpUtils.RespondError(c, 400, "VALIDATION_ERROR", err.Error())
	}

	pin, err := h.svc.CreatePin(c.Request().Context(), userID, &req)
	if err != nil {
		return classify(c, err)
	}
	return c.JSON(http.StatusCreated, pin)
}

// DeletePin godoc
// @Summary Delete a user pin
// @Tags UserPins
// @Security BearerAuth
// @Param id path int true "Pin ID"
// @Success 204 "No Content"
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 404 {object} httpUtils.ErrorResponse
// @Router /user/pins/{id} [delete]
func (h *Handler) DeletePin(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	pinID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid pin ID")
	}

	if err := h.svc.DeletePin(c.Request().Context(), userID, pinID); err != nil {
		return classify(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// Reorder godoc
// @Summary Reorder user pins
// @Tags UserPins
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body userpins.ReorderPinsReq true "New pin order"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /user/pins/reorder [put]
func (h *Handler) Reorder(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req userpins.ReorderPinsReq
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid input")
	}

	if err := h.svc.Reorder(c.Request().Context(), userID, req.Items); err != nil {
		return classify(c, err)
	}
	return httpUtils.RespondByMessage(c, 200, "PINS_REORDERED", "Pins reordered")
}

// GetPins godoc
// @Summary Get user pins
// @Tags UserPins
// @Produce json
// @Security BearerAuth
// @Param type query string true "Pin type (e.g. CHATLIST)"
// @Success 200 {array} userpins.UserPin
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 500 {object} httpUtils.ErrorResponse
// @Router /user/pins [get]
func (h *Handler) GetPins(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	pinType := c.QueryParam("type")
	if pinType == "" {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "type query param is required")
	}

	pins, err := h.svc.GetPins(c.Request().Context(), userID, pinType)
	if err != nil {
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", err.Error())
	}
	return c.JSON(http.StatusOK, pins)
}

func classify(c echo.Context, err error) error {
	msg := err.Error()
	switch {
	case strings.Contains(msg, "forbidden"):
		return httpUtils.RespondError(c, 403, "FORBIDDEN", msg)
	case strings.Contains(msg, "not found"):
		return httpUtils.RespondError(c, 404, "NOT_FOUND", msg)
	case strings.Contains(msg, "invalid"):
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", msg)
	default:
		return httpUtils.RespondError(c, 500, "INTERNAL_ERROR", msg)
	}
}