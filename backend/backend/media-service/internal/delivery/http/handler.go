package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	httpUtils "github.com/moera-sudo/backend/backend/media-service/internal/pkg/http"
	getUser "github.com/moera-sudo/backend/backend/media-service/internal/pkg/utils"
	"github.com/moera-sudo/backend/backend/media-service/internal/service"
	"github.com/moera-sudo/backend/backend/media-service/internal/domain/models"
)

type MediaHandler struct {
	svc *service.MediaService
}

func NewMediaHandler(svc *service.MediaService) *MediaHandler {
	return &MediaHandler{svc: svc}
}

// publicHost returns the client-facing host (incl. port) the request came in on,
// so presigned URLs are signed for the same origin the client can reach. The
// gateway forwards it as X-Forwarded-Host; falls back to the Host header.
// Empty → media service uses the configured S3_PUBLIC_ENDPOINT.
func publicHost(c echo.Context) string {
	if h := c.Request().Header.Get("X-Forwarded-Host"); h != "" {
		return h
	}
	return c.Request().Host
}

func (h *MediaHandler) RegisterRoutes(e *echo.Echo) {
	media := e.Group("/api/media")

	// Upload flow
	media.POST("/upload/init", h.InitUpload)
	media.POST("/:id/confirm", h.ConfirmUpload)

	// URLs (view = inline, download = attachment)
	media.GET("/:id/url", h.GetViewURL)
	media.GET("/:id/download-url", h.GetDownloadURL)

	// Metadata & management
	media.GET("/:id", h.GetMetadata)
	media.POST("/:id/link", h.LinkToEntity)
	media.DELETE("/:id", h.Delete)
}


// InitUpload godoc
// @Summary Initialize a file upload and get presigned URL
// @Tags Media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.InitUploadRequest true "Upload parameters"
// @Success 201 {object} models.InitUploadResponse
// @Failure 400 {object} httpUtils.ErrorResponse
// @Failure 401 {object} httpUtils.ErrorResponse
// @Router /api/media/upload/init [post]
func (h *MediaHandler) InitUpload(c echo.Context) error {
	userID, err := getUser.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	var req models.InitUploadRequest
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, 400, "BAD_REQUEST", "Invalid request body")
	}
	req.OwnerUserID = userID

	resp, err := h.svc.InitUpload(c.Request().Context(), req, publicHost(c))
	if err != nil {
		return httpUtils.RespondError(c, 400, "UPLOAD_INIT_FAILED", err.Error())
	}

	return c.JSON(201, resp)
}

