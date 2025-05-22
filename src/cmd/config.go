package main

import (
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	HTTPServer HTTPServer `envPrefix:"HTTP_SERVER_"`
	DB         DB         `envPrefix:"DB_"`
}

type HTTPServer struct {
	ListenAddr string `env:"LISTEN_ADDR,notEmpty"`
}

type DB struct {
	Host     string `env:"HOST,notEmpty"`
	Port     int    `env:"PORT,notEmpty"`
	User     string `env:"USER,notEmpty"`
	Password string `env:"PASSWORD,notEmpty"`
	Database string `env:"DATABASE,notEmpty"`
	SSLMode  string `env:"SSL_MODE,notEmpty"`
}

func loadConfigFromEnv() (Config, error) {
	cfg, err := env.ParseAs[Config]()
	if err != nil {
		return Config{}, fmt.Errorf("failed to parse config: %w", err)
	}

	return cfg, nil
}
