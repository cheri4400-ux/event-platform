package service

import "event-platform/internal/models"

type EventRepository interface {
	Save(event models.Event) error
	GetAll() ([]models.Event, error)
}
