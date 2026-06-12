package migrations

import (
	"database/sql"
	"log/slog"
	"os"
)

func Run(db *sql.DB) error {
	query, err := os.ReadFile(
		"internal/migrations/002_create_events.sql",
	)
	if err != nil {
		return err
	}

	_, err = db.Exec(string(query))
	if err != nil {
		return err
	}

	slog.Info(
		"migration completed",
		"file", "002_create_events.sql",
	)

	return nil
}
