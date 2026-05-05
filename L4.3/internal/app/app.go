// Package app defines the main application composition root and lifecycle management.
//
// It is responsible for bootstrapping all core components (configuration, logger,
// database, storage, services, HTTP server), wiring dependencies, handling OS signals,
// and coordinating graceful shutdown.
package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"L4.3/internal/config"
	"L4.3/internal/handler"
	"L4.3/internal/repository"
	"L4.3/internal/server"
	"L4.3/internal/service"
	"L4.3/pkg/logger"
	"github.com/pressly/goose/v3"
	"github.com/wb-go/wbf/dbpg"
)

const archiveInterval = 10 * time.Minute

// App represents the root application container.
type App struct {
	logger   logger.Logger       // structured application logger
	server   server.Server       // HTTP server handling incoming requests
	storage  *repository.Storage // data access layer (memory + archive)
	ctx      context.Context     // root context for cancellation
	cancel   context.CancelFunc  // cancels root context
	wg       *sync.WaitGroup     // coordinates goroutines lifecycle
	archiver *Archiver           // background archiving worker
}

// Boot initializes the application and wires all dependencies.
func Boot() *App {

	config, err := config.Load()
	if err != nil {
		log.Fatalf("app — failed to load configs: %v", err)
	}

	logger := logger.NewLogger(config.App.Logger)

	db, err := bootstrapDB(logger, config.App.Storage)
	if err != nil {
		logger.LogFatal("app — failed to connect to database", err, "layer", "app")
	}

	server, storage := wireApp(db, config.App, logger)
	archiver := NewArchiver(storage, logger, archiveInterval)

	ctx, cancel := newContext(logger)
	wg := new(sync.WaitGroup)

	return &App{
		logger:   logger,
		server:   server,
		storage:  storage,
		archiver: archiver,
		ctx:      ctx,
		cancel:   cancel,
		wg:       wg,
	}

}

// bootstrapDB establishes database connection and applies migrations.
func bootstrapDB(logger logger.Logger, config config.Storage) (*dbpg.DB, error) {

	db, err := repository.ConnectDB(config)
	if err != nil {
		return nil, err
	}

	logger.LogInfo("app — connected to database", "layer", "app")

	if err := goose.SetDialect(config.Dialect); err != nil {
		return nil, fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db.Master, config.MigrationsDir); err != nil {
		return nil, fmt.Errorf("failed to apply goose migrations: %w", err)
	}

	logger.Debug("app — migrations applied", "layer", "app")

	return db, nil

}

// wireApp constructs and wires core application layers.
func wireApp(db *dbpg.DB, config config.App, logger logger.Logger) (server.Server, *repository.Storage) {
	storage := repository.NewStorage(logger, config.Storage, db)
	service := service.NewService(config.Service, storage, logger)
	handler := handler.NewHandler(service, logger)
	server := server.NewServer(config.Server, handler, logger)
	return server, storage
}

// newContext creates a root context that is canceled on OS termination signals.
func newContext(logger logger.Logger) (context.Context, context.CancelFunc) {

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-sigCh
		logger.LogInfo("app — received signal "+sig.String()+", initiating graceful shutdown", "layer", "app")
		cancel()
	}()

	return ctx, cancel

}

// Run starts the application components and blocks until shutdown.
func (a *App) Run() {

	a.wg.Go(func() {
		if err := a.server.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.LogFatal("server run failed", err, "layer", "app")
		}
	})

	a.wg.Go(a.archiver.Start)

	<-a.ctx.Done()

	a.logger.LogInfo("app — shutting down...", "layer", "app")
	a.Stop()

	a.wg.Wait()

}

// Stop gracefully shuts down all application components.
func (a *App) Stop() {

	a.archiver.Stop()
	a.server.Shutdown()

	a.storage.Memory.Close()
	a.storage.Archive.Close()

	a.logger.Close()

}
