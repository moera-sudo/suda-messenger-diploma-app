package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	walletsvc "github.com/moera-sudo/backend/backend/transaction-service/internal/features/wallet/service"
	httputil "github.com/moera-sudo/backend/backend/transaction-service/internal/pkg/http"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/pkg/utils"
)

// Handler — HTTP-обёртка над wallet service.
type Handler struct {
	svc *walletsvc.Service
}

func NewHandler(svc *walletsvc.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes регистрирует все эндпоинты группы /wallet.
// GatewaySecret middleware должен быть подключён на group выше.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("/me", h.getMe)
	g.POST("/resolve", h.resolve)
	g.POST("/transfer", h.transfer)
	g.GET("/history", h.getHistory)
	g.GET("/channel/:channel_id", h.getChannelWallet)
	g.GET("/channel/:channel_id/treasury", h.getChannelTreasury)
	g.GET("/channel/:channel_id/donations", h.getChannelDonations)
	g.POST("/channel/:channel_id/withdraw", h.withdrawTreasury)
}

// getMe godoc
// @Summary Get my wallet
// @Description Returns the authenticated user's custodial wallet address and on-chain SUDA balance.
// @Tags Wallet
// @Produce json
// @Security BearerAuth
// @Success 200 {object} MeResponse
// @Failure 404 {object} httputil.ErrorResponse "wallet_not_found"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /wallet/me [get]
func (h *Handler) getMe(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	addr, balanceWei, err := h.svc.GetMyWalletAndBalance(c.Request().Context(), userID)
	if err != nil {
		if errors.Is(err, walletsvc.ErrWalletNotFound) {
			return httputil.RespondError(c, http.StatusNotFound, "wallet_not_found", "no wallet for this user")
		}
		log.Error().Err(err).Str("user_id", userID.String()).Msg("getMe failed")
		return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "failed to load wallet")
	}

	return c.JSON(http.StatusOK, MeResponse{
		Address:        addr,
		SudaBalanceWei: balanceWei,
		Decimals:       18,
	})
}

// resolve godoc
// @Summary Resolve username to wallet
// @Description Looks up a user by username and returns their user_id, display_name and wallet_address (if any).
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ResolveRequest true "Username to resolve"
// @Success 200 {object} ResolveResponse
// @Failure 400 {object} httputil.ErrorResponse "invalid_request"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /wallet/resolve [post]
func (h *Handler) resolve(c echo.Context) error {
	if _, err := utils.GetUserId(c); err != nil {
		return err
	}

	var req ResolveRequest
	if err := c.Bind(&req); err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "malformed JSON")
	}
	if err := c.Validate(&req); err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
	}

	resolved, err := h.svc.ResolveUsername(c.Request().Context(), req.Username)
	if err != nil {
		log.Error().Err(err).Str("username", req.Username).Msg("resolve failed")
		return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "resolve failed")
	}

	if !resolved.Found {
		return c.JSON(http.StatusOK, ResolveResponse{Found: false})
	}
	return c.JSON(http.StatusOK, ResolveResponse{
		Found:         true,
		UserID:        resolved.UserID,
		DisplayName:   resolved.DisplayName,
		WalletAddress: resolved.WalletAddress,
	})
}

// transfer godoc
// @Summary Transfer SUDA to a user
// @Description Custodial P2P transfer of SUDA tokens. Returns tx_hash immediately (async confirmation via WS event SUDA_RECEIVED/SUDA_SENT). If chat_id is provided, a system message is posted to that chat.
// @Tags Wallet
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body TransferRequest true "Transfer parameters"
// @Success 202 {object} TransferResponse
// @Failure 400 {object} httputil.ErrorResponse "invalid_request | cannot_transfer_to_self | recipient_has_no_wallet | insufficient_balance"
// @Failure 403 {object} httputil.ErrorResponse "sender_not_in_chat | recipient_not_in_chat"
// @Failure 404 {object} httputil.ErrorResponse "recipient_not_found | wallet_not_found | chat_not_found"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /wallet/transfer [post]
func (h *Handler) transfer(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req TransferRequest
	if err := c.Bind(&req); err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "malformed JSON")
	}
	if err := c.Validate(&req); err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
	}

	// Парсим опциональный chat_id для transfer-in-chat.
	var chatID uuid.UUID
	if req.ChatID != "" {
		parsed, err := uuid.Parse(req.ChatID)
		if err != nil {
			return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "chat_id must be UUID")
		}
		chatID = parsed
	}

	result, err := h.svc.Transfer(c.Request().Context(), walletsvc.TransferInput{
		FromUserID: userID,
		ToUsername: req.ToUsername,
		AmountWei:  req.AmountWei,
		Note:       req.Note,
		ChatID:     chatID, // uuid.Nil если не задан — обычный перевод
		RequestIP:  utils.GetClientIP(c),
		UserAgent:  utils.GetUserAgent(c),
	})

	if err != nil {
		return mapTransferError(c, err)
	}

	return c.JSON(http.StatusAccepted, TransferResponse{
		TxHash:      result.TxHash,
		FromAddress: result.FromAddress,
		ToAddress:   result.ToAddress,
		ToUserID:    result.ToUserID,
		AmountWei:   result.AmountWei,
		SubmittedAt: time.Now().UTC().Format(time.RFC3339),
		ChatID:      result.ChatID,
	})
}

