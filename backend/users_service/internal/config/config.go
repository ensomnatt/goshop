package config

import (
	"log"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Port        string `env:"PORT" envDefault:"8080"`
	DatabaseUrl string `env:"DATABASE_URL,required"`
}

func Load() *Config {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	return cfg
}
