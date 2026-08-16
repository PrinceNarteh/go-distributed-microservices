package config

type Config struct {
	App appConfig
	DB  dbConfig
}

type appConfig struct {
	Port string
}

type dbConfig struct {
	DNS string
}
