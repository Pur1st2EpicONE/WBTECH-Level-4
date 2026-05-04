// Package repository provides the abstraction for data storage operations,
// including CRUD operations on events and user event queries.
package repository

import (
	"fmt"
	"time"

	"L4.3/internal/config"
	"L4.3/internal/models"
	"L4.3/internal/repository/memory"
	"L4.3/internal/repository/postgres"
	"L4.3/pkg/logger"
	"github.com/wb-go/wbf/dbpg"
)

type Archive interface {
	SaveEvents(events []models.Event) error
	Close()
}

type Memory interface {
	CreateEvent(event *models.Event) (string, error)
	UpdateEvent(event *models.Event) error
	DeleteEvent(meta *models.Meta) error
	DeleteEvents(ids []string) error
	GetEventByID(eventID string) *models.Event
	CountUserEvents(userID int) (int, error)
	GetEvents(meta *models.Meta, period models.Period) ([]models.Event, error)
	GetExpiredEvents(before time.Time) ([]models.Event, error)
	Close()
}

type Storage struct {
	Archive Archive
	Memory  Memory
}

func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) *Storage {
	return &Storage{
		Archive: postgres.NewStorage(logger, config, db),
		Memory:  memory.NewStorage(config, logger),
	}
}

func ConnectDB(config config.Storage) (*dbpg.DB, error) {

	options := &dbpg.Options{
		MaxOpenConns:    config.MaxOpenConns,
		MaxIdleConns:    config.MaxIdleConns,
		ConnMaxLifetime: config.ConnMaxLifetime,
	}

	db, err := dbpg.New(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.Username, config.Password, config.DBName, config.SSLMode), nil, options)
	if err != nil {
		return nil, fmt.Errorf("database driver not found or DSN invalid: %w", err)
	}

	if err := db.Master.Ping(); err != nil {
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	return db, nil

}
