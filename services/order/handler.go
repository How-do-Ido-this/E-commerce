package order

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler maneja las peticiones HTTP para el servicio de órdenes.
type Handler struct {
	service *Service
}

// NewHandler crea un nuevo handler para el servicio de órdenes.
func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

// CreateOrder es el endpoint para crear una orden.
func (h *Handler) CreateOrder(c echo.Context) error {
	var o Order
	if err := c.Bind(&o); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}

	if err := h.service.CreateOrder(c.Request().Context(), &o); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, o)
}
