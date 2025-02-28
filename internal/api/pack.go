package api

import (
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
	"net/http"
	"re-partners-challenge/internal/clients"
	"re-partners-challenge/internal/constants"
	"re-partners-challenge/internal/models"
)

type PackI interface {
	UpdatePackSizes(c echo.Context) error
}

type Pack struct {
	cache *clients.Cache
}

func NewPackHandler(cache *clients.Cache) PackI {
	return &Pack{cache}
}

// UpdatePackSizes atualiza os tamanhos de pacotes disponíveis
// @Summary Update Package Sizes
// @Description Update available package sizes for order packing
// @Tags Pack
// @Accept json
// @Produce json
// @Param packSizes body models.PackSizes true "Package sizes request body"
// @Success 200 {string} string "success"
// @Failure 400 {object} map[string]string
// @Router /pack/sizes [put]
func (p *Pack) UpdatePackSizes(c echo.Context) error {
	log.Info().Msg("Updating pack sizes...")

	pack := models.PackSizes{}

	if err := c.Bind(&pack); err != nil {
		log.Error().Err(err).Msg("Failed to bind request body")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// verify if we already have pack sizes set, if not set the pack sizes received from request body
	_, found := p.cache.Get("pack_sizes")
	if !found {
		p.cache.Set(constants.PackSizesCacheKey, pack.Sizes, 0)
		return c.JSON(http.StatusOK, "success")
	}

	// Remove the old version of pack sizes and set the new sizes
	p.cache.Delete(constants.PackSizesCacheKey)
	p.cache.Set(constants.PackSizesCacheKey, pack.Sizes, 0)

	return c.JSON(http.StatusOK, "success")
}
