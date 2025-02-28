package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	echoSwagger "github.com/swaggo/echo-swagger"
	"os"
	"re-partners-challenge/config"
	_ "re-partners-challenge/docs"
	"re-partners-challenge/internal/api"
	"re-partners-challenge/internal/clients"
	"re-partners-challenge/internal/constants"
)

// @title Packing API
// @version 1.0
// @description API for calculating and managing package sizes for orders
// @BasePath /
// @schemes http
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

	// Serve Swagger's documentation
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Load the HTML file
	e.Static("/", "web")

	return e
}
