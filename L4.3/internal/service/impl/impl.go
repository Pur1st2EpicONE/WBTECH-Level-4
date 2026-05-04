// Package impl provides the concrete implementation of the Service interface.
// It contains the business logic for creating, updating, deleting, and retrieving events.
package impl

import (
	"fmt"
	"sort"

	"L4.3/internal/config"
	"L4.3/internal/errs"
	"L4.3/internal/models"
	"L4.3/internal/repository"
	"L4.3/pkg/logger"
)

// Service is the implementation of the event management service.
// It interacts with the repository layer to perform CRUD operations on events
// and enforces business rules such as user limits and validation checks.
type Service struct {
	Storage          repository.Storage // underlying storage for events
	logger           logger.Logger      // logger for service-level logging
	maxEventsPerUser int                // maximum number of events allowed per user
}

// NewService creates a new Service instance with the provided configuration, storage, and logger.
// The maxEventsPerUser field is set from the configuration and enforces limits on event creation.
func NewService(config config.Service, storage repository.Storage, logger logger.Logger) *Service {
	return &Service{Storage: storage, logger: logger, maxEventsPerUser: config.MaxEventsPerUser}
}

// CreateEvent validates and creates a new event for a user.
// Returns the generated event ID on success, or an error if creation fails.
// It enforces per-user limits and validates the event data before calling the repository.
func (s *Service) CreateEvent(event *models.Event) (string, error) {

	if err := validateCreate(event); err != nil {
		return "", err
	}

	count, err := s.Storage.CountUserEvents(event.Meta.UserID)
	if err != nil {
		return "", err
	}

	s.logger.Debug(fmt.Sprintf("service — user %d has %d remaining event slots", event.Meta.UserID, s.maxEventsPerUser-count), "UserID", event.Meta.UserID, "layer", "service.impl")

	if count >= s.maxEventsPerUser {
		return "", errs.ErrMaxEvents
	}

	return s.Storage.CreateEvent(event)

}

// UpdateEvent validates and updates an existing event.
// Returns an error if validation fails, the event does not exist, or the update cannot be applied.
func (s *Service) UpdateEvent(event *models.Event) error {

	if err := validateIDs(event.Meta.UserID, event.Meta.EventID); err != nil {
		return err
	}

	if err := validateUpdate(event, s.Storage.GetEventByID(event.Meta.EventID)); err != nil {
		return err
	}

	return s.Storage.UpdateEvent(event)

}

// DeleteEvent validates and deletes an event identified by the provided metadata.
// Returns an error if validation fails or the event cannot be deleted.
func (s *Service) DeleteEvent(meta *models.Meta) error {

	if err := validateIDs(meta.UserID, meta.EventID); err != nil {
		return err
	}

	if err := validateDelete(meta, s.Storage.GetEventByID(meta.EventID)); err != nil {
		return err
	}

	return s.Storage.DeleteEvent(meta)

}

// GetEvents retrieves all events for a user within the specified period (day, week, month).
// Events are returned in descending order by date. Returns an error if validation fails
// or if the repository fails to fetch events.
func (s *Service) GetEvents(meta *models.Meta, period models.Period) ([]models.Event, error) {

	if err := validateGet(meta); err != nil {
		return nil, err
	}

	events, err := s.Storage.GetEvents(meta, period)
	if err != nil {
		return nil, err
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].Meta.EventDate.After(events[j].Meta.EventDate)
	})

	return events, nil

}
