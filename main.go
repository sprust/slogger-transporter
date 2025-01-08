package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/commands"
	"slogger-transporter/internal/services/logging_service"
	"strings"
	"syscall"
)

var newApp = app.NewApp(context.Background())

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	initLogging()
}

func main() {
	args := os.Args
	argsLen := len(args)

	if argsLen == 1 {
		fmt.Println("Commands:")

		for key, command := range commands.GetCommands() {
			fmt.Printf(" %s %s - %s\n", key, command.Parameters(), command.Title())
		}

		os.Exit(0)
	}

	command, err := commands.GetCommand(args[1])

	if err != nil {
		panic(err)
	}

	newApp.AddFirstCloseListener(command)

	signals := make(chan os.Signal)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals

		slog.Warn("received stop signal")

		err = newApp.Close()

		if err != nil {
			panic(err)
		}
	}()

	err = command.Handle(&newApp, args[2:])

	if err != nil {
		panic(err)
	}

	slog.Warn("Exit")
}

func initLogging() {
	logLevels := strings.Split(newApp.GetConfig().GetLogLevels(), ",")

	var slogLevels []slog.Level

	if slices.Index(logLevels, "any") == -1 {
		for _, logLevel := range logLevels {
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

	customHandler, err := logging_service.NewCustomHandler(&newApp, logging_service.NewLevelPolicy(slogLevels))

	if err == nil {
		slog.SetDefault(slog.New(customHandler))
	} else {
		panic(err)
	}
}
