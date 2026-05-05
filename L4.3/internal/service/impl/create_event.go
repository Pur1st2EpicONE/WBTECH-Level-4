package impl

import (
	"fmt"
	"time"

	"L4.3/internal/errs"
	"L4.3/internal/models"
)

// CreateEvent validates and creates a new event in storage.
//
// Workflow:
//  1. validates input event
//  2. checks per-user event limit
//  3. persists event in memory storage
//  4. schedules reminder if configured
func (s *Service) CreateEvent(event *models.Event) (string, error) {

	if err := validateCreate(event); err != nil {
		return "", err
	}

	count, err := s.Storage.Memory.CountUserEvents(event.Meta.UserID)
	if err != nil {
		return "", err
	}

	s.logger.Debug(fmt.Sprintf("service — user %d has %d remaining event slots", event.Meta.UserID, s.maxEventsPerUser-count), "UserID", event.Meta.UserID, "layer", "service.impl")

	if count >= s.maxEventsPerUser {
		return "", errs.ErrMaxEvents
	}

	id, err := s.Storage.Memory.CreateEvent(event)
	if err != nil {
		return "", err
	}

	if event.Data.Reminder > 0 {
		remindAt := event.Meta.EventDate.Add(-event.Data.Reminder)
		if remindAt.After(time.Now()) {
			s.reminderCh <- models.Reminder{
				EventID:  id,
				UserID:   event.Meta.UserID,
				RemindAt: remindAt,
				Text:     event.Data.Text,
			}
		}
	}

	return id, nil

}
