package app

import (
	"time"

	"L4.3/internal/repository"
	"L4.3/pkg/logger"
)

// Archiver is a background worker responsible for moving expired events
// from primary (memory) storage to archive storage.
//
// It runs periodically using a ticker and performs a batch archive operation:
//  1. fetch expired events
//  2. persist them to archive storage
//  3. delete them from primary storage
//
// The worker is designed to be run in a separate goroutine.
type Archiver struct {
	storage *repository.Storage // access to memory and archive storages
	logger  logger.Logger       // structured logger
	ticker  *time.Ticker        // triggers periodic execution
	stop    chan struct{}       // signals worker termination
}

// NewArchiver constructs a new Archiver instance.
//
// interval defines how often the archiving job is executed.
// The ticker is started immediately.
func NewArchiver(storage *repository.Storage, logger logger.Logger, interval time.Duration) *Archiver {
	return &Archiver{
		storage: storage,
		logger:  logger,
		ticker:  time.NewTicker(interval),
		stop:    make(chan struct{}),
	}
}

// Start runs the archiving loop.
//
// It blocks and should typically be executed in a goroutine.
// The loop listens for:
//   - ticker events to trigger archiving
//   - stop signal to terminate execution
func (a *Archiver) Start() {
	for {
		select {
		case <-a.ticker.C:
			a.run()
		case <-a.stop:
			return
		}
	}
}

// run performs a single archiving cycle.
//
// Workflow:
//   - determines cutoff time (current UTC time)
//   - retrieves expired events from memory storage
//   - saves them to archive storage
//   - deletes them from memory storage
func (a *Archiver) run() {

	cutoff := time.Now().UTC()

	events, err := a.storage.Memory.GetExpiredEvents(cutoff)
	if err != nil {
		a.logger.LogError("archiver — failed to get expired events", err, "layer", "app")
		return
	}

	if len(events) == 0 {
		a.logger.Debug("archiver — nothing to archive", "layer", "app")
		return
	}

	if err := a.storage.Archive.SaveEvents(events); err != nil {
		a.logger.LogError("archiver — failed to save old events", err, "layer", "app")
		return
	}

	ids := make([]string, 0, len(events))
	for _, e := range events {
		ids = append(ids, e.Meta.EventID)
	}

	if err := a.storage.Memory.DeleteEvents(ids); err != nil {
		a.logger.LogError("archiver — failed to delete events", err, "layer", "app")
		return
	}

	a.logger.Debug("archiver — events archived", "count", len(events), "layer", "app")

}

// Stop signals the archiver to terminate and stops the ticker.
func (a *Archiver) Stop() {
	close(a.stop)
	a.ticker.Stop()
}
