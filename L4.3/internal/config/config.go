// Package config provides application configuration loading and default fallback.
//
// It uses Viper to read configuration from a config file (config.yaml or other supported formats)
// and ensures that all critical settings have sensible defaults if missing or empty.
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// App holds all configuration sections for the application.
type App struct {
	Logger  Logger  // Logger configuration
	Server  Server  // HTTP server configuration
	Service Service // Business logic / service configuration
	Storage Storage // Persistent storage configuration
}

// Logger contains configuration for the structured logger.
type Logger struct {
	LogDir string // Directory where logs are stored
	Debug  bool   // Enables debug logging if true
}

// Server contains configuration parameters for the HTTP server.
type Server struct {
	Port            string        // Port to listen on
	ReadTimeout     time.Duration // Maximum duration for reading a request
	WriteTimeout    time.Duration // Maximum duration before timing out writes
	MaxHeaderBytes  int           // Maximum size of request headers in bytes
	ShutdownTimeout time.Duration // Timeout for graceful server shutdown
}

// Service contains configuration for the business logic layer.
type Service struct {
	MaxEventsPerUser int // Maximum number of events a user can create
}

// Storage contains configuration for the storage layer.
type Storage struct {
	ExpectedUsers    int // Expected number of users for preallocation / sizing
	MaxEventsPerUser int // Maximum events per user in storage
	MaxEventsPerDay  int // Maximum events per day in storage
}

// Load reads the configuration from a file and returns an App instance.
//
// The configuration file must exist; if it cannot be read, an error is returned.
// For any fields missing or empty within the file, default values are applied
// to ensure the application has all required settings.
func Load() (App, error) {

	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return App{}, fmt.Errorf("viper: %v", err)
	}

	logger := loggerConfig()
	server := serverConfig()
	service := serviceConfig()
	storage := storageConfig()

	failsafe(&logger, &server, &service, &storage)

	return App{
		Logger:  logger,
		Server:  server,
		Service: service,
		Storage: storage,
	}, nil

}

// loggerConfig reads logger configuration from Viper.
func loggerConfig() Logger {
	return Logger{
		LogDir: viper.GetString("app.logger.log_directory"),
		Debug:  viper.GetBool("app.logger.debug_mode"),
	}
}

// serverConfig reads server configuration from Viper.
func serverConfig() Server {
	return Server{
		Port:            viper.GetString("app.server.port"),
		ReadTimeout:     viper.GetDuration("app.server.read_timeout"),
		WriteTimeout:    viper.GetDuration("app.server.write_timeout"),
		MaxHeaderBytes:  viper.GetInt("app.server.max_header_bytes"),
		ShutdownTimeout: viper.GetDuration("app.server.shutdown_timeout"),
	}
}

// serviceConfig reads service configuration from Viper.
func serviceConfig() Service {
	return Service{
		MaxEventsPerUser: viper.GetInt("app.service.max_events_per_user"),
	}
}

// storageConfig reads storage configuration from Viper.
func storageConfig() Storage {
	return Storage{
		ExpectedUsers:   viper.GetInt("app.storage.expected_users"),
		MaxEventsPerDay: viper.GetInt("app.storage.max_events_per_day"),
	}
}

// failsafe fills in default values for missing configuration fields.
//
// This ensures the application can still run even if parts of the config file
// are missing or empty. It prints informative messages for any field that
// is using a default value.
func failsafe(logger *Logger, server *Server, service *Service, storage *Storage) {

	if len(viper.AllSettings()) == 0 {

		fmt.Println("config file is empty, switching to default values")

		*logger = Logger{Debug: true}
		*server = Server{Port: "8080", ReadTimeout: 5 * time.Second, WriteTimeout: 10 * time.Second, MaxHeaderBytes: 1048576, ShutdownTimeout: 15 * time.Second}
		*service = Service{MaxEventsPerUser: 100}
		*storage = Storage{ExpectedUsers: 100, MaxEventsPerDay: 100, MaxEventsPerUser: 100}

		return

	}

	if !viper.IsSet("app.logger.debug_mode") {
		fmt.Println("logger.debug_mode missing, switching to default 'true'")
		logger.Debug = true
	}

	if !viper.IsSet("app.server.port") {
		fmt.Println("server.port missing, switching to default '8080'")
		server.Port = "8080"
	}
	if !viper.IsSet("app.server.read_timeout") {
		fmt.Println("server.read_timeout missing, switching to default 5s")
		server.ReadTimeout = 5 * time.Second
	}
	if !viper.IsSet("app.server.write_timeout") {
		fmt.Println("server.write_timeout missing, switching to default 10s")
		server.WriteTimeout = 10 * time.Second
	}
	if !viper.IsSet("app.server.max_header_bytes") {
		fmt.Println("server.max_header_bytes missing, switching to default 1MB")
		server.MaxHeaderBytes = 1048576
	}
	if !viper.IsSet("app.server.shutdown_timeout") {
		fmt.Println("server.shutdown_timeout missing, switching to default 15s")
		server.ShutdownTimeout = 15 * time.Second
	}

	if !viper.IsSet("app.service.max_events_per_user") {
		fmt.Println("service.max_events_per_user missing, switching to default 100")
		service.MaxEventsPerUser = 100
	}

	if !viper.IsSet("app.storage.expected_users") {
		fmt.Println("storage.expected_users missing, switching to default 100")
		storage.ExpectedUsers = 100
	}
	if !viper.IsSet("app.storage.max_events_per_day") {
		fmt.Println("storage.max_events_per_day missing, switching to default 100")
		storage.MaxEventsPerDay = 100
	}

}
