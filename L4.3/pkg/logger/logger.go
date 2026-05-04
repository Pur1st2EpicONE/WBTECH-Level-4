// Package logger provides a unified interface for structured logging across the application.
// It abstracts different logging levels (fatal, error, warning, info, debug) and allows
// consistent logging with additional contextual fields. The logger implementation can
// be swapped (currently using slog) without changing the application code.
package logger

import (
	"L4.3/internal/config"
	"L4.3/pkg/logger/async"
)

// Logger defines the interface for structured logging with multiple severity levels.
// It is used across the application to provide consistent logging with optional
// contextual fields (key-value pairs). Implementations can vary, but must respect
// the semantics of each log level and allow graceful shutdown via Close().
type Logger interface {
	// LogFatal logs a fatal message along with an associated error and optional
	// key-value arguments. Typically used for unrecoverable errors that should
	// terminate the application.
	LogFatal(msg string, err error, args ...any)

	// LogError logs an error message along with an associated error and optional
	// key-value arguments. Used for recoverable errors that need attention.
	LogError(msg string, err error, args ...any)

	// LogWarn logs a warning message with optional key-value arguments. Used
	// for non-critical issues that should be monitored.
	LogWarn(msg string, args ...any)

	// LogInfo logs an informational message with optional key-value arguments.
	// Useful for general runtime information, like startup/shutdown events or
	// major actions.
	LogInfo(msg string, args ...any)

	// Debug logs a debug message with optional key-value arguments. Intended
	// for development and troubleshooting, usually verbose.
	Debug(msg string, args ...any)

	// Close gracefully releases any underlying resources used by the logger,
	// such as open files or network connections.
	Close()
}

// NewLogger creates a new Logger instance using the provided configuration.
// It returns a logger implementation that satisfies the Logger interface.
func NewLogger(config config.Logger) Logger {
	return async.NewLogger(config)
}
