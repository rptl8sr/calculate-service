package app

import (
	"calculate-service/internal/config"
	"calculate-service/internal/logger"
)

type app struct {
}

type App interface {
	Run() error
}

func MustLoad() (App, error) {
	cfg, err := config.MustLoad()
	if err != nil {
		return nil, err
	}

	logger.Init(cfg.App.LogLevel)

	if cfg.App.Mode == config.Development {
		logger.Debug("Dev mode", "config", cfg)
	}

	logger.Info("App initialized", "mode", cfg.App.Mode, "logLevel", cfg.App.LogLevel)
	return &app{}, nil
}

func (a *app) Run() error {
	return nil
}
