package config

import (
	"fmt"
	"log"
	"os"

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
	}
}

func getEnv(key, fallback string) string {
	value, isExist := os.LookupEnv(key)
	if isExist {
		return value
	}
	return fallback
}
