package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug    bool `envconfig:"DEBUG" default:"true"`
	JSONLogs bool `envconfig:"JSON_LOGS" default:"true"`
	Port     int  `envconfig:"PORT" default:"8000"`
}

func New() (*Config, error) {
	config := &Config{}
	if err := envconfig.Process("TANK", config); err != nil {
		return nil, err
	}

	return config, nil
}
