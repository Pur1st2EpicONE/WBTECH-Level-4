// Package memory provides an in-memory implementation of the Storage interface.
package memory

import (
	"fmt"
	"sync"
	"time"

	"L4.3/internal/config"
	"L4.3/internal/models"
	"L4.3/pkg/logger"
	"github.com/google/uuid"
)

// Storage is an in-memory implementation of the repository.Storage interface.
// It stores events per user and per date, supports CRUD operations,
// and keeps auxiliary maps for fast lookup and user event counts.
//
// All methods are thread-safe using an internal RWMutex.
type Storage struct {
	db             map[int]map[string][]*models.Event // userID -> date string -> list of events
	eventsByID     map[string]*models.Event           // eventID -> event pointer
	userEventCount map[int]int                        // userID -> total number of events
	logger         logger.Logger                      // logger instance
	mu             sync.RWMutex                       // protects all maps
}

// NewStorage creates a new in-memory Storage instance.
// The initial map capacities are set based on the ExpectedUsers config value.
func NewStorage(config config.Storage, logger logger.Logger) *Storage {
	return &Storage{
		db:             make(map[int]map[string][]*models.Event, config.ExpectedUsers),
		eventsByID:     make(map[string]*models.Event, config.ExpectedUsers),
		userEventCount: make(map[int]int, config.ExpectedUsers),
		logger:         logger,
	}
}

// CreateEvent stores a new event in memory.
// Generates a unique UUID for the event and updates internal maps and counters.
func (s *Storage) CreateEvent(event *models.Event) (string, error) {

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, userExists := s.db[event.Meta.UserID]; !userExists {
		s.db[event.Meta.UserID] = make(map[string][]*models.Event)
		s.logger.Debug("repository — new user created", "UserID", event.Meta.UserID, "layer", "repository.memory")
	}

	eventDate := format(event.Meta.EventDate)
	event.Meta.EventID = uuid.New().String()

	s.db[event.Meta.UserID][eventDate] = append(s.db[event.Meta.UserID][eventDate], event)
	s.eventsByID[event.Meta.EventID] = event
	s.userEventCount[event.Meta.UserID]++

	return event.Meta.EventID, nil

}

// UpdateEvent updates an existing event's data or moves it to a new date.
// Thread-safe with write lock. Updates are logged.
func (s *Storage) UpdateEvent(new *models.Event) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	current := s.eventsByID[new.Meta.EventID]

	if current.Data != new.Data {
		updateData(&current.Data, &new.Data)
		s.logger.Debug("repository — event data updated", "UserID", new.Meta.UserID, "EventID", new.Meta.EventID, "layer", "repository.memory")
	}

	if !new.Meta.NewDate.IsZero() && !new.Meta.NewDate.Equal(current.Meta.EventDate) {

		newDate := format(new.Meta.NewDate)
		oldDate := format(current.Meta.EventDate)

		dayEvents := s.db[current.Meta.UserID][oldDate]

		for i, e := range dayEvents {

			if e.Meta.EventID == current.Meta.EventID {
				copy(dayEvents[i:], dayEvents[i+1:])
				dayEvents[len(dayEvents)-1] = nil
				dayEvents = dayEvents[:len(dayEvents)-1]
				break
			}

		}

		if len(dayEvents) == 0 {
			delete(s.db[current.Meta.UserID], oldDate)
		} else {
			s.db[current.Meta.UserID][oldDate] = dayEvents
		}

		current.Meta.EventDate = new.Meta.NewDate
		s.db[current.Meta.UserID][newDate] = append(s.db[current.Meta.UserID][newDate], current)

		s.logger.Debug("repository — event meta updated", "UserID", new.Meta.UserID, "EventID", new.Meta.EventID, "layer", "repository.memory")

	}

	return nil

}

// DeleteEvent removes an event from memory and updates counters.
// Uses write lock for thread safety.
func (s *Storage) DeleteEvent(meta *models.Meta) error {

	s.mu.Lock()
	defer s.mu.Unlock()

	current := s.eventsByID[meta.EventID]
	date := format(current.Meta.EventDate)

	userID := current.Meta.UserID
	dayEvents := s.db[userID][date]

	for i, e := range dayEvents {

		if e.Meta.EventID == meta.EventID {
			copy(dayEvents[i:], dayEvents[i+1:])
			dayEvents[len(dayEvents)-1] = nil
			dayEvents = dayEvents[:len(dayEvents)-1]
			break
		}

	}

	if len(dayEvents) == 0 {
		delete(s.db[userID], date)
	} else {
		s.db[userID][date] = dayEvents
	}

	s.userEventCount[userID]--
	delete(s.eventsByID, meta.EventID)

	return nil

}

