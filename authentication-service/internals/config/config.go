package config

import "time"

type Config struct {
	App appConfig
	DB  dbConfig
	JWT jwtConfig
}

type appConfig struct {
	Port string
}

type dbConfig struct {
	DNS string
}

type jwtConfig struct {
	Secret    string
	ExpiresAt time.Duration
}
