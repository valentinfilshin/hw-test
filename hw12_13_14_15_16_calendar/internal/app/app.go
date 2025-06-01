package app

import (
	"context"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage"
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
}

func New(logger Logger, storage Storage) *App {
	return &App{}
}

func (a *App) CreateEvent(ctx context.Context, id, title string) error {
	// TODO
	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
