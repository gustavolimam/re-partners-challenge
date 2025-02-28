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

// CreateOrder cria um pedido para calcular o melhor empacotamento
// @Summary Create Order
// @Description Calculate the best packing options for an order
// @Tags Order
// @Accept json
// @Produce json
// @Param order body models.Order true "Order request body"
// @Success 200 {object} models.OrderResponse
// @Failure 400 {object} map[string]string
// @Router /order [post]
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
