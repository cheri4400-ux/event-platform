// @title Event Platform API
// @version 1.0
// @description Event Tracking Platform
// @host localhost:8080
// @BasePath /

package main

import (
	"context"
	_ "event-platform/docs"
	"event-platform/internal/config"
	"event-platform/internal/handler"
	"event-platform/internal/kafka"
	"event-platform/internal/postgres"
	"event-platform/internal/repository"
	"event-platform/internal/service"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// 1. Настраиваем structured logger
	logger := slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}),
	)

	slog.SetDefault(logger)

	cfg := config.Load()

	producer := kafka.NewProducer(
		cfg.KafkaBrokers,
		cfg.KafkaTopic,
	)

	db, err := postgres.NewPostgresDB(cfg)
	if err != nil {
		slog.Error("failed to connect postgres",
			"error", err,
		)
		os.Exit(1)
	}

	defer db.Close()

	eventRepo := repository.NewEventRepository(db)

	eventService := service.NewEventService(eventRepo)

	eventHandler := handler.NewEventHandler(
		eventService,
		producer,
	)

	// 2. HTTP mux
	mux := http.NewServeMux()

	mux.Handle(
		"/swagger/",
		httpSwagger.WrapHandler,
	)

	mux.Handle(
		"/metrics",
		promhttp.Handler(),
	)

	// 3. Health check endpoint
	mux.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			eventHandler.GetEvents(w, r)

		case http.MethodPost:
			eventHandler.CreateEvent(w, r)

		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	slog.Info("API starting", "port", 8080)

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			slog.Error("server failed",
				"error", err,
			)
			os.Exit(1)
		}
	}()

	sig := <-signalChan

	slog.Info("shutdown signal received",
		"signal", sig.String(),
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		slog.Error("shutdown failed",
			"error", err,
		)
	}

	err = producer.Close()
	if err != nil {
		slog.Error("failed to close producer",
			"error", err,
		)
	}

	slog.Info("api stopped gracefully")

}
