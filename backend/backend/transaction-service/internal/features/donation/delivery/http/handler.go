package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	donationsvc "github.com/moera-sudo/backend/backend/transaction-service/internal/features/donation/service"
	httputil "github.com/moera-sudo/backend/backend/transaction-service/internal/pkg/http"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/pkg/utils"
)

// Handler — HTTP-обёртка над donation service.
type Handler struct {
	svc *donationsvc.Service
}

func NewHandler(svc *donationsvc.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes регистрирует эндпоинты группы /donate.
// GatewaySecret middleware подключается на group выше.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.POST("", h.donate)
}

// donate godoc
// @Summary Donate SUDA to a user or channel
// @Description Send a SUDA donation. Exactly one of to_username or to_channel_id must be set. Returns tx_hash immediately; observer indexes donation asynchronously and emits DONATION_SENT/DONATION_RECEIVED WS events.
// @Tags Donation
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body DonateRequest true "Donation parameters"
// @Success 202 {object} DonateResponse
// @Failure 400 {object} httputil.ErrorResponse "invalid_request | cannot_donate_to_self | recipient_has_no_wallet | insufficient_balance"
// @Failure 403 {object} httputil.ErrorResponse "sender_not_in_chat | recipient_not_in_chat"
// @Failure 404 {object} httputil.ErrorResponse "recipient_not_found | channel_wallet_not_found | wallet_not_found | chat_not_found"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /donate [post]
func (h *Handler) donate(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req DonateRequest
	if err := c.Bind(&req); err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "malformed JSON")
	}
	if err := c.Validate(&req); err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
	}

	var toChannelID uuid.UUID
	if req.ToChannelID != "" {
		parsed, err := uuid.Parse(req.ToChannelID)
		if err != nil {
			return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "to_channel_id must be UUID")
		}
		toChannelID = parsed
	}

	var chatID uuid.UUID
	if req.ChatID != "" {
		parsed, err := uuid.Parse(req.ChatID)
		if err != nil {
			return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "chat_id must be UUID")
		}
		chatID = parsed
	}

	result, err := h.svc.Donate(c.Request().Context(), donationsvc.DonationInput{
		FromUserID:  userID,
		ToUsername:  req.ToUsername,
		ToChannelID: toChannelID,
		AmountWei:   req.AmountWei,
		Message:     req.Message,
		ChatID:      chatID,
		RequestIP:   utils.GetClientIP(c),
		UserAgent:   utils.GetUserAgent(c),
	})
	if err != nil {
		return mapDonationError(c, err)
	}

	return c.JSON(http.StatusAccepted, DonateResponse{
		TxHash:      result.TxHash,
		FromAddress: result.FromAddress,
		ToAddress:   result.ToAddress,
		AmountWei:   result.AmountWei,
		IsChannel:   result.IsChannel,
		ToUserID:    result.ToUserID,
		ToChannelID: result.ToChannelID,
		SubmittedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

// mapDonationError приводит доменные ошибки сервиса к HTTP-кодам.
func mapDonationError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, donationsvc.ErrInvalidInput):
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
	case errors.Is(err, donationsvc.ErrSelfDonation):
		return httputil.RespondError(c, http.StatusBadRequest, "cannot_donate_to_self", "you cannot donate to yourself")
	case errors.Is(err, donationsvc.ErrRecipientNotFound):
		return httputil.RespondError(c, http.StatusNotFound, "recipient_not_found", "no user with this username")
	case errors.Is(err, donationsvc.ErrRecipientHasNoWallet):
		return httputil.RespondError(c, http.StatusBadRequest, "recipient_has_no_wallet", "recipient does not have a wallet yet")
	case errors.Is(err, donationsvc.ErrChannelWalletNotFound):
		return httputil.RespondError(c, http.StatusNotFound, "channel_wallet_not_found", "this channel has no treasury wallet yet")
	case errors.Is(err, donationsvc.ErrWalletNotFound):
		return httputil.RespondError(c, http.StatusNotFound, "wallet_not_found", "your wallet is not initialized")
	case errors.Is(err, donationsvc.ErrInsufficientBalance):
		return httputil.RespondError(c, http.StatusBadRequest, "insufficient_balance", "your SUDA balance is too low")
	case errors.Is(err, donationsvc.ErrChatNotFound):
		return httputil.RespondError(c, http.StatusNotFound, "chat_not_found", "chat does not exist")
	case errors.Is(err, donationsvc.ErrSenderNotInChat):
		return httputil.RespondError(c, http.StatusForbidden, "sender_not_in_chat", "you are not a member of this chat")
	case errors.Is(err, donationsvc.ErrRecipientNotInChat):
		return httputil.RespondError(c, http.StatusForbidden, "recipient_not_in_chat", "recipient is not a member of this chat")
	default:
		log.Error().Err(err).Msg("donation failed")
		return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "donation failed")
	}
}
