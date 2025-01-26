package config

import "log/slog"

type Config struct {
	LogConfig LogConfig
	RmqConfig RmqConfig
}

type LogConfig struct {
	Levels   []slog.Level
	DirPath  string
	KeepDays int
}

type RmqConfig struct {
	User string
	Pass string
	Host string
	Port string
}
