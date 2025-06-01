package logger

import (
	"log/slog"
	"os"

	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/config"
)

type Logger struct {
	*slog.Logger
}

func (l *Logger) Info(msg string) {
	l.Logger.Info(msg)
}

func (l *Logger) Error(msg string) {
	l.Logger.Error(msg)
}

func New(cfg config.LoggerConf) *Logger {
	var level slog.Level

	err := level.UnmarshalText([]byte(cfg.Level))
	if err != nil {
		level = slog.LevelInfo
	}

	slogOptions := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, slogOptions))

	return &Logger{logger}
}
