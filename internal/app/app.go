package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"slices"
	"slogger-transporter/internal/commands"
	"slogger-transporter/internal/config"
	"slogger-transporter/internal/services/errs"
	"slogger-transporter/internal/services/logging_service"
	"strings"
	"syscall"
)

type App struct {
	closeListeners     []io.Closer
	lastCloseListeners []io.Closer
	commands           map[string]commands.CommandInterface
}

func NewApp(commands map[string]commands.CommandInterface) App {
	app := App{
		commands: commands,
	}

	return app
}

func (a *App) Start(commandName string, args []string) {
	if commandName == "" {
		fmt.Println("Commands:")

		for key, command := range commands.GetCommands() {
			fmt.Printf(" %s %s - %s\n", key, command.Parameters(), command.Title())
		}

		os.Exit(0)
	}

	a.initLogging()

	command, ok := a.commands[commandName]

	if !ok {
		panic(errs.Err(errors.New("command not found")))
	}

	a.AddFirstCloseListener(command)

	signals := make(chan os.Signal)

	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signals

		slog.Warn("received stop signal")

		err := a.Close()

		if err != nil {
			panic(err)
		}
	}()

	err := command.Handle(context.Background(), args)

	if err != nil {
		panic(err)
	}

	slog.Warn("Exit")
}

func (a *App) Close() error {
	slog.Warn("Closing app...")

	for _, listener := range append(a.closeListeners, a.lastCloseListeners...) {
		err := listener.Close()

		if err != nil {
			return errs.Err(err)
		}
	}

	return nil
}

func (a *App) AddFirstCloseListener(listener io.Closer) {
	a.closeListeners = append(a.closeListeners, listener)
}

func (a *App) AddLastCloseListener(listener io.Closer) {
	a.lastCloseListeners = append(a.closeListeners, listener)
}

func (a *App) initLogging() {
	logLevels := strings.Split(config.GetConfig().GetLogLevels(), ",")

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

	customHandler, err := logging_service.NewCustomHandler(logging_service.NewLevelPolicy(slogLevels))

	if err == nil {
		slog.SetDefault(slog.New(customHandler))
	} else {
		panic(err)
	}

	a.AddLastCloseListener(customHandler)
}
