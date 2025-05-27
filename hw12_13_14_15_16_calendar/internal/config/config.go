package config

import (
	"context"
	"errors"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

type Config struct {
	Logger LoggerConf
}

type LoggerConf struct {
	Level     string `config:"level"`
	AddSource bool   `config:"add_source"`
}

func LoadConfig(configFile string) (*Config, error) {
	cfg := &Config{}
	err := confita.NewLoader(file.NewBackend(configFile)).Load(context.Background(), &cfg)
	if err != nil {
		return nil, errors.New("failed to load config")
	}
	return cfg, nil
}
