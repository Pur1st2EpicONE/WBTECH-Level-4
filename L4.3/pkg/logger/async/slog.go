package async

import (
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"

	"L4.3/internal/config"
)

type entry struct {
	level slog.Level
	msg   string
	args  []any
	err   error
}

type Logger struct {
	ch    chan entry
	wg    sync.WaitGroup
	done  chan struct{}
	sync  *slog.Logger
	file  *os.File
	debug bool
}

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

func (l *Logger) LogFatal(msg string, err error, args ...any) {
	l.ch <- entry{level: slog.LevelError, msg: msg, err: err, args: args}
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		args = append(args, "err", err.Error())
	}
	l.sync.Error(msg, args...)
	os.Exit(1)
}

func (l *Logger) LogError(msg string, err error, args ...any) {
	l.ch <- entry{level: slog.LevelError, msg: msg, err: err, args: args}
}

func (l *Logger) LogWarn(msg string, args ...any) {
	l.ch <- entry{level: slog.LevelWarn, msg: msg, args: args}
}

func (l *Logger) LogInfo(msg string, args ...any) {
	l.ch <- entry{level: slog.LevelInfo, msg: msg, args: args}
}

func (l *Logger) Debug(msg string, args ...any) {
	if !l.debug {
		return
	}
	l.ch <- entry{level: slog.LevelDebug, msg: msg, args: args}
}

func (l *Logger) Close() {
	close(l.done)
	l.wg.Wait()
	if l.file != nil && l.file != os.Stdout {
		_ = l.file.Close()
	}
}

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
