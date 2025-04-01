package config

import (
	"fmt"
	"os"
)

type Config struct {
	PostgresURL  string
	RabbitMQURL  string
}

func Load() (*Config, error) {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		return nil, fmt.Errorf("POSTGRES_URL not set")
	}

	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		return nil, fmt.Errorf("RABBITMQ_URL not set")
	}

	return &Config{
		PostgresURL:  postgresURL,
		RabbitMQURL:  rabbitMQURL,
	}, nil
}