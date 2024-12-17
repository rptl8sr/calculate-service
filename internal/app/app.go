package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"calculate-service/internal/config"
	"calculate-service/internal/logger"
)

type app struct {
	server *http.Server
}

type App interface {
	Run(ctx context.Context) error
}

func MustLoad() (App, error) {
	cfg, err := config.MustLoad()
	if err != nil {
		return nil, err
	}

	logger.Init(cfg.App.LogLevel)

	if cfg.App.Mode == config.Development {
		logger.Debug("Dev mode", "config", cfg)
	} else {
		logger.Info("App initialized",
			"mode", cfg.App.Mode,
			"logLevel", cfg.App.LogLevel,
			"version", cfg.App.Version,
			"port", cfg.App.Port,
		)
	}

	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", cfg.App.Port),
	}

	return &app{server: srv}, nil
}

func (a *app) Run(ctx context.Context) error {
	go func() {
		<-ctx.Done()

		logger.Info("Shutting down server by context...")
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		err := a.server.Shutdown(ctxWithTimeout)
		if err != nil {
			logger.Error("Error shutting down server", "error", err)
			return
		}
	}()

	logger.Info("Server started", "address", a.server.Addr)
	err := a.server.ListenAndServe()

	switch {
	case errors.Is(err, http.ErrServerClosed), errors.Is(err, context.Canceled):
		logger.Info("Server shutting down gracefully")
		return nil
	default:
		logger.Error("Server error", "error", err)
		return err
	}
}
