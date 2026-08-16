package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var Env = initConfig()

func initConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file")
	}

	return &Config{
		App: appConfig{
			Port: fmt.Sprintf(":%s", getEnv("APP_PORT", "8080")),
		},
		DB: dbConfig{
			DNS: getEnv("DNS", ""),
		},
		JWT: jwtConfig{
			Secret:    getEnv("JWT_SECRET", "top-secret-to-generate-token"),
			ExpiresAt: getEnvAsDuration("JWT_SECRET", time.Hour),
		},
	}
}

func getEnv(key, fallback string) string {
	value, isExist := os.LookupEnv(key)
	if isExist {
		return value
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	durationAsStr := getEnv(key, "")
	if durationAsStr == "" {
		return fallback
	}

	duration, err := time.ParseDuration(durationAsStr)
	if err != nil {
		return fallback
	}

	return duration
}
