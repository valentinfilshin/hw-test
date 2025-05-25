package memorystorage

import (
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage"
	"sync"
)

type Storage struct {
	events map[int]storage.Event
	mu     sync.RWMutex //nolint:unused
	nextID int
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) AddEvent(event storage.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[s.nextID] = event
	s.nextID++

	return nil
}

func (s *Storage) GetEvents(userID int) ([]storage.Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]storage.Event, 0)

	for _, event := range s.events {
		if event.UserID == userID {
			result = append(result, event)
		}
	}

	return result, nil
}

func (s *Storage) ChangeEvent(event storage.Event) error {
	return nil
}

func (s *Storage) RemoveEvent(id string) error {
	return nil
}
