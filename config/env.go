package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type Envs struct {
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDB       string `env:"POSTGRES_DB"`
	ServerAddr       string `env:"SERVER_ADDRESS"`
	MigrationPath    string `env:"MIGRATION_PATH"`
	LogLevel         string `env:"LOG_LEVEL" default:"error"`
	AppEnv           string `env:"APP_ENV" default:"development"`
}

func parseEnv() (*Envs, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("No .env file found, using system environment variables")
	}
	e := Envs{}
	if err := env.Parse(&e); err != nil {
		return nil, err
	}
	return &e, nil
}
