package http

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"

	purchasesvc "github.com/moera-sudo/backend/backend/transaction-service/internal/features/purchase/service"
	httputil "github.com/moera-sudo/backend/backend/transaction-service/internal/pkg/http"
	"github.com/moera-sudo/backend/backend/transaction-service/internal/pkg/utils"
)

// Handler — HTTP-обёртка над purchase service.
type Handler struct {
	svc *purchasesvc.Service
}

func NewHandler(svc *purchasesvc.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes регистрирует все эндпоинты группы /purchase.
// GatewaySecret middleware подключается на group выше.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	g.GET("/packages", h.listPackages)
	g.POST("/initiate", h.initiate)
	g.POST("/:id/confirm", h.confirm)
	g.GET("/history", h.history)
}

// listPackages godoc
// @Summary List SUDA purchase packages
// @Description Public catalog of SUDA purchase packages with fiat prices.
// @Tags Purchase
// @Produce json
// @Security BearerAuth
// @Success 200 {object} PackagesResponse
// @Router /purchase/packages [get]
func (h *Handler) listPackages(c echo.Context) error {
	if _, err := utils.GetUserId(c); err != nil {
		return err
	}

	pkgs := h.svc.ListPackages(c.Request().Context())
	dto := make([]PackageDTO, 0, len(pkgs))
	for _, p := range pkgs {
		dto = append(dto, PackageDTO{
			Code:          p.Code,
			Title:         p.Title,
			SudaAmount:    p.SudaAmount,
			SudaAmountWei: p.SudaAmountWei().String(),
			FiatPrice:     p.FiatPrice,
			FiatCurrency:  p.FiatCurrency,
		})
	}
	return c.JSON(http.StatusOK, PackagesResponse{Packages: dto})
}

// initiate godoc
// @Summary Initiate a SUDA purchase
// @Description Creates a PENDING purchase. Returns purchase_id for follow-up Confirm call.
// @Tags Purchase
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body InitiateRequest true "Package code + payment method"
// @Success 201 {object} InitiateResponse
// @Failure 400 {object} httputil.ErrorResponse "invalid_request | unknown_package"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /purchase/initiate [post]
func (h *Handler) initiate(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	var req InitiateRequest
	if err := c.Bind(&req); err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "malformed JSON")
	}
	if err := c.Validate(&req); err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
	}

	result, err := h.svc.Initiate(c.Request().Context(), purchasesvc.InitiateInput{
		UserID:        userID,
		PackageCode:   req.PackageCode,
		PaymentMethod: req.PaymentMethod,
	})
	if err != nil {
		return mapServiceError(c, err)
	}

	return c.JSON(http.StatusCreated, InitiateResponse{
		PurchaseID:    result.PurchaseID.String(),
		PackageCode:   result.PackageCode,
		SudaAmountWei: result.SudaAmountWei,
		FiatAmount:    result.FiatAmount,
		FiatCurrency:  result.FiatCurrency,
		Status:        result.Status,
		PaymentMethod: result.PaymentMethod,
	})
}

// confirm godoc
// @Summary Confirm a SUDA purchase
// @Description Simulates a payment for an existing PENDING purchase. Treasury broadcasts the SUDA Transfer; observer creates a PURCHASE_COMPLETED WS event for the buyer.
// @Tags Purchase
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Purchase UUID"
// @Param request body ConfirmRequest false "Card details (not validated in simulation)"
// @Success 200 {object} ConfirmResponse
// @Failure 400 {object} httputil.ErrorResponse "invalid_request"
// @Failure 403 {object} httputil.ErrorResponse "forbidden"
// @Failure 404 {object} httputil.ErrorResponse "purchase_not_found"
// @Failure 409 {object} httputil.ErrorResponse "already_handled"
// @Failure 424 {object} httputil.ErrorResponse "wallet_not_found"
// @Failure 429 {object} httputil.ErrorResponse "rate_limited"
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /purchase/{id}/confirm [post]
func (h *Handler) confirm(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	purchaseIDStr := c.Param("id")
	purchaseID, err := uuid.Parse(purchaseIDStr)
	if err != nil {
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", "id must be UUID")
	}

	// Bind/validate чисто символический — card_number/cvv мы не валидируем.
	var req ConfirmRequest
	_ = c.Bind(&req)

	result, err := h.svc.Confirm(c.Request().Context(), purchasesvc.ConfirmInput{
		PurchaseID: purchaseID,
		UserID:     userID,
		RequestIP:  utils.GetClientIP(c),
		UserAgent:  utils.GetUserAgent(c),
	})
	if err != nil {
		return mapServiceError(c, err)
	}

	return c.JSON(http.StatusOK, ConfirmResponse{
		PurchaseID:    result.PurchaseID.String(),
		Status:        result.Status,
		TxHash:        result.TxHash,
		SudaAmountWei: result.SudaAmountWei,
		CompletedAt:   time.Now().UTC().Format(time.RFC3339),
	})
}

