package storage

import (
	"errors"
	"time"
)

var (
	ErrEventNotExist = errors.New("event not exists")
	ErrEventExists   = errors.New("event already exists")
	ErrNotFound      = errors.New("events not found")
)

type Event struct {
	ID           string
	Title        string
	StartTime    time.Time
	Duration     time.Duration
	Description  string
	UserID       int
	NotifyBefore time.Duration
}

func (e *Event) IsBetween(from, to time.Time) bool {
	return (e.StartTime.After(from) || e.StartTime.Equal(from)) &&
		(e.StartTime.Before(to) || e.StartTime.Equal(to))
}
