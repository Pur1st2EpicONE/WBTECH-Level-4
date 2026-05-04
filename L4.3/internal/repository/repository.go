// Package repository provides the abstraction for data storage operations,
// including CRUD operations on events and user event queries.
package repository

import (
	"L4.3/internal/config"
	"L4.3/internal/models"
	"L4.3/internal/repository/memory"
	"L4.3/pkg/logger"
)

// Storage defines the interface for interacting with the application's
// persistent or in-memory storage layer.
type Storage interface {
	// CreateEvent stores a new event and returns its unique ID.
	CreateEvent(event *models.Event) (string, error)

	// UpdateEvent updates an existing event identified by its ID.
	UpdateEvent(event *models.Event) error

	// DeleteEvent removes an event based on metadata (user ID + event ID).
	DeleteEvent(meta *models.Meta) error

	// GetEventByID retrieves an event by its unique ID.
	// Returns nil if no event is found.
	GetEventByID(eventID string) *models.Event

	// CountUserEvents returns the number of events associated with a user.
	CountUserEvents(userID int) (int, error)

	// GetEvents retrieves all events for a user filtered by a given period
	// (day, week, month).
	GetEvents(meta *models.Meta, period models.Period) ([]models.Event, error)

	// Close cleans up any resources held by the storage.
	Close()
}

// NewStorage creates a new Storage instance. If db is nil, it returns
// an in-memory implementation. Panics if an unsupported storage type is provided.
func NewStorage(db any, config config.Storage, logger logger.Logger) Storage {
	if db == nil {
		return memory.NewStorage(config, logger)
	} else {
		panic("unsupported storage type")
	}
}
