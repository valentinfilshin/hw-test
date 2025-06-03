package config

import (
	"context"
	"fmt"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
}

type LoggerConf struct {
	Level     string `config:"level"`
	AddSource bool   `config:"add_source"`
}

type StorageConf struct {
	Type     string `config:"type"`
	DB       string `config:"db"`
	User     string `config:"user"`
	Password string `config:"password"`
	Host     string `config:"host"`
	Port     string `config:"port"`
}

func LoadConfig(configFile string) (*Config, error) {
	cfg := &Config{}
	err := confita.NewLoader(file.NewBackend(configFile)).Load(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return cfg, nil
}