// ConfirmUpload godoc
// @Summary Confirm that file was uploaded to S3
// @Tags Media
// @Produce json
// @Security BearerAuth
// @Param id path string true "Media ID (UUID)"
// @Success 200 {object} models.ConfirmUploadResponse
// @Failure 401 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Failure 404 {object} httpUtils.ErrorResponse
// @Router /api/media/{id}/confirm [post]
func (h *MediaHandler) ConfirmUpload(c echo.Context) error {
	userID, err := getUser.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	resp, err := h.svc.ConfirmUpload(c.Request().Context(), c.Param("id"), userID)
	if err != nil {
		code, errCode := classifyError(err)
		return httpUtils.RespondError(c, code, errCode, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

// GetViewURL godoc
// @Summary Get a presigned URL for inline viewing
// @Tags Media
// @Produce json
// @Security BearerAuth
// @Param id path string true "Media ID (UUID)"
// @Success 200 {object} models.GetURLResponse
// @Failure 401 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Failure 404 {object} httpUtils.ErrorResponse
// @Router /api/media/{id}/url [get]
func (h *MediaHandler) GetViewURL(c echo.Context) error {
	userID, err := getUser.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	resp, err := h.svc.GetViewURL(c.Request().Context(), c.Param("id"), userID, publicHost(c))
	if err != nil {
		code, errCode := classifyError(err)
		return httpUtils.RespondError(c, code, errCode, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}

// GetDownloadURL godoc
// @Summary Get a presigned URL for downloading
// @Tags Media
// @Produce json
// @Security BearerAuth
// @Param id path string true "Media ID (UUID)"
// @Success 200 {object} models.GetURLResponse
// @Failure 401 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Failure 404 {object} httpUtils.ErrorResponse
// @Router /api/media/{id}/download-url [get]
func (h *MediaHandler) GetDownloadURL(c echo.Context) error {
	userID, err := getUser.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	resp, err := h.svc.GetDownloadURL(c.Request().Context(), c.Param("id"), userID, publicHost(c))
	if err != nil {
		code, errCode := classifyError(err)
		return httpUtils.RespondError(c, code, errCode, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}


// GetMetadata godoc
// @Summary Get media file metadata
// @Tags Media
// @Produce json
// @Security BearerAuth
// @Param id path string true "Media ID (UUID)"
// @Success 200 {object} models.MediaMetadata
// @Failure 401 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Failure 404 {object} httpUtils.ErrorResponse
// @Router /api/media/{id} [get]
func (h *MediaHandler) GetMetadata(c echo.Context) error {
	userID, err := getUser.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Invalid User ID")
	}

	resp, err := h.svc.GetMetadata(c.Request().Context(), c.Param("id"), userID)
	if err != nil {
		code, errCode := classifyError(err)
		return httpUtils.RespondError(c, code, errCode, err.Error())
	}

	return c.JSON(http.StatusOK, resp)
}


// LinkToEntity godoc
// @Summary Link media to an entity (message, chat avatar, etc.)
// @Tags Media
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Media ID (UUID)"
// @Param request body object true "entity_type and entity_id"
// @Success 200 {object} httpUtils.MessageRespond
// @Failure 400 {object} httpUtils.ErrorResponse
// @Router /api/media/{id}/link [post]
func (h *MediaHandler) LinkToEntity(c echo.Context) error {
	var req struct {
		EntityType string `json:"entity_type" validate:"required"`
		EntityID   string `json:"entity_id" validate:"required"`
	}
	if err := c.Bind(&req); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "Invalid request body")
	}

	if req.EntityType == "" || req.EntityID == "" {
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", "entity_type and entity_id are required")
	}

	if err := h.svc.LinkToEntity(c.Request().Context(), c.Param("id"), req.EntityType, req.EntityID); err != nil {
		return httpUtils.RespondError(c, http.StatusBadRequest, "LINK_FAILED", err.Error())
	}

	return httpUtils.RespondByMessage(c, http.StatusOK, "LINKED", "Media linked to entity successfully")
}


// Delete godoc
// @Summary Delete a media file
// @Tags Media
// @Security BearerAuth
// @Param id path string true "Media ID (UUID)"
// @Success 204 "No Content"
// @Failure 401 {object} httpUtils.ErrorResponse
// @Failure 403 {object} httpUtils.ErrorResponse
// @Failure 404 {object} httpUtils.ErrorResponse
// @Router /api/media/{id} [delete]
func (h *MediaHandler) Delete(c echo.Context) error {
	userID, err := getUser.GetUserId(c)
	if err != nil {
		return httpUtils.RespondError(c, 401, "UNAUTHORIZED", "Missing user ID")
	}

	if err := h.svc.Delete(c.Request().Context(), c.Param("id"), userID); err != nil {
		code, errCode := classifyError(err)
		return httpUtils.RespondError(c, code, errCode, err.Error())
	}

	return c.NoContent(http.StatusNoContent)
}

func classifyError(err error) (int, string) {
	if err == nil {
		return http.StatusInternalServerError, "INTERNAL_ERROR"
	}
	msg := err.Error()

	if strings.HasPrefix(msg, "forbidden") {
		return http.StatusForbidden, "FORBIDDEN"
	}
	if strings.Contains(msg, "not found") {
		return http.StatusNotFound, "NOT_FOUND"
	}
	if strings.Contains(msg, "deleted") {
		return http.StatusGone, "DELETED"
	}
	if strings.Contains(msg, "not ready") || strings.Contains(msg, "PENDING") {
		return http.StatusConflict, "NOT_READY"
	}

	return http.StatusBadRequest, "BAD_REQUEST"
}