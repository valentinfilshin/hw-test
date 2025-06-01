package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/app"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/valentinfilshin/hw-test/hw12_13_14_15_calendar/internal/storage/memory"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	fmt.Println(configFile)
	// 1. Загружаем конфигурацию
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Создаем логгер
	logg := logger.New(cfg.Logger)

	// 3. Создаем хранилище по условию?
	storage := memorystorage.New()

	// 4. Бизнес-логика
	calendar := app.New(logg, storage)

	// 5. Запускаем сервер
	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
