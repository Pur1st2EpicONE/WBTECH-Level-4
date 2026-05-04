package app

import (
	"time"

	"L4.3/internal/repository"
	"L4.3/pkg/logger"
)

type Archiver struct {
	storage *repository.Storage
	logger  logger.Logger
	ticker  *time.Ticker
	stop    chan struct{}
}

func NewArchiver(storage *repository.Storage, logger logger.Logger, interval time.Duration) *Archiver {
	return &Archiver{
		storage: storage,
		logger:  logger,
		ticker:  time.NewTicker(interval),
		stop:    make(chan struct{}),
	}
}

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

func (a *Archiver) Stop() {
	close(a.stop)
	a.ticker.Stop()
}
