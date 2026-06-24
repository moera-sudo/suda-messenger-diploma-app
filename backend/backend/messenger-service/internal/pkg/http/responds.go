package http

import (
	"github.com/labstack/echo/v4"
)

type MessageRespond struct {
	Action string `json:"action"`
	Message string `json:"message,omitempty"`
}

func RespondByMessage(c echo.Context, code int, action string, message string) error {
	return c.JSON(code, MessageRespond{
		Action: action,
		Message: message,
	})
}