// history godoc
// @Summary Get my purchase history
// @Description Paginated list of purchases for the authenticated user (any status).
// @Tags Purchase
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Page size (1..100)" default(50)
// @Param offset query int false "Offset" default(0)
// @Success 200 {object} HistoryResponse
// @Failure 500 {object} httputil.ErrorResponse "internal_error"
// @Router /purchase/history [get]
func (h *Handler) history(c echo.Context) error {
	userID, err := utils.GetUserId(c)
	if err != nil {
		return err
	}

	limit := parseIntDefault(c.QueryParam("limit"), 50, 1, 100)
	offset := parseIntDefault(c.QueryParam("offset"), 0, 0, 1_000_000)

	items, err := h.svc.GetHistory(c.Request().Context(), userID, limit, offset)
	if err != nil {
		log.Error().Err(err).Str("user_id", userID.String()).Msg("purchase history failed")
		return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "failed to load history")
	}

	dto := make([]HistoryItem, 0, len(items))
	for _, p := range items {
		item := HistoryItem{
			ID:            p.ID.String(),
			PackageCode:   p.PackageCode,
			SudaAmountWei: p.SudaAmountWei.String(),
			FiatAmount:    p.FiatAmount,
			FiatCurrency:  p.FiatCurrency,
			Status:        p.Status,
			PaymentMethod: p.PaymentMethod,
			TxHash:        p.TxHash,
			FailureReason: p.FailureReason,
			CreatedAt:     p.CreatedAt.UTC().Format(time.RFC3339),
		}
		if p.CompletedAt != nil {
			item.CompletedAt = p.CompletedAt.UTC().Format(time.RFC3339)
		}
		dto = append(dto, item)
	}

	return c.JSON(http.StatusOK, HistoryResponse{Items: dto})
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

func mapServiceError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, purchasesvc.ErrInvalidInput):
		return httputil.RespondError(c, http.StatusBadRequest, "invalid_request", err.Error())
	case errors.Is(err, purchasesvc.ErrPackageUnknown):
		return httputil.RespondError(c, http.StatusBadRequest, "unknown_package", "this package is not available")
	case errors.Is(err, purchasesvc.ErrNotFound):
		return httputil.RespondError(c, http.StatusNotFound, "purchase_not_found", "purchase does not exist")
	case errors.Is(err, purchasesvc.ErrForbidden):
		return httputil.RespondError(c, http.StatusForbidden, "forbidden", "you do not own this purchase")
	case errors.Is(err, purchasesvc.ErrAlreadyHandled):
		return httputil.RespondError(c, http.StatusConflict, "already_handled", "this purchase was already confirmed or failed")
	case errors.Is(err, purchasesvc.ErrWalletNotFound):
		return httputil.RespondError(c, http.StatusFailedDependency, "wallet_not_found", "your wallet is not initialized")
	case errors.Is(err, purchasesvc.ErrRateLimited):
		return httputil.RespondError(c, http.StatusTooManyRequests, "rate_limited", "too many purchases in a short period")
	default:
		log.Error().Err(err).Msg("purchase service failed")
		return httputil.RespondError(c, http.StatusInternalServerError, "internal_error", "purchase failed")
	}
}

