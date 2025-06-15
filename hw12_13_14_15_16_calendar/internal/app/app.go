package app

import (
	"context"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage"
	"time"
)

type App struct { // TODO
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Storage interface {
	AddEvent(event storage.Event) error
	ChangeEvent(event storage.Event) error
	RemoveEvent(id string) error
	GetEvents(userID int, from, to time.Time) ([]storage.Event, error)
}

func New(_ Logger, _ Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(_ context.Context, _, _ string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
