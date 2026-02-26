package config

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	DBPath             string
	JWTSecret          string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

var App Config

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables directly")
	}

	App = Config{
		Port:               getEnv("PORT", "8080"),
		DBPath:             getEnv("DB_PATH", "api.db"),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		AccessTokenExpiry:  parseDuration("ACCESS_TOKEN_EXPIRY", "15m"),
		RefreshTokenExpiry: parseDuration("REFRESH_TOKEN_EXPIRY", "168h"),
	}

	if App.JWTSecret == "" {
		log.Fatal("JWT_SECRET environment variable is required")
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// helper to parse duration strings like "15m" or "168h"
func parseDuration(key, defaultValue string) time.Duration {
	value := getEnv(key, defaultValue)
	duration, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Invalid duration for %s, using default %s", key, defaultValue)
		duration, _ = time.ParseDuration(defaultValue)
	}
	return duration
}
