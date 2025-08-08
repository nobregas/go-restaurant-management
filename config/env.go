package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	PUBLIC_HOST string
	PORT        string

	DB_ADDRESS  string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string

	JWT_SECRET string
	JWT_EXPIRE int64 // In seconds
}

var Envs = initConfig()

func initConfig() Config {
	godotenv.Load()

	return Config{
		PUBLIC_HOST: getEnv("PUBLIC_HOST", "localhost"),
		PORT:        getEnv("PORT", "8080"),
		DB_ADDRESS:  getEnv("DB_ADDRESS", "localhost"),
		DB_USER:     getEnv("DB_USER", "root"),
		DB_PASSWORD: getEnv("DB_PASSWORD", ""),
		DB_NAME:     getEnv("DB_NAME", "ecommerce"),
		JWT_SECRET:  getEnv("JWT_SECRET", "secret"),
		JWT_EXPIRE:  getEnvAsInt("JWT_EXPIRE", 1*60*60),
	}
}

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}
	return fallback
}
