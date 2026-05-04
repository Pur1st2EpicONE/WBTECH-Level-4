package postgres

import (
	"L4.3/internal/config"
	"L4.3/pkg/logger"
	"github.com/wb-go/wbf/dbpg"
)

type Storage struct {
	db     *dbpg.DB
	logger logger.Logger
	config config.Storage
}

func NewStorage(logger logger.Logger, config config.Storage, db *dbpg.DB) *Storage {
	return &Storage{db: db, logger: logger, config: config}
}
