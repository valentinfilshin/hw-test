package internalhttp

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	*http.Server
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Application interface { // TODO
}

func NewServer(logger Logger, app Application) *Server {
	// TODO add config params
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	// Вариант 1
	// проблема в том, что я работаю с интерфейсом логера, а не с самим slog
	// возможно стоит slog вынести глобально? чтобы не прокидывать его везде зависимостью
	// что в этом может быть плохого?

	// Вариант 2
	// можно просто на основании метода интерфейса собрать данные в структуру и вывести их,
	// но из-за этого мне везде приходится таскать логгер
	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)

		_, err := w.Write([]byte("Hello, world!"))
		if err != nil {
			return
		}
	})

	var srv *http.Server

	srv = &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{srv}
}

func (s *Server) Start(ctx context.Context) error {
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {

		}
	}()

	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.Stop(ctx)
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
