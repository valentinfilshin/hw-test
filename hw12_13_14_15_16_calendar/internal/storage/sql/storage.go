package sqlstorage

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // Драйвер для PostgreSQL
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	db  *sql.DB
	dsn string
}

func New(dsn string) *Storage {
	return &Storage{dsn: dsn}
}

func (s *Storage) Connect() error {
	// Создаем пул подключений
	db, err := sql.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("failed to load driver: %w", err)
	}

	// Проверяем коннект к БД
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to db: %w", err)
	}

	// Настройки пула соединений
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

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

func (s *Storage) AddEvent(event storage.Event) error {
	query := `INSERT INTO events (title, start_time, end_time, description, user_id, notify_before) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := s.db.Exec(query, event.Title, event.StartTime, event.EndTime, event.Description, event.UserID,
		event.NotifyBefore)
	if err != nil {
		return fmt.Errorf("failed to add event: %w", err)
	}

	return nil
}

func (s *Storage) ChangeEvent(event storage.Event) error {
	query := `UPDATE events SET title = $1, start_time = $2, end_time = $3, description = $4, user_id = $5, 
	notify_before = $6 WHERE id = $7`

	_, err := s.db.Exec(query, event.Title, event.StartTime, event.EndTime, event.Description, event.UserID,
		event.NotifyBefore, event.ID)
	if err != nil {
		return fmt.Errorf("failed to change event: %w", err)
	}

	return nil
}

func (s *Storage) RemoveEvent(id string) error {
	query := `DELETE FROM events WHERE id = $1`

	_, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to remove event: %w", err)
	}

	return nil
}

func (s *Storage) GetEvents(userID int, from, to time.Time) ([]storage.Event, error) {
	query := `SELECT id, title, start_time, end_time, description, user_id, 
       EXTRACT(EPOCH FROM notify_before::interval)::int AS notify_before_seconds
	FROM events
	WHERE user_id = $1 AND start_time >= $2 and start_time <= $3`

	rows, err := s.db.Query(query, userID, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get events: %w", err)
	}
	defer rows.Close()

	var result []storage.Event
	for rows.Next() {
		var event storage.Event
		var userID int64
		var id, title, description, notifyBefore string
		var startTime, endTime time.Time

		if err := rows.Scan(&id, &title, &startTime, &endTime, &description, &userID, &notifyBefore); err != nil {
			return nil, fmt.Errorf("failed to get events: %w", err)
		}

		event.ID = id
		event.Title = title
		event.StartTime = startTime
		event.EndTime = endTime
		event.Description = description
		event.UserID = int(userID)
		event.NotifyBefore, err = time.ParseDuration(notifyBefore + "s")
		if err != nil {
			return nil, fmt.Errorf("failed to parse notify_before: %w", err)
		}

		result = append(result, event)
	}

	return result, nil
}
