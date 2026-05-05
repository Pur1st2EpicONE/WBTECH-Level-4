// Package async provides an asynchronous structured logger built on top of slog.
// It uses a buffered channel and a background worker to decouple log producers
// from I/O operations, reducing latency in hot paths.
package async

import (
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"L4.3/internal/config"
)

// entry represents a single log record passed to the async worker.
type entry struct {
	level slog.Level // log level (debug, info, warn, error)
	msg   string     // log message
	args  []any      // structured key-value arguments
	err   error      // optional error to attach
}

// Logger is an asynchronous logger that buffers log entries and writes them
// using a background worker.
//
// It wraps a synchronous slog.Logger and ensures non-blocking logging
// for most operations. On shutdown, it drains all pending log entries.
type Logger struct {
	ch    chan entry // buffered channel for log entries
	wg    sync.WaitGroup
	done  chan struct{} // signal to stop worker and flush logs
	sync  *slog.Logger  // underlying synchronous logger
	file  *os.File      // log file (or stdout)
	debug bool          // enables debug-level logging
}

// NewLogger initializes a new asynchronous Logger based on the provided config.
//
// If LogDir is set, logs are written to a file inside that directory.
// Otherwise, logs are written to stdout.
//
// A background worker is started immediately.
//
// The logger uses JSON formatting via slog.
func NewLogger(cfg config.Logger) *Logger {

	var logDest *os.File
	if cfg.LogDir == "" {
		logDest = os.Stdout
	} else {
		logDest = openFile(cfg.LogDir)
		if logDest == nil {
			logDest = os.Stdout
		}
	}

	level := slog.LevelInfo
	if cfg.Debug {
		level = slog.LevelDebug
	}
	handler := slog.NewJSONHandler(logDest, &slog.HandlerOptions{Level: level})
	syncLogger := slog.New(handler)

	al := &Logger{
		ch:    make(chan entry, 10000),
		done:  make(chan struct{}),
		sync:  syncLogger,
		file:  logDest,
		debug: cfg.Debug,
	}

	al.wg.Add(1)
	go al.worker()
	return al
}

// worker is a background goroutine that processes log entries from the channel.
//
// It continues processing until a shutdown signal is received via done,
// after which it drains the channel before exiting.
func (l *Logger) worker() {
	defer l.wg.Done()
	for {
		select {
		case e := <-l.ch:
			l.write(e)
		case <-l.done:
			for {
				select {
				case e := <-l.ch:
					l.write(e)
				default:
					return
				}
			}
		}
	}
}

// write sends a single log entry to the underlying slog.Logger.
func (l *Logger) write(e entry) {
	if e.level == slog.LevelDebug && !l.debug {
		return
	}
	args := e.args
	if e.err != nil {
		args = append(args, "err", e.err.Error())
	}
	switch e.level {
	case slog.LevelError:
		l.sync.Error(e.msg, args...)
	case slog.LevelWarn:
		l.sync.Warn(e.msg, args...)
	case slog.LevelInfo:
		l.sync.Info(e.msg, args...)
	case slog.LevelDebug:
		l.sync.Debug(e.msg, args...)
	}
}

// LogFatal logs an error-level message and terminates the application.
//
// It attempts to enqueue the log entry, waits briefly to allow the worker
// to process it, then logs synchronously and calls os.Exit(1).
//
// This method should be used only for unrecoverable errors.
func (l *Logger) LogFatal(msg string, err error, args ...any) {
	l.ch <- entry{level: slog.LevelError, msg: msg, err: err, args: args}
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		args = append(args, "err", err.Error())
	}
	l.sync.Error(msg, args...)
	os.Exit(1)
}

// LogError enqueues an error-level log message.
func (l *Logger) LogError(msg string, err error, args ...any) {
	l.ch <- entry{level: slog.LevelError, msg: msg, err: err, args: args}
}

// LogWarn enqueues a warning-level log message.
func (l *Logger) LogWarn(msg string, args ...any) {
	l.ch <- entry{level: slog.LevelWarn, msg: msg, args: args}
}

// LogInfo enqueues an info-level log message.
func (l *Logger) LogInfo(msg string, args ...any) {
	l.ch <- entry{level: slog.LevelInfo, msg: msg, args: args}
}

// Debug enqueues a debug-level log message if debug mode is enabled.
func (l *Logger) Debug(msg string, args ...any) {
	if !l.debug {
		return
	}
	l.ch <- entry{level: slog.LevelDebug, msg: msg, args: args}
}

// Close gracefully shuts down the logger.
//
// It signals the worker to stop, waits for all pending log entries
// to be processed, and closes the log file if necessary.
func (l *Logger) Close() {
	close(l.done)
	l.wg.Wait()
	if l.file != nil && l.file != os.Stdout {
		_ = l.file.Close()
	}
}

// openFile ensures the log directory exists and opens the log file.
//
// If any error occurs, it logs the issue using the default slog logger
// and returns nil, indicating fallback to stdout.
func openFile(logDir string) *os.File {
	if err := os.MkdirAll(logDir, 0777); err != nil {
		slog.Error("logger — failed to create log directory switching to stdout", "error", err)
		return nil
	}
	logPath := filepath.Join(logDir, "app.log")
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		slog.Error("logger — failed to create log file switching to stdout", "error", err)
		return nil
	}
	return logFile
}
