// Package impl provides the concrete implementation of the Service interface.
// It contains the business logic for creating, updating, deleting, and retrieving events.
package impl

import (
	"L4.3/internal/config"
	"L4.3/internal/models"
	"L4.3/internal/repository"
	"L4.3/pkg/logger"
)

// Service is the implementation of the event management service.
// It interacts with the repository layer to perform CRUD operations on events
// and enforces business rules such as user limits and validation checks.
type Service struct {
	Storage          *repository.Storage
	logger           logger.Logger
	maxEventsPerUser int
	reminderCh       chan models.Reminder
	stopCh           chan struct{}
}

// NewService creates a new Service instance with the provided configuration, storage, and logger.
// The maxEventsPerUser field is set from the configuration and enforces limits on event creation.
func NewService(config config.Service, storage *repository.Storage, logger logger.Logger) *Service {

	s := &Service{
		Storage:          storage,
		logger:           logger,
		maxEventsPerUser: config.MaxEventsPerUser,
		reminderCh:       make(chan models.Reminder, 100),
		stopCh:           make(chan struct{}),
	}

	go s.reminderWorker()

	return s

}
