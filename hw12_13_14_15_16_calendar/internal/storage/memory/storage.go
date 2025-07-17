package memorystorage

import (
	"errors"
	"sync"
	"time"

	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	events map[string]storage.Event
	mu     sync.RWMutex
}

func New() *Storage {
	return &Storage{
		events: make(map[string]storage.Event),
	}
}

func (s *Storage) AddEvent(event storage.Event) error {
	if event.ID == "" {
		return storage.ErrEmptyEventID
	}

	events, err := s.GetEvents(event.UserID, event.StartTime, event.EndTime)
	if err != nil && !errors.Is(err, storage.ErrEventsNotFound) {
		return err
	}

	if len(events) > 0 {
		return storage.ErrDateBusy
	}

	s.mu.Lock()
	defer s.mu.Unlock()

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
		if event.UserID == userID && event.IntersectsWith(from, to) {
			result = append(result, event)
		}
	}

	if len(result) == 0 {
		return nil, storage.ErrEventsNotFound
	}

	return result, nil
}

func (s *Storage) ChangeEvent(event storage.Event) error {
	events, err := s.GetEvents(event.UserID, event.StartTime, event.EndTime)
	if err != nil && !errors.Is(err, storage.ErrEventsNotFound) {
		return err
	}

	if len(events) > 1 || len(events) == 1 && events[0].ID != event.ID {
		return storage.ErrDateBusy
	}

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