// getHistory godoc
// @Summary Get my transaction history
// @Description Returns confirmed on-chain transactions involving the user's wallet (transfers, donations, purchases).
// @Tags Wallet
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Page size (1..100)" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} HistoryResponse
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /wallet/history [get]
func (h *Handler) getHistory(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	limit := parseIntDefault(c.QueryParam("limit"), 50, 1, 100)
	offset := parseIntDefault(c.QueryParam("offset"), 0, 0, 1_000_000)

	items, err := h.svc.GetHistory(c.Request().Context(), userID, limit, offset)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("get history failed")
		return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "failed to load history")
	}

	dtoItems := make([]HistoryItem, 0, len(items))
	for _, it := range items {
		dto := HistoryItem{
			ID:          it.ID.String(),
			TxHash:      it.TxHash,
			LogIndex:    it.LogIndex,
			FromAddress: it.FromAddress,
			ToAddress:   it.ToAddress,
			AmountWei:   it.Amount.String(),
			Type:        it.Type,
			Status:      it.Status,
			BlockNumber: it.BlockNumber,
			ConfirmedAt: it.ConfirmedAt.UTC().Format(time.RFC3339),
		}
		if it.FromUserID != nil {
			dto.FromUserID = it.FromUserID.String()
		}
		if it.ToUserID != nil {
			dto.ToUserID = it.ToUserID.String()
		}
		if it.RelatedEntityType != nil {
			dto.RelatedEntityType = *it.RelatedEntityType
		}
		if it.RelatedEntityID != nil {
			dto.RelatedEntityID = it.RelatedEntityID.String()
		}
		if it.Note != nil {
			dto.Note = *it.Note
		}
		dtoItems = append(dtoItems, dto)
	}

	return c.JSON(http.StatusOK, HistoryResponse{Items: dtoItems})
}

// getChannelWallet godoc
// @Summary Get channel treasury wallet
// @Description Returns the channel's treasury wallet address and on-chain SUDA balance. Available to channel OWNER or ADMIN.
// @Tags Wallet
// @Produce json
// @Security BearerAuth
// @Param channel_id path string true "Channel UUID"
// @Success 200 {object} ChannelWalletResponse
// @Failure 400 {object} httputil.ErrorResponse "invalid_request"
// @Failure 403 {object} httputil.ErrorResponse "forbidden"
// @Failure 404 {object} httputil.ErrorResponse "wallet_not_found"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /wallet/channel/{channel_id} [get]
func (h *Handler) getChannelWallet(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	channelIDStr := c.Param("channel_id")
	channelID, err := uuid.Parse(channelIDStr)
	if err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "channel_id must be UUID")
	}

	addr, balanceWei, err := h.svc.GetChannelWalletAndBalance(c.Request().Context(), userID, channelID)
	if err != nil {
		switch {
		case errors.Is(err, walletsvc.ErrWalletNotFound):
			return httputil.RespondError(c, http.StatusNotFound, "wallet_not_found", "channel has no wallet")
		case errors.Is(err, walletsvc.ErrForbidden):
			return httputil.RespondError(c, http.StatusForbidden, "forbidden", "you are not allowed to view this channel's wallet")
		default:
			log.Error().Err(err).Str("user_id", userID.String()).Str("channel_id", channelID.String()).Msg("get channel wallet failed")
			return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "failed to load channel wallet")
		}
	}

	return c.JSON(http.StatusOK, ChannelWalletResponse{
		ChannelID:      channelID.String(),
		Address:        addr,
		SudaBalanceWei: balanceWei,
		Decimals:       18,
	})
}

