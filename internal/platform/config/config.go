package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	HTTPPort    string
	DatabaseURL string
	RedisAddr   string
	NATSURL     string
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system environment")
	}

	return Config{
		AppEnv:      getEnv("APP_ENV", "development"),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6380"),
		NATSURL:     getEnv("NATS_URL", "nats://localhost:4223"),
	}
}

func getEnv(key string, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	return val
}
