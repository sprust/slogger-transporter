package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"os/signal"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/commands"
	"slogger-transporter/internal/services/logging_service"
	"syscall"
)

var customHandler *logging_service.CustomHandler

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	logLevel := os.Getenv("LOG_LEVEL")

	var slogLevel slog.Level

	if logLevel != "" {
		switch logLevel {
		case "any":
		case "debug":
			slogLevel = slog.LevelDebug
		case "info":
			slogLevel = slog.LevelInfo
		case "warn":
			slogLevel = slog.LevelWarn
		case "error":
			slogLevel = slog.LevelError
		default:
			panic(fmt.Errorf("unknown log level: %s", logLevel))
		}
	}

	var err error

	customHandler, err = logging_service.NewCustomHandler(&slogLevel)

	if err == nil {
		slog.SetDefault(slog.New(customHandler))
	} else {
		slog.Error(err.Error())
	}
}

func main() {
	defer func() {
		if customHandler != nil {
			err := customHandler.Close()

			if err != nil {
				panic(err)
			}
		}
	}()

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

	newApp := app.NewApp(context.Background())

	newApp.AddCloseListener(command)

	signals := make(chan os.Signal)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals

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
