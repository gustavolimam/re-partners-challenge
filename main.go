package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"re-partners-challenge/config"
	"re-partners-challenge/internal/api"
	"re-partners-challenge/internal/clients"
	"re-partners-challenge/internal/constants"
)

func main() {
	config.LoadConfig()
	logLevel, err := zerolog.ParseLevel(config.Cfg.LogLevel)
	if err != nil {
		logLevel = zerolog.InfoLevel
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(logLevel)

	log.Info().Msg("Starting server...")

	cache := loadMemCache()

	e := loadServer(cache)

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%s", config.Cfg.AppPort)))
}

func loadMemCache() *clients.Cache {
	log.Info().Msg("Creating memory cache...")
	cache := clients.NewCache()

	cache.Set(constants.PackSizesCacheKey, constants.PackSizesDefault, 0)

	return cache
}

func loadServer(cache *clients.Cache) *echo.Echo {
	log.Info().Msg("Creating server...")
	e := echo.New()

	e.HideBanner = true

	api.LoadOrderRoutes(e, cache)
	api.LoadPackRoutes(e, cache)

	return e
}
