package config

import (
	"context"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BackAddr    string        `envconfig:"BACK_ADDR"`
	BackTimeout time.Duration `envconfig:"BACK_TIMEOUT"`
	TargetBits  int           `envconfig:"TARGET_BITS"`
}

func NewConfig(ctx context.Context) *Config {
	cfg := &Config{}

	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, ignore outside the local")
	}

	err = envconfig.Process("", cfg)
	if err != nil {
		log.Fatalf("envconfig err: %v", err.Error())
	}
	log.Println("envconfig ok")

	return cfg
}
