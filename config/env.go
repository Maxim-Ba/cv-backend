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
	LogLevel         string `env:"LOG_LEVEL" envDefault:"error"`
	AppEnv           string `env:"APP_ENV" envDefault:"development"`
	AllowedOrigin    string `env:"ALLOWED_ORIGIN" envDefault:"http://localhost:4200"`
	AdminUser        string `env:"ADMIN_USER" envDefault:"admin"`
	AdminPassword    string `env:"ADMIN_PASSWORD" envDefault:"admin"`
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
