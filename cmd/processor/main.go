package main

import (
	"context"
	"encoding/json"
	"event-platform/internal/config"
	"event-platform/internal/migrations"
	"event-platform/internal/models"
	"event-platform/internal/postgres"
	"event-platform/internal/repository"
	"event-platform/internal/service"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/segmentio/kafka-go"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		slog.Error("failed to load .env",
			"error", err,
		)
		os.Exit(1)
	}

	// logger
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	slog.SetDefault(logger)

	cfg := config.Load()

	slog.Info(
		"config loaded",
		"kafka_topic", cfg.KafkaTopic,
		"postgres_db", cfg.PostgresDB,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// postgres
	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		slog.Error(
			"failed to connect postgres",
			"error", err,
		)
		os.Exit(1)
	}

	err = migrations.Run(db)
	if err != nil {
		slog.Error(
			"migration failed",
			"error", err,
		)
		os.Exit(1)
	}

	defer db.Close()

	eventRepo := repository.NewEventRepository(db)

	eventService := service.NewEventService(eventRepo)

	// kafka reader
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "events",
		GroupID: "event-processors",
	})

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	go func() {
		sig := <-signalChan

		slog.Info("shutdown signal received",
			"signal", sig.String(),
		)

		cancel()

		err := reader.Close()
		if err != nil {
			slog.Error("failed to close kafka reader",
				"error", err,
			)
		}

		err = db.Close()
		if err != nil {
			slog.Error("failed to close postgres",
				"error", err,
			)
		}

		slog.Info("processor stopped gracefully")

		os.Exit(0)
	}()

	defer reader.Close()

	slog.Info("processor started")

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			slog.Error("failed to read message",
				"error", err,
			)
			continue
		}

		slog.Info("message received from kafka",
			"key", string(msg.Key),
			"value", string(msg.Value),
		)

		var event models.Event

		err = json.Unmarshal(msg.Value, &event)
		if err != nil {
			slog.Error("failed to unmarshal event",
				"error", err,
			)
			continue
		}

		slog.Info("event processed",
			"user_id", event.UserID,
			"event_type", event.EventType,
			"page", event.Page,
			"amount", event.Amount,
		)

		err = eventService.Process(event)
		if err != nil {
			slog.Error("failed to save event",
				"error", err,
			)
			continue
		}

		slog.Info("event saved to postgres")
	}
}
