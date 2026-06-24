package http

import (
	"github.com/labstack/echo/v4"
)

// ErrorResponse — стандартизированный JSON-ответ при ошибках.
// Формат идентичен messenger-service.
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

// RespondError отправляет ошибку клиенту.
//
//	errCode — машинный код: "wallet_not_found", "insufficient_balance"
//	message — человекочитаемое описание (опционально)
func RespondError(c echo.Context, code int, errCode string, message string) error {
	return c.JSON(code, ErrorResponse{
		Error:   errCode,
		Message: message,
	})
}