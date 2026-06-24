package http

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"

	httpUtils "github.com/moera-sudo/backend/backend/messenger-service/internal/pkg/http"
)

func classifyAndRespond(c echo.Context, err error) error {
	if err == nil {
		return httpUtils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Unknown error")
	}

	msg := err.Error()

	switch {
	case strings.Contains(msg, "gating"):
		return httpUtils.RespondError(c, http.StatusForbidden, "GATING_REQUIRED", msg)
	case strings.Contains(msg, "forbidden") || strings.Contains(msg, "cannot"):
		return httpUtils.RespondError(c, http.StatusForbidden, "FORBIDDEN", msg)
	case strings.Contains(msg, "not found"):
		return httpUtils.RespondError(c, http.StatusNotFound, "NOT_FOUND", msg)
	case strings.Contains(msg, "not a member"):
		return httpUtils.RespondError(c, http.StatusForbidden, "NOT_MEMBER", msg)
	case strings.Contains(msg, "blocked"):
		return httpUtils.RespondError(c, http.StatusForbidden, "BLOCKED", msg)
	case strings.Contains(msg, "already taken"):
		return httpUtils.RespondError(c, http.StatusConflict, "CONFLICT", msg)
	case strings.Contains(msg, "Invalid"):
		return httpUtils.RespondError(c, http.StatusBadRequest, "BAD_REQUEST", msg)
	default:
		return httpUtils.RespondError(c, http.StatusInternalServerError, "INTERNAL_ERROR", msg)
	}
}