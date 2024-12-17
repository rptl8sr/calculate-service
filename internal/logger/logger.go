package logger

import (
	"log/slog"
	"os"
)

type logger struct {
	log *slog.Logger
}

var (
	globalLogger logger
)

func Init(lvl slog.Level) {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: lvl}))
	globalLogger.log = log
}

type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

func Debug(msg string, args ...interface{}) {
	globalLogger.log.Debug(msg, args...)
}

func Info(msg string, args ...interface{}) {
	globalLogger.log.Info(msg, args...)
}

func Warn(msg string, args ...interface{}) {
	globalLogger.log.Warn(msg, args...)
}

func Error(msg string, args ...interface{}) {
	globalLogger.log.Error(msg, args...)
}