// getChannelTreasury godoc
// @Summary Get channel treasury statistics
// @Description Aggregated donation stats for a channel: on-chain balance, total donations count/amount, top donors. OWNER/ADMIN only.
// @Tags Treasury
// @Produce json
// @Security BearerAuth
// @Param channel_id path string true "Channel UUID"
// @Success 200 {object} TreasuryResponse
// @Failure 400 {object} httputil.ErrorResponse "invalid_request"
// @Failure 403 {object} httputil.ErrorResponse "forbidden"
// @Failure 404 {object} httputil.ErrorResponse "wallet_not_found"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /wallet/channel/{channel_id}/treasury [get]
func (h *Handler) getChannelTreasury(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	channelID, err := uuid.Parse(c.Param("channel_id"))
	if err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "channel_id must be UUID")
	}

	stats, err := h.svc.GetChannelTreasury(c.Request().Context(), userID, channelID)
	if err != nil {
		return mapTreasuryError(c, err)
	}

	topDonors := make([]TopDonorDTO, 0, len(stats.TopDonors))
	for _, d := range stats.TopDonors {
		topDonors = append(topDonors, TopDonorDTO{
			UserID:        d.UserID,
			Username:      d.Username,
			DisplayName:   d.DisplayName,
			AmountWei:     d.TotalWei,
			DonationCount: d.DonationCount,
		})
	}

	return c.JSON(http.StatusOK, TreasuryResponse{
		ChannelID:           stats.ChannelID,
		Address:             stats.Address,
		SudaBalanceWei:      stats.SudaBalanceWei,
		Decimals:            18,
		TotalDonationsCount: stats.TotalDonationsCount,
		TotalDonationsWei:   stats.TotalDonationsWei,
		TopDonors:           topDonors,
	})
}

// getChannelDonations godoc
// @Summary List channel donations
// @Description Paginated list of donations to a channel with resolved usernames. OWNER/ADMIN only.
// @Tags Treasury
// @Produce json
// @Security BearerAuth
// @Param channel_id path string true "Channel UUID"
// @Param limit query int false "Page size (1..100)" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} DonationsListResponse
// @Failure 400 {object} httputil.ErrorResponse "invalid_request"
// @Failure 403 {object} httputil.ErrorResponse "forbidden"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /wallet/channel/{channel_id}/donations [get]
func (h *Handler) getChannelDonations(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	channelID, err := uuid.Parse(c.Param("channel_id"))
	if err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "channel_id must be UUID")
	}

	limit := parseIntDefault(c.QueryParam("limit"), 50, 1, 100)
	offset := parseIntDefault(c.QueryParam("offset"), 0, 0, 1_000_000)

	list, err := h.svc.GetChannelDonations(c.Request().Context(), userID, channelID, limit, offset)
	if err != nil {
		return mapTreasuryError(c, err)
	}

	items := make([]DonationDTO, 0, len(list.Items))
	for _, d := range list.Items {
		items = append(items, DonationDTO{
			ID:              d.ID,
			FromUserID:      d.FromUserID,
			FromUsername:    d.FromUsername,
			FromDisplayName: d.FromDisplayName,
			FromAddress:     d.FromAddress,
			AmountWei:       d.AmountWei,
			Message:         d.Message,
			TxHash:          d.TxHash,
			CreatedAt:       d.CreatedAt.UTC().Format(time.RFC3339),
		})
	}

	return c.JSON(http.StatusOK, DonationsListResponse{Items: items, Total: list.Total})
}

