package config

import (
	"fmt"
	"log/slog"
	"math"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

type Mode string

var (
	Development Mode = "development"
	Production  Mode = "production"
)

type Config struct {
	App App
}

type App struct {
	Port       int        `env:"PORT" envDefault:"8080"`
	APIVersion string     `env:"API_VERSION" envDefault:"v1"`
	Version    string     `env:"APP_VERSION" envDefault:"v1.0.0"`
	Name       string     `env:"APP_NAME" envDefault:"Calculate"`
	Mode       Mode       `env:"APP_MODE" envDefault:"production"`
	LogLevel   slog.Level `env:"LOG_LEVEL" envDefault:"info"`
}

func MustLoad() (*Config, error) {
	var config Config

	err := cleanenv.ReadEnv(&config)
	if err != nil {
		return nil, err
	}

	if config.App.Mode != Production && config.App.Mode != Development {
		return nil, fmt.Errorf("invalid APP_MODE env value: %s", config.App.Mode)
	}

	if config.App.Port <= 0 || config.App.Port > math.MaxUint16 {
		return nil, fmt.Errorf("invalid PORT env value: %d", config.App.Port)
	}

	return &config, nil
}
