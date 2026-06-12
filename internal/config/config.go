package config

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	KafkaBrokers []string
	KafkaTopic   string

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
}

func Load() *Config {

	err := godotenv.Load()
	if err != nil {
		slog.Warn(".env file not found")
	}

	cfg := &Config{
		KafkaBrokers: []string{
			getEnv("KAFKA_BROKERS", "localhost:9092"),
		},
		KafkaTopic: getEnv("KAFKA_TOPIC", "events"),

		PostgresHost:     getEnv("POSTGRES_HOST", "localhost"),
		PostgresPort:     getEnv("POSTGRES_PORT", "5432"),
		PostgresUser:     getEnv("POSTGRES_USER", "admin"),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", "admin"),
		PostgresDB:       getEnv("POSTGRES_DB", "events_db"),
	}

	slog.Info("config loaded",
		"kafka_topic", cfg.KafkaTopic,
		"postgres_db", cfg.PostgresDB,
	)

	return cfg
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}
