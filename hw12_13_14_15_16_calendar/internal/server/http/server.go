package internalhttp

import (
	"context"
	"errors"
	"fmt"
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

func NewServer(logger Logger, addr string, app Application) *Server {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(LoggerMiddleware)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(3 * time.Second)

		_, err := w.Write([]byte("Hello, world!"))
		if err != nil {
			return
		}
	})

	var srv *http.Server

	srv = &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{srv}
}

func (s *Server) Start() error {
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	err := s.Shutdown(ctx)
	if err != nil {
		return err
	}

	return nil
}
