package repository

import (
	"database/sql"
	"event-platform/internal/models"
)

type EventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
}

func (r *EventRepository) Save(event models.Event) error {
	_, err := r.db.Exec(`
		INSERT INTO events (
			user_id,
			event_type,
			page,
			amount
		)
		VALUES ($1, $2, $3, $4)
	`,
		event.UserID,
		event.EventType,
		event.Page,
		event.Amount,
	)

	return err
}
func (r *EventRepository) GetAll() ([]models.Event, error) {
	rows, err := r.db.Query(`
		SELECT
			id,
			user_id,
			event_type,
			page,
			amount,
			created_at
		FROM events
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var events []models.Event

	for rows.Next() {
		var event models.Event

		err := rows.Scan(
			&event.ID,
			&event.UserID,
			&event.EventType,
			&event.Page,
			&event.Amount,
			&event.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}
