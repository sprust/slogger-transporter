package app

import "log/slog"

type Config struct {
	logLevels   []slog.Level
	logDirPath  string
	logKeepDays int
}

func NewConfig(logLevels []slog.Level, logDirPath string, logKeepDays int) *Config {
	return &Config{
		logLevels:   logLevels,
		logDirPath:  logDirPath,
		logKeepDays: logKeepDays,
	}
}
