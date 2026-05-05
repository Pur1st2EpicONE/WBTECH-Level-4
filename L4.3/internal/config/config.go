package config

import (
	"fmt"
	"os"
	"time"

	wbf "github.com/wb-go/wbf/config"
)

type Config struct {
	App App `mapstructure:"app"`
}

type App struct {
	Logger  Logger  `mapstructure:"logger"`
	Server  Server  `mapstructure:"server"`
	Service Service `mapstructure:"service"`
	Storage Storage `mapstructure:"storage"`
}

type Logger struct {
	LogDir string `mapstructure:"log_directory"`
	Debug  bool   `mapstructure:"debug_mode"`
}

type Server struct {
	Port            string        `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type Service struct {
	MaxEventsPerUser int `mapstructure:"max_events_per_user"`
}

type Storage struct {
	ExpectedUsers   int `mapstructure:"expected_users"`
	MaxEventsPerDay int `mapstructure:"max_events_per_day"`

	Dialect       string `mapstructure:"goose_dialect"`
	MigrationsDir string `mapstructure:"goose_migrations_directory"`

	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`

	MaxOpenConns    int           `mapstructure:"max_open_conns"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`

	QueryRetryStrategy RetryStrategy `mapstructure:"query_retry_strategy"`
	TxRetryStrategy    RetryStrategy `mapstructure:"tx_retry_strategy"`
}

type RetryStrategy struct {
	Attempts int           `mapstructure:"attempts"`
	Delay    time.Duration `mapstructure:"delay"`
	Backoff  float64       `mapstructure:"backoff"`
}

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

func loadEnvs(conf *Config) {
	conf.App.Storage.Username = os.Getenv("DB_USER")
	conf.App.Storage.Password = os.Getenv("DB_PASSWORD")
}
