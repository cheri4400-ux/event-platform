package models

import (
	"errors"
	"time"
)

type Event struct {
	ID        int       `json:"id"`
	UserID    string    `json:"user_id"`
	EventType string    `json:"event_type"`
	Page      string    `json:"page"`
	Amount    float64   `json:"amount,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func (e *Event) Validate() error {
	if e.UserID == "" {
		return errors.New("user_id is required")
	}

	if e.EventType == "" {
		return errors.New("event_type is required")
	}

	return nil
}
