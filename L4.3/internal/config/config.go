// Package config provides application configuration structures and loading logic.
//
// It supports loading configuration from multiple sources:
//   - YAML config file
//   - environment file (.env)
//   - OS environment variables (highest priority for sensitive data)
//
// The package relies on mapstructure tags for unmarshalling.
package config

import (
	"fmt"
	"os"
	"time"

	wbf "github.com/wb-go/wbf/config"
)

// Config is the root configuration structure of the application.
type Config struct {
	App App `mapstructure:"app"`
}

// App groups all application-level configuration sections.
type App struct {
	Logger  Logger  `mapstructure:"logger"`
	Server  Server  `mapstructure:"server"`
	Service Service `mapstructure:"service"`
	Storage Storage `mapstructure:"storage"`
}

// Logger defines logging configuration.
type Logger struct {
	LogDir string `mapstructure:"log_directory"` // directory for log files (empty = stdout)
	Debug  bool   `mapstructure:"debug_mode"`    // enables debug-level logging
}

// Server defines HTTP server configuration.
type Server struct {
	Port            string        `mapstructure:"port"`             // listening port
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`     // maximum duration for reading request
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`    // maximum duration for writing response
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"` // maximum size of request headers
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"` // graceful shutdown timeout
}

// Service defines business logic constraints.
type Service struct {
	MaxEventsPerUser int `mapstructure:"max_events_per_user"` // per-user event limit
}

// Storage defines configuration for both in-memory and persistent storage.
type Storage struct {
	ExpectedUsers      int           `mapstructure:"expected_users"`             // capacity hint for in-memory storage
	MaxEventsPerDay    int           `mapstructure:"max_events_per_day"`         // per-day event limit
	Dialect            string        `mapstructure:"goose_dialect"`              // migration dialect
	MigrationsDir      string        `mapstructure:"goose_migrations_directory"` // path to migrations
	Host               string        `mapstructure:"host"`                       // DB host
	Port               string        `mapstructure:"port"`                       // DB port
	Username           string        `mapstructure:"username"`                   // DB username (overridden by env)
	Password           string        `mapstructure:"password"`                   // DB password (overridden by env)
	DBName             string        `mapstructure:"dbname"`                     // database name
	SSLMode            string        `mapstructure:"sslmode"`                    // SSL mode
	MaxOpenConns       int           `mapstructure:"max_open_conns"`             // max open connections
	MaxIdleConns       int           `mapstructure:"max_idle_conns"`             // max idle connections
	ConnMaxLifetime    time.Duration `mapstructure:"conn_max_lifetime"`          // connection lifetime
	QueryRetryStrategy RetryStrategy `mapstructure:"query_retry_strategy"`       // retry policy for queries
	TxRetryStrategy    RetryStrategy `mapstructure:"tx_retry_strategy"`          // retry policy for transactions
}

// RetryStrategy defines retry behavior for database operations.
type RetryStrategy struct {
	Attempts int           `mapstructure:"attempts"` // number of retry attempts
	Delay    time.Duration `mapstructure:"delay"`    // initial delay between attempts
	Backoff  float64       `mapstructure:"backoff"`  // exponential backoff multiplier
}

// Load reads configuration from files and environment variables.
func Load() (Config, error) {

	cfg := wbf.New()

	if err := cfg.LoadConfigFiles("./config.yaml"); err != nil {
		return Config{}, err
	}

	if err := cfg.LoadEnvFiles(".env"); err != nil && !cfg.GetBool("docker") {
		return Config{}, err
	}

	var conf Config

	if err := cfg.Unmarshal(&conf); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	loadEnvs(&conf)

	return conf, nil

}

// loadEnvs overrides configuration fields with OS environment variables.
func loadEnvs(conf *Config) {
	conf.App.Storage.Username = os.Getenv("DB_USER")
	conf.App.Storage.Password = os.Getenv("DB_PASSWORD")
}
