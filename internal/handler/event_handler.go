package handler

import (
	"context"
	"encoding/json"
	"event-platform/internal/models"
	"event-platform/internal/service"
	"log/slog"
	"net/http"
)

type EventHandler struct {
	service  *service.EventService
	producer Producer
}

func NewEventHandler(
	service *service.EventService,
	producer Producer,
) *EventHandler {
	return &EventHandler{
		service:  service,
		producer: producer,
	}
}

// GetEvents godoc
//
// @Summary Get all events
// @Description Returns all saved events
// @Tags events
// @Produce json
// @Success 200 {array} models.Event
// @Router /events [get]
func (h *EventHandler) GetEvents(
	w http.ResponseWriter,
	r *http.Request,
) {
	events, err := h.service.GetAll()
	if err != nil {
		slog.Error(
			"failed to get events",
			"error", err,
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	err = json.NewEncoder(w).Encode(events)
	if err != nil {
		slog.Error(
			"failed to encode events",
			"error", err,
		)
	}
}

// CreateEvent godoc
//
// @Summary Create event
// @Description Send event to Kafka
// @Tags events
// @Accept json
// @Produce plain
// @Param event body models.Event true "Event"
// @Success 200 {string} string "event sent to kafka"
// @Failure 400 {string} string "invalid request"
// @Router /events [post]
func (h *EventHandler) CreateEvent(
	w http.ResponseWriter,
	r *http.Request,
) {
	var event models.Event

	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid json"))
		return
	}

	err = event.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err = h.producer.Send(
		context.Background(),
		event.UserID,
		event,
	)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("event sent to kafka"))
}
