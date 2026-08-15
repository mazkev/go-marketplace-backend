package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort          string
	DatabaseDriver   string
	DatabaseDSN      string
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	port := getEnv("PORT", "8080")
	dbDriver := getEnv("DB_DRIVER", "sqlite")
	dbDSN := getEnv("DB_DSN", "marketplace.db")
	jwtSecret := getEnv("JWT_SECRET", "super-secret-jwt-key-change-in-production-12345")

	return &Config{
		AppPort:          port,
		DatabaseDriver:   dbDriver,
		DatabaseDSN:      dbDSN,
		JWTSecret:        jwtSecret,
		JWTAccessExpiry:  24 * time.Hour,
		JWTRefreshExpiry: 7 * 24 * time.Hour,
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
