package http

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	"github.com/moera-sudo/backend/backend/transaction-service/internal/features/gating"
	gatingsvc "github.com/moera-sudo/backend/backend/transaction-service/internal/features/gating/service"
	httputil "github.com/moera-sudo/backend/backend/transaction-service/internal/pkg/http"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/pkg/utils"
)

// Handler — HTTP-обёртка над gating service.
type Handler struct {
	svc *gatingsvc.Service
}

func NewHandler(svc *gatingsvc.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes регистрирует эндпоинты группы /gating.
// GatewaySecret middleware подключается на group выше.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("/rule", h.createRule)
	g.GET("/rule/:chat_id", h.getRule)
	g.DELETE("/rule/:chat_id", h.deleteRule)
}

// createRule godoc
// @Summary Create or update a gating rule
// @Description Sets a token-gating rule for a channel (minimum SUDA balance and/or required NFT collection). Channel OWNER only.
// @Tags Gating
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateRuleRequest true "Gating rule parameters"
// @Success 200 {object} RuleResponse
// @Failure 400 {object} httputil.ErrorResponse "invalid_request"
// @Failure 403 {object} httputil.ErrorResponse "not_owner"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /gating/rule [post]
func (h *Handler) createRule(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req CreateRuleRequest
	if err := c.Bind(&req); err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "malformed JSON")
	}
	if err := c.Validate(&req); err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
	}

	chatID, err := uuid.Parse(req.ChatID)
	if err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "chat_id must be UUID")
	}

	var nftCollectionID *uuid.UUID
	if req.RequiredNFTCollectionID != "" {
		parsed, err := uuid.Parse(req.RequiredNFTCollectionID)
		if err != nil {
			return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "required_nft_collection_id must be UUID")
		}
		nftCollectionID = &parsed
	}

	err = h.svc.CreateRule(c.Request().Context(), gatingsvc.RuleInput{
		OwnerID:                 userID,
		ChatID:                  chatID,
		MinSudaBalanceWei:       req.MinSudaBalanceWei,
		RequiredNFTCollectionID: nftCollectionID,
		SubscriptionPriceWei:    req.SubscriptionPriceWei,
	})
	if err != nil {
		return mapGatingError(c, err)
	}

	rule, err := h.svc.GetRule(c.Request().Context(), chatID)
	if err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Msg("gating: read-back after create failed")
		return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "rule saved but failed to read it back")
	}
	return c.JSON(http.StatusOK, toRuleResponse(rule))
}

// getRule godoc
// @Summary Get gating rule for a channel
// @Description Returns the current gating rule (if any). Available to any authenticated user — the threshold is public.
// @Tags Gating
// @Produce json
// @Security BearerAuth
// @Param chat_id path string true "Channel UUID"
// @Success 200 {object} RuleResponse
// @Failure 400 {object} httputil.ErrorResponse "invalid_request"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /gating/rule/{chat_id} [get]
func (h *Handler) getRule(c echo.Context) error {
	if _, err := utils.GetUserId(c); err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "chat_id must be UUID")
	}

	rule, err := h.svc.GetRule(c.Request().Context(), chatID)
	if err != nil {
		log.Error().Err(err).Str("chat_id", chatID.String()).Msg("gating: get rule failed")
		return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "failed to load rule")
	}
	return c.JSON(http.StatusOK, toRuleResponse(rule))
}

// deleteRule godoc
// @Summary Delete gating rule
// @Description Removes the gating rule for a channel. Channel OWNER only.
// @Tags Gating
// @Produce json
// @Security BearerAuth
// @Param chat_id path string true "Channel UUID"
// @Success 200 {object} httputil.MessageRespond
// @Failure 400 {object} httputil.ErrorResponse "invalid_request"
// @Failure 403 {object} httputil.ErrorResponse "not_owner"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /gating/rule/{chat_id} [delete]
func (h *Handler) deleteRule(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	chatID, err := uuid.Parse(c.Param("chat_id"))
	if err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "chat_id must be UUID")
	}

	if err := h.svc.DeleteRule(c.Request().Context(), userID, chatID); err != nil {
		return mapGatingError(c, err)
	}
	return httputil.RespondByMessage(c, http.StatusOK, "deleted", "gating rule removed")
}

func toRuleResponse(r *gatingsvc.Rule) RuleResponse {
	resp := RuleResponse{
		ChatID:               r.ChatID.String(),
		MinSudaBalanceWei:    r.MinSudaBalanceWei,
		SubscriptionPriceWei: r.SubscriptionPriceWei,
		HasRule:              r.HasRule,
	}
	if r.RequiredNFTCollectionID != nil {
		resp.RequiredNFTCollectionID = r.RequiredNFTCollectionID.String()
	}
	return resp
}

func mapGatingError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, gating.ErrInvalidInput):
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
	case errors.Is(err, gating.ErrForbidden):
		return httputil.RespondError(c, http.StatusForbidden, "not_owner", "only the channel owner can manage gating rules")
	default:
		log.Error().Err(err).Msg("gating operation failed")
		return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "gating operation failed")
	}
}
