// Package config provides application configuration loading
// from YAML files using Viper.
package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// Config represents full application configuration.
type Config struct {
	Logger Logger
	Server Server
}

// Logger holds logging configuration.
type Logger struct {
	LogDir string
	Debug  bool
}

// Server holds HTTP server configuration options.
type Server struct {
	Port            string
	ReadTimeout     time.Duration // Maximum duration for reading a request
	WriteTimeout    time.Duration // Maximum duration before timing out writes
	MaxHeaderBytes  int           // Maximum size of request headers in bytes
	ShutdownTimeout time.Duration // Timeout for graceful server shutdown
}

// Load reads configuration from config.yaml.
func Load() (Config, error) {

	viper.AddConfigPath(".")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return Config{}, fmt.Errorf("viper: %v", err)
	}

	logger := loggerConfig()
	server := serverConfig()

	return Config{
		Logger: logger,
		Server: server,
	}, nil

}

// loggerConfig builds logger configuration from viper.
func loggerConfig() Logger {
	return Logger{
		LogDir: viper.GetString("logger.log_directory"),
		Debug:  viper.GetBool("logger.debug_mode"),
	}
}

// serverConfig builds server configuration from viper.
func serverConfig() Server {
	return Server{
		Port:            viper.GetString("server.port"),
		ReadTimeout:     viper.GetDuration("server.read_timeout"),
		WriteTimeout:    viper.GetDuration("server.write_timeout"),
		MaxHeaderBytes:  viper.GetInt("server.max_header_bytes"),
		ShutdownTimeout: viper.GetDuration("server.shutdown_timeout"),
	}
}
