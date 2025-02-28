package api

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"re-partners-challenge/internal/clients"
	"re-partners-challenge/internal/models"
	"re-partners-challenge/internal/services"
)

type OrderI interface {
	CreateOrder(c echo.Context) error
}

type Order struct {
	service services.OrderI
}

func NewOrderHandler(cache *clients.Cache) OrderI {
	return &Order{
		service: services.NewOrderService(cache),
	}
}

func (o *Order) CreateOrder(c echo.Context) error {
	log.Info().Msg("Receiving a new order...")

	order := models.Order{}

	if err := c.Bind(&order); err != nil {
		log.Error().Err(err).Msg("Failed to bind request body")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	packs, err := o.service.CalculateOrderPackQty(order)
	if err != nil {
		log.Error().Err(err).Msg("Failed to calculate pack qty")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, models.OrderResponse{
		OrderItems: order.Items,
		OrderPacks: packs,
	})
}
