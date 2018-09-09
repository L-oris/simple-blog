package env

import (
	"github.com/L-oris/yabb/logger"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type config struct {
	Port string `env:"PORT" envDefault:"8080"`
}

// Vars contains all the env variables in use for the project
var Vars config

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Warning("File .env not found, reading configuration from ENV")
	}

	if err := env.Parse(&Vars); err != nil {
		logger.Log.Fatal("Failed to parse ENV")
	}

	logger.Log.Debug("env variables loaded successfully")
}
