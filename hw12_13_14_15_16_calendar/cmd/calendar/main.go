package main

import (
	"context"
	"flag"
	"github.com/google/uuid"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/server/http"
	str "github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage/sql"
	"log"
	"os/signal"
	"syscall"
	"time"
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
	logg := logger.New(cfg.Logger)

	// 3. Создаем хранилища
	var storage app.Storage
	if cfg.Storage.Type == "postgres" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		postgresqlStorage := sqlstorage.New(cfg.Storage.DSN)

		err := postgresqlStorage.Connect(ctx)
		if err != nil {
			logg.Error("failed to connect to database: " + err.Error())
			return
		}

		logg.Info("connected to database")

		testUuid := uuid.New().String()

		event := str.Event{
			ID:           testUuid,
			Title:        "123",
			StartTime:    time.Time{},
			EndTime:      time.Time{},
			Description:  "123",
			UserID:       1,
			NotifyBefore: 1,
		}

		err = postgresqlStorage.AddEvent(event)
		if err != nil {
			logg.Error("failed to add event: " + err.Error())
		}

		changedEvent := str.Event{
			ID:           testUuid,
			Title:        "123456",
			StartTime:    time.Time{},
			EndTime:      time.Time{},
			Description:  "123456",
			UserID:       1,
			NotifyBefore: 1,
		}

		err = postgresqlStorage.ChangeEvent(changedEvent)
		if err != nil {
			logg.Error("failed to change event: " + err.Error())
		}

		err = postgresqlStorage.RemoveEvent(testUuid)
		if err != nil {
			logg.Error("failed to remove event: " + err.Error())
		}

		defer func() {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			err := postgresqlStorage.Close(ctx)
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
	logg.Info("calendar is stopped")

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = server.Stop(ctx)
	if err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}
}
