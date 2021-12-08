package config

import "os"

type Config struct {
	HTTPPort string
	// Postgres
	Grpc string
}

func SetConfig() *Config {
	var config Config

	config.HTTPPort = os.Getenv("HTTPPort")
	if config.HTTPPort == "" {
		config.HTTPPort = ":8080"
	}

	config.Grpc = os.Getenv("GRPC")
	if config.Grpc == "" {
		config.Grpc = ":9000"
	}

	return &Config{
		HTTPPort: config.HTTPPort,
		Grpc:     config.Grpc,
	}
}
