package api

import (
	"github.com/labstack/echo/v4"
	"re-partners-challenge/internal/clients"
)

func LoadOrderRoutes(e *echo.Echo, cache *clients.Cache) {
	handler := NewOrderHandler(cache)

	e.POST("/order", handler.CreateOrder)
}

func LoadPackRoutes(e *echo.Echo, cache *clients.Cache) {
	handler := NewPackHandler(cache)

	e.PUT("/pack/sizes", handler.UpdatePackSizes)
}
