package logger

import (
	"log/slog"
	"os"
	"strings"
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

func New(lvl string, addSource bool) *Logger {
	sLvl := slog.LevelInfo

	switch strings.ToLower(lvl) {
	case "debug":
		sLvl = slog.LevelDebug
	case "info":
		sLvl = slog.LevelInfo
	case "warn":
		sLvl = slog.LevelWarn
	case "error":
		sLvl = slog.LevelError
	}

	slogOptions := &slog.HandlerOptions{
		Level:     sLvl,
		AddSource: addSource,
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, slogOptions))

	return &Logger{logger}
}
