package http

import (
	"github.com/labstack/echo/v4"
)

// MessageRespond — стандартизированный JSON-ответ для success-операций.
// Формат идентичен messenger-service, чтобы фронтенд получал одинаковые
// структуры от всех сервисов и мог переиспользовать общий код парсинга.
type MessageRespond struct {
	Action  string `json:"action"`
	Message string `json:"message,omitempty"`
}

// RespondByMessage отправляет клиенту success-ответ с action-кодом и
// опциональным человекочитаемым сообщением.
//
//	action  — машинный код события: "wallet_created", "transfer_submitted"
//	message — человекочитаемое описание (опционально)
func RespondByMessage(c echo.Context, code int, action string, message string) error {
	return c.JSON(code, MessageRespond{
		Action:  action,
		Message: message,
	})
}