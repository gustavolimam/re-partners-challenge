package config

import (
	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
	"log"
)

type config struct {
	AppPort  string `env:"APP_PORT"`
	LogLevel string `env:"LOG_LEVEL"`
}

var Cfg config

func LoadConfig() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
	_, err := env.UnmarshalFromEnviron(&Cfg)
	if err != nil {
		log.Fatal(err)
	}
}
