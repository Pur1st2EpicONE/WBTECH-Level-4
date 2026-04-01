// Package config provides functionality for loading configuration
// parameters from a config file using the Viper library.
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config defines runtime configuration parameters for the sort utility.
// Fields correspond to keys in the configuration file.
type Config struct {
	ChunkSize int `mapstructure:"chunk_size"` // Number of lines per chunk.
	Workers   int `mapstructure:"workers"`    // Number of concurrent workers for sorting.
}

// Load reads the configuration file from the current working directory.
func Load() (*Config, error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	cfg := new(Config)
	if err := viper.ReadInConfig(); err != nil {
		return cfg, fmt.Errorf("failed to load config: %v", err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return cfg, nil
}
