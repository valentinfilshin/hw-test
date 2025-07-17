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
	ErrDateBusy       = errors.New("event date is busy")
	ErrStartAfterEnd  = errors.New("start time can't be after end time")
	ErrEndBeforeStart = errors.New("end time can't be before start time")
)

type Event struct {
	ID           string        `db:"id"`
	Title        string        `db:"title"`
	StartTime    time.Time     `db:"start_time"`
	EndTime      time.Time     `db:"end_time"`
	Description  string        `db:"description"`
	UserID       int           `db:"user_id"`
	NotifyBefore time.Duration `db:"notify_before"`
}

func (e *Event) IntersectsWith(from, to time.Time) bool {
	return e.StartTime.Before(to) && e.EndTime.After(from)
}
