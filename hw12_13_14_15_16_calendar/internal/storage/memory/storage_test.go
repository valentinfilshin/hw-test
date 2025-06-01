package memorystorage

import (
	"errors"
	"fmt"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage"
	"sync"
	"testing"
	"time"
)

func TestMemoryStorage_AddEvent(t *testing.T) {
	s := setupTestStorage(t)

	cases := []struct {
		name        string
		event       storage.Event
		expectedErr error
	}{
		{
			name:        "",
			event:       storage.Event{},
			expectedErr: storage.ErrEmptyEventId,
		},
		{
			name: "add first event",
			event: storage.Event{
				ID:           "test1",
				Title:        "Простое событие",
				StartTime:    time.Date(2025, 6, 17, 10, 0, 0, 0, time.UTC),
				Duration:     2 * time.Hour,
				Description:  "",
				UserID:       1,
				NotifyBefore: 2 * time.Hour,
			},
		},
		{
			name: "add existing event",
			event: storage.Event{
				ID:           "test1",
				Title:        "Простое событие",
				StartTime:    time.Date(2025, 6, 17, 10, 0, 0, 0, time.UTC),
				Duration:     2 * time.Hour,
				Description:  "",
				UserID:       1,
				NotifyBefore: 2 * time.Hour,
			},
			expectedErr: storage.ErrEventExists,
		},
		{
			name: "add second event",
			event: storage.Event{
				ID:           "test2",
				Title:        "Простое событие 2",
				StartTime:    time.Date(2025, 5, 17, 10, 0, 0, 0, time.UTC),
				Duration:     2 * time.Hour,
				Description:  "",
				UserID:       1,
				NotifyBefore: 2 * time.Hour,
			},
		},
		{
			name: "add third event",
			event: storage.Event{
				ID:           "test3",
				Title:        "Простое событие 3",
				StartTime:    time.Date(2025, 5, 16, 10, 0, 0, 0, time.UTC),
				Duration:     2 * time.Hour,
				Description:  "",
				UserID:       2,
				NotifyBefore: 2 * time.Hour,
			},
		},
		{
			name: "add forth event",
			event: storage.Event{
				ID:           "test4",
				Title:        "Простое событие 4",
				StartTime:    time.Date(2025, 5, 16, 10, 0, 0, 0, time.UTC),
				Duration:     2 * time.Hour,
				Description:  "",
				UserID:       1,
				NotifyBefore: 2 * time.Hour,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			err := s.AddEvent(tc.event)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMemoryStorage_UpdateEvent(t *testing.T) {
	s := setupTestStorage(t)

	changeCases := []struct {
		name        string
		event       storage.Event
		expectedErr error
	}{
		{
			name: "change first event",
			event: storage.Event{
				ID:           "event1",
				Title:        "Простое измененное событие",
				StartTime:    time.Date(2025, 6, 17, 10, 0, 0, 0, time.UTC),
				Duration:     2 * time.Hour,
				Description:  "",
				UserID:       1,
				NotifyBefore: 2 * time.Hour,
			},
		},
		{
			name: "change not exist event",
			event: storage.Event{
				ID:           "not exist",
				Title:        "Простое событие",
				StartTime:    time.Date(2025, 6, 17, 10, 0, 0, 0, time.UTC),
				Duration:     2 * time.Hour,
				Description:  "",
				UserID:       1,
				NotifyBefore: 2 * time.Hour,
			},
			expectedErr: storage.ErrEventNotExist,
		},
	}

	for _, tc := range changeCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.ChangeEvent(tc.event)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMemoryStorage_DeleteEvent(t *testing.T) {
	s := setupTestStorage(t)

	removeCases := []struct {
		name        string
		eventID     string
		expectedErr error
	}{
		{
			name:    "delete second event",
			eventID: "event2",
		},
		{
			name:        "delete not exist event",
			eventID:     "not exist",
			expectedErr: storage.ErrEventNotExist,
		},
	}

	for _, tc := range removeCases {
		t.Run(tc.name, func(t *testing.T) {
			err := s.RemoveEvent(tc.eventID)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}
}

func TestMemoryStorage_ListEvents(t *testing.T) {
	s := setupTestStorage(t)

	getCases := []struct {
		name        string
		userId      int
		eventsCount int
		startTime   time.Time
		endTime     time.Time
		expectedErr error
	}{
		{
			name:        "get events",
			userId:      1,
			eventsCount: 2,
			startTime:   time.Date(2025, 5, 16, 10, 0, 0, 0, time.UTC),
			endTime:     time.Date(2025, 6, 17, 10, 0, 0, 0, time.UTC),
		},
		{
			name:        "get events with not exist user",
			userId:      100,
			startTime:   time.Date(2025, 5, 16, 10, 0, 0, 0, time.UTC),
			endTime:     time.Date(2025, 5, 17, 10, 0, 0, 0, time.UTC),
			expectedErr: storage.ErrEventsNotFound,
		},
		{
			name:        "get events with not exist period",
			userId:      1,
			startTime:   time.Date(2025, 4, 16, 10, 0, 0, 0, time.UTC),
			endTime:     time.Date(2025, 4, 17, 10, 0, 0, 0, time.UTC),
			expectedErr: storage.ErrEventsNotFound,
		},
	}

	for _, tc := range getCases {
		t.Run(tc.name, func(t *testing.T) {
			events, err := s.GetEvents(tc.userId, tc.startTime, tc.endTime)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
			if len(events) != tc.eventsCount {
				t.Errorf("expected %d events, got %d", tc.eventsCount, len(events))
			}
		})
	}
}

func TestMemoryStorage_ConcurrentAccess(t *testing.T) {
	s := setupTestStorage(t)

	const iterations = 100
	var wg sync.WaitGroup

	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()
			_, _ = s.GetEvents(i, time.Now(), time.Now().Add(24*time.Hour))
		}()
	}

	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		go func() {
			defer wg.Done()
			err := s.AddEvent(storage.Event{
				ID:        fmt.Sprintf("concurrent-event-%d", i),
				Title:     fmt.Sprintf("Конкурентное событие %d", i),
				StartTime: time.Now().Add(time.Duration(i) * time.Minute),
				UserID:    i,
			})
			if err != nil {
				return
			}
		}()
	}

	wg.Wait()
}

func setupTestStorage(t *testing.T) *Storage {
	t.Helper()
	s := New()

	events := []storage.Event{
		{
			ID:           "event1",
			Title:        "Событие 1",
			StartTime:    time.Date(2025, 6, 17, 10, 0, 0, 0, time.UTC),
			Duration:     2 * time.Hour,
			Description:  "",
			UserID:       1,
			NotifyBefore: 2 * time.Hour,
		},
		{
			ID:           "event2",
			Title:        "Событие 2",
			StartTime:    time.Date(2025, 6, 16, 10, 0, 0, 0, time.UTC),
			Duration:     2 * time.Hour,
			Description:  "",
			UserID:       1,
			NotifyBefore: 2 * time.Hour,
		},
		{
			ID:           "event3",
			Title:        "Событие 3",
			StartTime:    time.Date(2025, 5, 14, 10, 0, 0, 0, time.UTC),
			Duration:     2 * time.Hour,
			Description:  "",
			UserID:       1,
			NotifyBefore: 2 * time.Hour,
		},
		{
			ID:           "event4",
			Title:        "Событие 4",
			StartTime:    time.Date(2025, 6, 16, 10, 0, 0, 0, time.UTC),
			Duration:     2 * time.Hour,
			Description:  "",
			UserID:       2,
			NotifyBefore: 2 * time.Hour,
		},
		{
			ID:           "event5",
			Title:        "Событие 5",
			StartTime:    time.Date(2025, 6, 16, 10, 0, 0, 0, time.UTC),
			Duration:     2 * time.Hour,
			Description:  "",
			UserID:       2,
			NotifyBefore: 2 * time.Hour,
		},
	}

	for _, e := range events {
		if err := s.AddEvent(e); err != nil {
			t.Fatalf("не удалось добавить тестовое событие: %v", err)
		}
	}

	return s
}
