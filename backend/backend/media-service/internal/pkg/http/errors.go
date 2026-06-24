package http

import (
	"github.com/labstack/echo/v4"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Message string `json:"message,omitempty"`
}

func RespondError (c echo.Context, code int, errCode string, message string) error {
	return c.JSON(code, ErrorResponse{
		Error: errCode, 
		Message: message,
	})
}

