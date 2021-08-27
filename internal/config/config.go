package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug             bool          `envconfig:"DEBUG" default:"false"`
	JSONLogs          bool          `envconfig:"JSON_LOGS" default:"false"`
	Port              int           `envconfig:"PORT" default:"8000"`
	Timeout           time.Duration `envconfig:"TIMEOUT" default:"30s"`         // seconds
	TimeoutPerHost    time.Duration `envconfig:"TIMEOUT_PER_HOST" default:"5s"` // seconds
	StartRPS          int           `envconfig:"START_RPS" default:"20"`
	IncreasingStepRPS int           `envconfig:"INCREASING_STEP_RPS" default:"10"`
}

func New() (*Config, error) {
	config := &Config{}
	if err := envconfig.Process("TANK", config); err != nil {
		return nil, err
	}

	return config, nil
}
