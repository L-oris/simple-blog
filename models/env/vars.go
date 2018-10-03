package env

import (
	"github.com/L-oris/yabb/logger"
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type config struct {
	Port                   string `env:"PORT" envDefault:"8080"`
	DB                     string `env:"DB,required"`
	GoogleCloudCredentials string `env:"GOOGLE_APPLICATION_CREDENTIALS,required"`
	BucketName             string `env:"GOOGLE_CLOUD_BUCKET,required"`
}

// Vars stores all env variables in use for the project
var Vars config

func init() {
	if err := godotenv.Load(); err != nil {
		logger.Log.Warningf("File .env not found, reading configuration from ENV: %s", err.Error())
	}

	if err := env.Parse(&Vars); err != nil {
		logger.Log.Fatalf("Failed to parse ENV: %s", err.Error())
	}

	logger.Log.Debug("env variables loaded successfully")
}
