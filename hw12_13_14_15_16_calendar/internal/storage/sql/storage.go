package sqlstorage

import (
	"context"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Драйвер для Postgres
	"github.com/jmoiron/sqlx"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db  *sqlx.DB
	dsn string
}

func New(dsn string) *Storage {
	return &Storage{dsn: dsn}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	s.db = db

	return nil
}

func (s *Storage) Close() error {
	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}

// TODO что будет если одновременно два запроса попробуют забронировать свободное пересекающееся время

func (s *Storage) AddEvent(event storage.Event) error {
	// Проверка на существование других с пересекающимся временем перед вставкой + transaction?
	_, err := s.db.Exec(
		"INSERT INTO events (id, title, start_time, end_time, description, user_id, notify_before) "+
			"VALUES ($1, $2, $3, $4, $5, $6, $7)",
		event.ID, event.Title, event.StartTime, event.EndTime, event.Description, event.UserID, event.NotifyBefore)
	if err != nil {
		return fmt.Errorf("failed to add event: %w", err)
	}

	return nil
}

func (s *Storage) ChangeEvent(event storage.Event) error {
	// Проверка на существование других с пересекающимся временем перед обновлением + transaction?
	_, err := s.db.Exec(
		"UPDATE events SET title = $1, start_time = $2, end_time = $3, description = $4, user_id = $5, notify_before = $6 "+
			"WHERE id = $7",
		event.Title, event.StartTime, event.EndTime, event.Description, event.UserID, event.NotifyBefore, event.ID)
	if err != nil {
		return fmt.Errorf("failed to change event: %w", err)
	}

	return nil
}

func (s *Storage) RemoveEvent(id string) error {
	_, err := s.db.Exec("DELETE FROM events WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to remove event: %w", err)
	}

	return nil
}

func (s *Storage) GetEvents(userID int, from, to time.Time) ([]storage.Event, error) {
	// TODO маппинг в структуры
	result, err := s.db.Exec("SELECT * FROM events WHERE user_id = $1 "+
		"AND (start_time >= $2 OR end_time <= $3)", userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}

	fmt.Println(result)

	return nil, nil
}
