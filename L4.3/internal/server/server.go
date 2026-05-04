// Package server provides an abstraction over HTTP servers.
// It defines a Server interface and a constructor to create a new server instance.
package server

import (
	"net/http"

	"L4.3/internal/config"
	"L4.3/internal/server/httpserver"
	"L4.3/pkg/logger"
)

// Server defines the behavior expected from an HTTP server instance.
// Any implementation must provide methods to start the server and to shut it down gracefully.
type Server interface {
	// Run starts the HTTP server and blocks until the server exits or an error occurs.
	// Returns a non-nil error if the server fails to start or stops unexpectedly.
	Run() error

	// Shutdown gracefully stops the server, waiting for active requests to finish.
	// Should be called when the server needs to stop in response to a signal or application shutdown.
	Shutdown()
}

// NewServer creates and returns a new Server instance.
// It takes server configuration, an HTTP handler, and a logger as input parameters.
// The returned Server is ready to be run using its Run() method.
func NewServer(config config.Server, handler http.Handler, logger logger.Logger) Server {
	return httpserver.NewServer(config, handler, logger)
}
