package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"syscall"
	"time"

	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	// 1. Загружаем конфигурацию
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Создаем логгер
	logg := logger.New(cfg.Logger.Level, cfg.Logger.AddSource)

	// 3. Создаем хранилище
	var storage app.Storage
	if cfg.Storage.Type == "postgres" {
		postgresqlStorage := sqlstorage.New(cfg.Storage.DSN)

		err := postgresqlStorage.Connect()
		if err != nil {
			logg.Error("failed to connect to database: " + err.Error())
			return
		}

		logg.Info("connected to database")

		defer func() {
			err := postgresqlStorage.Close()
			if err != nil {
				logg.Error("failed to close database: " + err.Error())
			}
		}()

		storage = postgresqlStorage
	} else {
		storage = memorystorage.New()
	}

	// 4. Бизнес-логика
	calendar := app.New(logg, storage)

	// 5. Запускаем сервер
	server := internalhttp.NewServer(logg, cfg.Server.Addr, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	logg.Info("calendar is running")

	go func() {
		if err := server.Start(); err != nil {
			logg.Error("failed to start http server: " + err.Error())
		}
	}()

	<-ctx.Done()

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Stop(ctx)
	if err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}

	logg.Info("calendar is stopped")
}
