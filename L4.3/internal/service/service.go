// Package service defines the business logic layer for managing user events.
// It provides an interface for creating, updating, deleting, and retrieving events.
package service

import (
	"L4.3/internal/config"
	"L4.3/internal/models"
	"L4.3/internal/repository"
	"L4.3/internal/service/impl"
	"L4.3/pkg/logger"
)

// Service defines the interface for the event management service.
// It encapsulates all business logic related to events, including creation,
// updates, deletion, and retrieval for different time periods.
type Service interface {
	// CreateEvent creates a new event for a user and returns the generated event ID.
	// Returns an error if the creation fails or validation rules are violated.
	CreateEvent(event *models.Event) (string, error)

	// UpdateEvent updates an existing event's data or date.
	// Returns an error if the event does not exist or no changes are detected.
	UpdateEvent(event *models.Event) error

	// DeleteEvent removes an event identified by the provided metadata.
	// Returns an error if the event does not exist or cannot be deleted.
	DeleteEvent(meta *models.Meta) error

	// GetEvents retrieves all events for a user within a specified period (day, week, month).
	// Returns a slice of events and an error if retrieval fails.
	GetEvents(meta *models.Meta, period models.Period) ([]models.Event, error)
}

// NewService creates a new Service implementation using the provided configuration,
// repository storage, and logger. The returned Service implements all event management
// operations defined in the Service interface.
func NewService(config config.Service, storage *repository.Storage, logger logger.Logger) Service {
	return impl.NewService(config, storage, logger)
}
