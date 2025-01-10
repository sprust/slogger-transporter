package app

import "log/slog"

type Config struct {
	logLevels  []slog.Level
	logDitPath string
}

func NewConfig(logLevels []slog.Level, logDitPath string) *Config {
	return &Config{
		logLevels:  logLevels,
		logDitPath: logDitPath,
	}
}
