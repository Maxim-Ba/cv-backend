package config

import (
	"github.com/caarlos0/env/v11"
)

type Envs struct {
	PostgresHost     string `env:"POSTGRES_HOST"`
	PostgresPort     string `env:"POSTGRES_PORT"`
	PostgresUser     string `env:"POSTGRES_USER"`
	PostgresPassword string `env:"POSTGRES_PASSWORD"`
	PostgresDB       string `env:"POSTGRES_DB"`
	ServerAddr       string `env:"SERVER_ADDRESS"`
	MigrationPath    string `env:"MIGRATION_PATH"`
}

func parseEnv() (*Envs, error) {
	e := Envs{}
	if err := env.Parse(&e); err != nil {
		return nil, err
	}
	return &e, nil
}
