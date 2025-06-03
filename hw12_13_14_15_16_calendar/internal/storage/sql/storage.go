package sqlstorage

import (
	"context"
	"database/sql"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage"
	"time"

	_ "github.com/jackc/pgx/v5"
)

type Storage struct {
	db *sql.DB
}

func New() *Storage {

	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) AddEvent(event storage.Event) error {

}

func (s *Storage) ChangeEvent(event storage.Event) error {

}

func (s *Storage) RemoveEvent(id string) error {

}

func (s *Storage) GetEvents(userID int, from, to time.Time) ([]storage.Event, error) {

}