// GetEventByID retrieves an event by its ID. Returns nil if not found.
// Thread-safe using read lock.
func (s *Storage) GetEventByID(eventID string) *models.Event {

	s.mu.RLock()
	defer s.mu.RUnlock()

	if event, eventFound := s.eventsByID[eventID]; eventFound {
		return event
	}

	return nil

}

// CountUserEvents returns the total number of events for a given user.
func (s *Storage) CountUserEvents(userID int) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.userEventCount[userID], nil
}

// GetEvents retrieves all events for a user filtered by period: day, week, or month.
// Returns empty slice if no events exist for the period.
func (s *Storage) GetEvents(meta *models.Meta, period models.Period) ([]models.Event, error) {

	s.mu.RLock()
	defer s.mu.RUnlock()

	if period != models.Day && period != models.Week && period != models.Month {
		return nil, fmt.Errorf("unknown period: %s", period)
	}

	allUserEvents, eventsFound := s.db[meta.UserID]
	if !eventsFound {
		return []models.Event{}, nil
	}

	switch period {

	case models.Day:
		return s.getEventsForDay(allUserEvents, meta)

	case models.Week:
		return s.getEventsForWeek(allUserEvents, meta)

	default:
		return s.getEventsForMonth(allUserEvents, meta)

	}

}

// getEventsForDay returns all events for a specific day for a user.
// It extracts events from the provided map, keyed by date strings (YYYY-MM-DD),
// and returns a slice of Event structs. If no events are found for the given day,
// an empty slice is returned. Thread safety must be ensured by the caller.
func (s *Storage) getEventsForDay(allUserEvents map[string][]*models.Event, meta *models.Meta) ([]models.Event, error) {

	dayEvents, eventsFound := allUserEvents[format(meta.EventDate)]
	if !eventsFound {
		return []models.Event{}, nil
	}

	res := make([]models.Event, len(dayEvents))
	for i, event := range dayEvents {
		res[i] = *event
	}

	return res, nil

}

// getEventsForWeek returns all events for the ISO week containing meta.EventDate for a user.
// It iterates over all user events, compares their ISO week and year to the target week,
// and collects matching events. Returns an empty slice if no events are found.
func (s *Storage) getEventsForWeek(allUserEvents map[string][]*models.Event, meta *models.Meta) ([]models.Event, error) {

	var res []models.Event
	targetYear, targetWeek := meta.EventDate.ISOWeek()

	for _, dayEvents := range allUserEvents {

		for _, event := range dayEvents {

			eventYear, entryWeek := event.Meta.EventDate.ISOWeek()
			if eventYear == targetYear && entryWeek == targetWeek {
				res = append(res, *event)
			}

		}

	}

	return res, nil

}

// getEventsForMonth returns all events for the month containing meta.EventDate for a user.
// It iterates over all user events and compares the year and month of each event
// to the target month. Matching events are collected and returned as a slice.
// Returns an empty slice if no events are found for the month.
func (s *Storage) getEventsForMonth(allUserEvents map[string][]*models.Event, meta *models.Meta) ([]models.Event, error) {

	var res []models.Event

	targetYear := meta.EventDate.Year()
	targetMonth := meta.EventDate.Month()

	for _, dayEvents := range allUserEvents {

		for _, event := range dayEvents {

			eventYear := event.Meta.EventDate.Year()
			eventMonth := event.Meta.EventDate.Month()

			if eventYear == targetYear && eventMonth == targetMonth {
				res = append(res, *event)
			}

		}

	}

	return res, nil

}

// Close clears all in-memory data and logs the shutdown.
func (s *Storage) Close() {

	s.mu.Lock()
	defer s.mu.Unlock()

	s.db = nil
	s.eventsByID = nil
	s.userEventCount = nil

	s.logger.LogInfo("in-memory storage — cleared and stopped", "layer", "repository.memory")

}

// updateData updates the event's textual data.
func updateData(current *models.Data, new *models.Data) {
	current.Text = new.Text
}

// format formats time.Time as a string in YYYY-MM-DD format.
func format(date time.Time) string {
	return date.Format("2006-01-02")
}
