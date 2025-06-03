package storage

import (
	"errors"
	"time"
)

var (
	ErrEventNotExist  = errors.New("event not exists")
	ErrEventExists    = errors.New("event already exists")
	ErrEventsNotFound = errors.New("events not found")
	ErrEmptyEventID   = errors.New("event id can't be empty")
	ErrDateBusy       = errors.New("date is busy")
	// TODO ошибки валидации, а не добавления в БД?
	ErrStartAfterEnd  = errors.New("start time can't be after end time")
	ErrEndBeforeStart = errors.New("end time can't be before start time")
)

type Event struct {
	ID           string
	Title        string
	StartTime    time.Time
	EndTime      time.Time
	Description  string
	UserID       int
	NotifyBefore time.Duration
}

func (e *Event) IntersectsWith(from, to time.Time) bool {
	return e.StartTime.Before(to) && e.EndTime.After(from)
}
