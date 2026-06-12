package service

import (
	"event-platform/internal/models"
	"log/slog"
)

type EventService struct {
	repo EventRepository
}

func NewEventService(repo EventRepository) *EventService {
	return &EventService{
		repo: repo,
	}
}

func (s *EventService) Process(event models.Event) error {
	slog.Info("processing event",
		"user_id", event.UserID,
		"event_type", event.EventType,
	)

	err := s.repo.Save(event)
	if err != nil {
		return err
	}

	slog.Info("event processed successfully")

	return nil
}

func (s *EventService) GetAll() ([]models.Event, error) {
	return s.repo.GetAll()
}
