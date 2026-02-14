package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBPath    string
	JWTSecret string
}

var App Config

func Load() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables directly")
	}

	App = Config{
		Port:      getEnv("PORT", "8080"),
		DBPath:    getEnv("DB_PATH", "api.db"),
		JWTSecret: getEnv("JWT_SECRET", ""),
	}

	// jwt secret is critical
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
