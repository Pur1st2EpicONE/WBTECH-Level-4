// Package app defines the main application structure and lifecycle management.
//
// It handles application initialization, context and signal management, server startup,
// graceful shutdown, and resource cleanup. The App struct encapsulates all components
// required to run the calendar service, including logger, server, storage, and context.
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

// App represents the main application instance, managing its components and lifecycle.
type App struct {
	logger   logger.Logger // Structured logger used throughout the application for info, warning, error, and debug logs
	server   server.Server // HTTP server instance that handles incoming requests
	storage  *repository.Storage
	ctx      context.Context    // Context used for cancellation and graceful shutdown
	cancel   context.CancelFunc // Function to cancel the application context and trigger shutdown
	wg       *sync.WaitGroup    // WaitGroup to synchronize goroutines during server run and shutdown
	archiver *Archiver
}

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
	archiver := NewArchiver(storage, logger, 10*time.Second)

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

func wireApp(db *dbpg.DB, config config.App, logger logger.Logger) (server.Server, *repository.Storage) {
	storage := repository.NewStorage(logger, config.Storage, db)
	service := service.NewService(config.Service, storage, logger)
	handler := handler.NewHandler(service, logger)
	server := server.NewServer(config.Server, handler, logger)
	return server, storage
}

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

func (a *App) Stop() {
	a.archiver.Stop()
	a.server.Shutdown()
	a.storage.Memory.Close()
	a.storage.Archive.Close()
	a.logger.Close()
}
