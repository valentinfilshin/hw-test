package config

import (
	"context"
	"fmt"

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

func NewConfig(configFile string) Config {
	cfg := Config{}
	err := confita.NewLoader(file.NewBackend(configFile)).Load(context.Background(), &cfg)
	if err != nil {
		// TODO как правильно логировать ошибки при старте программы?
		_ = fmt.Errorf(
			"failed to load config: %v",
			err,
		)
	}
	return cfg
}
