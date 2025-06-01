package memorystorage

import (
	"errors"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage"
	"testing"
	"time"
)

func TestStorage(t *testing.T) {
	createCases := []struct {
		name        string
		event       storage.Event
		expectedErr error
	}{
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

	s := New()

	for _, tc := range createCases {
		t.Run(tc.name, func(t *testing.T) {

			err := s.AddEvent(tc.event)
			if !errors.Is(err, tc.expectedErr) {
				t.Errorf("expected error %v, got %v", tc.expectedErr, err)
			}
		})
	}

	removeCases := []struct {
		name        string
		eventID     string
		expectedErr error
	}{
		{
			name:    "delete second event",
			eventID: "test2",
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

	changeCases := []struct {
		name        string
		event       storage.Event
		expectedErr error
	}{
		{
			name: "change first event",
			event: storage.Event{
				ID:           "test1",
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
			endTime:     time.Date(2025, 6, 30, 10, 0, 0, 0, time.UTC),
		},
		{
			name:        "get events with not exist user",
			userId:      100,
			startTime:   time.Date(2025, 5, 16, 10, 0, 0, 0, time.UTC),
			endTime:     time.Date(2025, 5, 17, 10, 0, 0, 0, time.UTC),
			expectedErr: storage.ErrNotFound,
		},
		{
			name:        "get events with not exist period",
			userId:      1,
			startTime:   time.Date(2025, 4, 16, 10, 0, 0, 0, time.UTC),
			endTime:     time.Date(2025, 4, 17, 10, 0, 0, 0, time.UTC),
			expectedErr: storage.ErrNotFound,
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

// TODO проверка парралельной работы?!?
// TODO проверка что получаем данные только своего пользователя
// TODO проверка пограничных ситуаций по времени/дате
