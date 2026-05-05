// Package app contains the application bootstrap logic,
// including dependency wiring, signal handling, and graceful shutdown.
package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"L4.4/internal/config"
	"L4.4/internal/handler"
	"L4.4/internal/logger"
	"L4.4/internal/server"
)

// App represents the main application container holding
// runtime dependencies and lifecycle context.
type App struct {
	logger  logger.Logger
	logFile *os.File
	server  server.Server
	ctx     context.Context
	cancel  context.CancelFunc
}

// Boot initializes application dependencies and returns a ready-to-run App instance.
func Boot() *App {

	config, err := config.Load()
	if err != nil {
		log.Fatalf("app — failed to load configs: %v", err)
	}

	logger, logFile := logger.NewLogger(config.Logger)

	return wireApp(logger, logFile, config)

}

// wireApp builds application dependencies and connects core components.
func wireApp(logger logger.Logger, logFile *os.File, config config.Config) *App {

	ctx, cancel := newContext(logger)
	server := server.NewServer(logger, config.Server, handler.NewHandler())

	return &App{
		logger:  logger,
		logFile: logFile,
		server:  server,
		ctx:     ctx,
		cancel:  cancel,
	}

}

// newContext creates a cancellable context and listens for OS termination signals.
func newContext(logger logger.Logger) (context.Context, context.CancelFunc) {

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := <-sigCh
		sigString := sig.String()
		if sig == syscall.SIGTERM {
			sigString = "terminate" // sig.String() returns the SIGTERM string in past tense for some reason
		}
		logger.LogInfo("app — received signal "+sigString+", initiating graceful shutdown", "layer", "app")
		cancel()
	}()

	return ctx, cancel

}

// Run starts the HTTP server and blocks until shutdown signal is received.
func (a *App) Run() {

	go func() {
		if err := a.server.Run(); err != nil {
			a.logger.LogFatal("server run failed", err, "layer", "app")
		}
	}()

	<-a.ctx.Done()

	a.Stop()

}

// Stop gracefully shuts down the server and releases resources.
func (a *App) Stop() {

	a.server.Shutdown()

	if a.logFile != nil && a.logFile != os.Stdout {
		_ = a.logFile.Close()
	}

}
