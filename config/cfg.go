package config

type Config struct {
	Secret           string
	ServerAddr       string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	MigrationPath    string
	LogLevel         string
	AppEnv           string
}

var cfg Config

func init() {
	envs, err := parseEnv()
	if err != nil {
		panic(err)
	} else {
		cfg = Config{
			ServerAddr:       envs.ServerAddr,
			PostgresHost:     envs.PostgresHost,
			PostgresPort:     envs.PostgresPort,
			PostgresUser:     envs.PostgresUser,
			PostgresPassword: envs.PostgresPassword,
			PostgresDB:       envs.PostgresDB,
			MigrationPath:    envs.MigrationPath,
			LogLevel:         envs.LogLevel,
			AppEnv:           envs.AppEnv,
		}
	}
}
func GetConfig() *Config {
	return &cfg
}
