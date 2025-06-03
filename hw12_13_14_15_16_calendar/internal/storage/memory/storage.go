package memorystorage

import (
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage"
	"sync"
	"time"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

// TODO добавить добавление не заблокированного слота

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) AddEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if event.ID == "" {
		return storage.ErrEmptyEventId
	}

	if _, ok := s.events[event.ID]; ok {
		return storage.ErrEventExists
	}
	s.events[event.ID] = event

	return nil
}

func (s *Storage) GetEvents(userID int, from, to time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storage.Event, 0)

	for _, event := range s.events {
		if event.UserID == userID && event.IsBetween(from, to) {
			result = append(result, event)
		}
	}

	if len(result) == 0 {
		return nil, storage.ErrEventsNotFound
	}

	return result, nil
}

func (s *Storage) ChangeEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[event.ID]; !ok {
		return storage.ErrEventNotExist
	}

	s.events[event.ID] = event

	return nil
}

func (s *Storage) RemoveEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.events[id]; !ok {
		return storage.ErrEventNotExist
	}

	delete(s.events, id)

	return nil
}
