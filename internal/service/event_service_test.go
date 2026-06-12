package service

import (
	"event-platform/internal/models"
	"testing"
)

func TestProcessEvent(t *testing.T) {

	repo := &MockRepo{}

	service := NewEventService(repo)

	event := models.Event{
		UserID:    "123",
		EventType: "purchase",
		Page:      "checkout",
		Amount:    100,
	}

	err := service.Process(event)

	if err != nil {
		t.Fatalf(
			"expected nil error, got %v",
			err,
		)
	}
}
