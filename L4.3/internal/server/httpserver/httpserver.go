// Package httpserver provides a concrete implementation of an HTTP server.
// It wraps the standard net/http.Server and adds graceful shutdown and logging capabilities.
package httpserver

import (
	"context"
	"net/http"
	"time"

	"L4.3/internal/config"
	"L4.3/pkg/logger"
)

// HttpServer represents an HTTP server with logging and graceful shutdown support.
// It wraps an underlying http.Server and keeps track of shutdown timeout and logger.
type HttpServer struct {
	srv             *http.Server  // srv is the underlying HTTP server that handles incoming requests.
	shutdownTimeout time.Duration // shutdownTimeout specifies how long to wait for active connections to finish during shutdown.
	logger          logger.Logger // logger is used to log server events, errors, and shutdown information.
}

// NewServer creates a new HttpServer instance with the specified configuration, HTTP handler, and logger.
// The configuration provides the port, read/write timeouts, max header size, and shutdown timeout.
func NewServer(config config.Server, handler http.Handler, logger logger.Logger) *HttpServer {
	server := new(HttpServer)
	server.srv = &http.Server{
		Addr:           ":" + config.Port,
		Handler:        handler,
		ReadTimeout:    config.ReadTimeout,
		WriteTimeout:   config.WriteTimeout,
		MaxHeaderBytes: config.MaxHeaderBytes,
	}
	server.shutdownTimeout = config.ShutdownTimeout
	server.logger = logger
	return server
}

// Run starts the HTTP server and begins listening for requests.
// It blocks until the server is stopped or an error occurs. Errors are returned to the caller.
func (s *HttpServer) Run() error {
	s.logger.LogInfo("server — receiving requests", "layer", "server")
	return s.srv.ListenAndServe()
}

// Shutdown gracefully stops the server, allowing active connections to complete within the configured timeout.
// Logs either a successful shutdown or any errors encountered during the shutdown process.
func (s *HttpServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()
	if err := s.srv.Shutdown(ctx); err != nil {
		s.logger.LogError("server — failed to shutdown gracefully", err, "layer", "server")
	} else {
		s.logger.LogInfo("server — shutdown complete", "layer", "server")
	}
}