// withdrawTreasury godoc
// @Summary Withdraw funds from a channel treasury
// @Description Channel OWNER withdraws SUDA from the channel treasury to their own wallet. Returns tx_hash immediately (async confirmation). OWNER only.
// @Tags Treasury
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param channel_id path string true "Channel UUID"
// @Param request body WithdrawRequest true "Withdraw parameters"
// @Success 202 {object} WithdrawResponse
// @Failure 400 {object} httputil.ErrorResponse "invalid_request | insufficient_balance"
// @Failure 403 {object} httputil.ErrorResponse "forbidden"
// @Failure 404 {object} httputil.ErrorResponse "wallet_not_found"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /wallet/channel/{channel_id}/withdraw [post]
func (h *Handler) withdrawTreasury(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	channelID, err := uuid.Parse(c.Param("channel_id"))
	if err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "channel_id must be UUID")
	}

	var req WithdrawRequest
	if err := c.Bind(&req); err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "malformed JSON")
	}
	if err := c.Validate(&req); err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
	}

	result, err := h.svc.WithdrawTreasury(
		c.Request().Context(), userID, channelID, req.AmountWei,
		utils.GetClientIP(c), utils.GetUserAgent(c),
	)
	if err != nil {
		switch {
		case errors.Is(err, walletsvc.ErrInvalidInput):
			return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
		case errors.Is(err, walletsvc.ErrInsufficientBalance):
			return httputil.RespondError(c, http.StatusBadRequest, "insufficient_balance", "treasury balance is too low")
		case errors.Is(err, walletsvc.ErrForbidden):
			return httputil.RespondError(c, http.StatusForbidden, "forbidden", "only the channel owner can withdraw")
		case errors.Is(err, walletsvc.ErrWalletNotFound):
			return httputil.RespondError(c, http.StatusNotFound, "wallet_not_found", "channel treasury or owner wallet not found")
		default:
			log.Error().Err(err).Str("user_id", userID.String()).Str("channel_id", channelID.String()).Msg("treasury withdraw failed")
			return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "withdraw failed")
		}
	}

	return c.JSON(http.StatusAccepted, WithdrawResponse{
		TxHash:      result.TxHash,
		FromAddress: result.FromAddress,
		ToAddress:   result.ToAddress,
		AmountWei:   result.AmountWei,
		SubmittedAt: time.Now().UTC().Format(time.RFC3339),
	})
}

// mapTreasuryError приводит доменные ошибки treasury к HTTP-кодам.
func mapTreasuryError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, walletsvc.ErrForbidden):
		return httputil.RespondError(c, http.StatusForbidden, "forbidden", "you are not allowed to view this channel's treasury")
	case errors.Is(err, walletsvc.ErrWalletNotFound):
		return httputil.RespondError(c, http.StatusNotFound, "wallet_not_found", "channel has no treasury wallet")
	default:
		log.Error().Err(err).Msg("treasury request failed")
		return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "failed to load treasury data")
	}
}

// Helpers

func parseIntDefault(s string, def, min, max int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	if n < min {
		return min
	}
	if n > max {
		return max
	}
	return n
}

// mapTransferError приводит доменные ошибки сервиса к HTTP-кодам.
func mapTransferError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, walletsvc.ErrInvalidInput):
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
	case errors.Is(err, walletsvc.ErrSelfTransfer):
		return httputil.RespondError(c, http.StatusBadRequest, "cannot_transfer_to_self", "you cannot transfer SUDA to yourself")
	case errors.Is(err, walletsvc.ErrRecipientNotFound):
		return httputil.RespondError(c, http.StatusNotFound, "recipient_not_found", "no user with this username")
	case errors.Is(err, walletsvc.ErrRecipientHasNoWallet):
		return httputil.RespondError(c, http.StatusBadRequest, "recipient_has_no_wallet", "recipient does not have a wallet yet")
	case errors.Is(err, walletsvc.ErrInsufficientBalance):
		return httputil.RespondError(c, http.StatusBadRequest, "insufficient_balance", "your SUDA balance is too low")
	case errors.Is(err, walletsvc.ErrWalletNotFound):
		return httputil.RespondError(c, http.StatusNotFound, "wallet_not_found", "your wallet is not initialized")

	// ── transfer-in-chat ошибки ──
	case errors.Is(err, walletsvc.ErrChatNotFound):
		return httputil.RespondError(c, http.StatusNotFound, "chat_not_found", "chat does not exist")
	case errors.Is(err, walletsvc.ErrSenderNotInChat):
		return httputil.RespondError(c, http.StatusForbidden, "sender_not_in_chat", "you are not a member of this chat")
	case errors.Is(err, walletsvc.ErrRecipientNotInChat):
		return httputil.RespondError(c, http.StatusForbidden, "recipient_not_in_chat", "recipient is not a member of this chat")

	default:
		log.Error().Err(err).Msg("transfer failed")
		return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "transfer failed")
	}
}