package http

import "github.com/labstack/echo/v4"

// Handler — заглушка HTTP-handler для quest.
type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

// RegisterRoutes — пустая регистрация. Маршруты появятся в итерации Г.
func (h *Handler) RegisterRoutes(g *echo.Group) {
	_ = g
}