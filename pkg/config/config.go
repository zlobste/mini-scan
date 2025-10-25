package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Database  DatabaseConfig  `envconfig:""`
	PubSub    PubSubConfig    `envconfig:""`
	Processor ProcessorConfig `envconfig:""`
}

type DatabaseConfig struct {
	Host     string `envconfig:"POSTGRES_HOST" default:"postgres"`
	Port     int    `envconfig:"POSTGRES_PORT" default:"5432"`
	User     string `envconfig:"POSTGRES_USER" default:"postgres"`
	Password string `envconfig:"POSTGRES_PASSWORD" default:"postgres"`
	DB       string `envconfig:"POSTGRES_DB" default:"scans"`
}

type PubSubConfig struct {
	ProjectID    string `envconfig:"PUBSUB_PROJECT_ID" default:"test-project"`
	EmulatorHost string `envconfig:"PUBSUB_EMULATOR_HOST"`
	Subscription string `envconfig:"PROCESSOR_SUBSCRIPTION" default:"scan-sub"`
}

type ProcessorConfig struct {
	LogLevel     string        `envconfig:"PROCESSOR_LOG_LEVEL" default:"info"`
	BatchSize    int           `envconfig:"PROCESSOR_BATCH_SIZE" default:"100"`
	BatchTimeout time.Duration `envconfig:"PROCESSOR_BATCH_TIMEOUT" default:"5s"`
}

func Load() (*Config, error) {
	cfg := &Config{}
	if err := envconfig.Process("", cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (dc DatabaseConfig) ConnectionString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dc.User,
		dc.Password,
		dc.Host,
		dc.Port,
		dc.DB,
	)
}
