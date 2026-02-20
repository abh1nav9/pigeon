package config

import (
	"log"
	"os"
)

type Config struct {
	AppEnv      string
	HTTPPort    string
	PostgresURL string
	RedisURL    string
	JWTSecret   string
}

func Load() *Config {
	cfg := &Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
		PostgresURL: getEnv("POSTGRES_URL", ""),
		RedisURL:    getEnv("REDIS_URL", ""),
		JWTSecret:   getEnv("JWT_SECRET", ""),
	}

	if cfg.PostgresURL == "" {
		log.Fatal("POSTGRES_URL is required")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}