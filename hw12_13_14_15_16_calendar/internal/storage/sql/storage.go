package sqlstorage

import (
	"context"
	"fmt"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage"

	_ "github.com/jackc/pgx/v5"
	"github.com/jmoiron/sqlx"
	"time"
)

type Storage struct {
	db  *sqlx.DB
	dsn string
}

func New(dsn string) *Storage {
	return &Storage{dsn: dsn}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Connect("postgres", dsn)

	if err != nil {
		return fmt.Errorf("could not connect to database: %w", err)
	}

	s.db = db

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	err := s.db.Close()

	if err != nil {
		return fmt.Errorf("failed to close database: %w", err)
	}

	return nil
}

func (s *Storage) AddEvent(event storage.Event) error {
	return nil
}

func (s *Storage) ChangeEvent(event storage.Event) error {
	return nil
}

func (s *Storage) RemoveEvent(id string) error {
	return nil
}

func (s *Storage) GetEvents(userID int, from, to time.Time) ([]storage.Event, error) {
	return nil, nil
}
