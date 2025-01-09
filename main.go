package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"slices"
	"slogger-transporter/internal/commands"
	"slogger-transporter/internal/config"
	"slogger-transporter/pkg/foundation/app"
	"strings"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	args := os.Args

	var commandName string
	var commandArgs []string

	if len(args) > 1 {
		commandName = args[1]
	}

	if len(args) > 2 {
		commandArgs = args[2:]
	}

	newApp := app.NewApp(commands.GetCommands(), getLogLevels())

	newApp.Start(commandName, commandArgs)
}

func getLogLevels() []slog.Level {
	logLevels := strings.Split(config.GetConfig().GetLogLevels(), ",")

	var slogLevels []slog.Level

	if slices.Index(logLevels, "any") == -1 {
		for _, logLevel := range logLevels {
			if logLevel == "" {
				continue
			}

			switch logLevel {
			case "debug":
				slogLevels = append(slogLevels, slog.LevelDebug)
			case "info":
				slogLevels = append(slogLevels, slog.LevelInfo)
			case "warn":
				slogLevels = append(slogLevels, slog.LevelWarn)
			case "error":
				slogLevels = append(slogLevels, slog.LevelError)
			default:
				panic(fmt.Errorf("unknown log level: %s", logLevel))
			}
		}
	}

	return slogLevels
}
