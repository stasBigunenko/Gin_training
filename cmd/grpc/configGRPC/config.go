package configGRPC

import "os"

type Config struct {
	TcpPort string
	// Postgres
	PostgresHost string
	PostgresPort string
	PostgresUser string
	PostgresPsw  string
	PostgresDB   string
	PostgresSSL  string
}

func SetConfig() *Config {
	var config Config

	config.TcpPort = os.Getenv("TcpPort")
	if config.TcpPort == "" {
		config.TcpPort = "127.0.0.1:9000"
	}

	config.PostgresHost = os.Getenv("POSTGRES_HOST")
	if config.PostgresHost == "" {
		config.PostgresHost = "localhost"
	}

	config.PostgresPort = os.Getenv("POSTGRES_PORT")
	if config.PostgresPort == "" {
		config.PostgresPort = "5432"
	}

	config.PostgresUser = os.Getenv("POSTGRES_USER")
	if config.PostgresUser == "" {
		config.PostgresUser = "postgres"
	}

	config.PostgresPsw = os.Getenv("POSTGRES_PASSWORD")
	if config.PostgresPsw == "" {
		config.PostgresPsw = "postgres"
	}

	config.PostgresDB = os.Getenv("POSTGRES_DATABASE")
	if config.PostgresDB == "" {
		config.PostgresDB = "postgres"
	}

	config.PostgresSSL = os.Getenv("POSTGRES_SSL")
	if config.PostgresSSL == "" {
		config.PostgresSSL = "disable"
	}

	return &Config{
		TcpPort:      config.TcpPort,
		PostgresHost: config.PostgresHost,
		PostgresPort: config.PostgresPort,
		PostgresUser: config.PostgresUser,
		PostgresPsw:  config.PostgresPsw,
		PostgresDB:   config.PostgresDB,
		PostgresSSL:  config.PostgresSSL,
	}
}
