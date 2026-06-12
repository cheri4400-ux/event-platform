package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"event-platform/internal/service"
)

type MockProducer struct{}

func (m *MockProducer) Send(
	ctx context.Context,
	key string,
	value any,
) error {
	return nil
}

func TestCreateEvent(t *testing.T) {

	mockRepo := &service.MockRepo{}
	mockProducer := &MockProducer{}

	eventService := service.NewEventService(
		mockRepo,
	)

	handler := NewEventHandler(
		eventService,
		mockProducer,
	)

	body := `{
		"user_id":"123",
		"event_type":"purchase",
		"page":"checkout",
		"amount":100
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.CreateEvent(
		rec,
		req,
	)

	if rec.Code != http.StatusOK {
		t.Fatalf(
			"expected 200, got %d",
			rec.Code,
		)
	}
}

func TestCreateEvent_InvalidJSON(t *testing.T) {

	mockRepo := &service.MockRepo{}
	mockProducer := &MockProducer{}

	eventService := service.NewEventService(
		mockRepo,
	)

	handler := NewEventHandler(
		eventService,
		mockProducer,
	)

	body := `{"user_id":`

	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.CreateEvent(
		rec,
		req,
	)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected 400, got %d",
			rec.Code,
		)
	}
}

func TestCreateEvent_ValidationError(t *testing.T) {

	mockRepo := &service.MockRepo{}
	mockProducer := &MockProducer{}

	eventService := service.NewEventService(
		mockRepo,
	)

	handler := NewEventHandler(
		eventService,
		mockProducer,
	)

	body := `{
		"event_type":"purchase"
	}`

	req := httptest.NewRequest(
		http.MethodPost,
		"/events",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.CreateEvent(
		rec,
		req,
	)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf(
			"expected 400, got %d",
			rec.Code,
		)
	}
}
