package app

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"slogger/pkg/foundation/commands"
	"slogger/pkg/foundation/errs"
	"slogger/pkg/foundation/logging"
	"syscall"
)

type App struct {
	commands           map[string]commands.CommandInterface
	config             *Config
	closeListeners     []io.Closer
	lastCloseListeners []io.Closer
}

func NewApp(commands map[string]commands.CommandInterface, config *Config) App {
	app := App{
		commands: commands,
		config:   config,
	}

	return app
}

func (a *App) Start(commandName string, args []string) {
	if commandName == "" {
		fmt.Println("Commands:")

		for key, command := range a.commands {
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
	customHandler, err := logging.NewCustomHandler(
		logging.NewLevelPolicy(a.config.logLevels),
		a.config.logDirPath,
		a.config.logKeepDays,
	)

	if err == nil {
		slog.SetDefault(slog.New(customHandler))
	} else {
		panic(err)
	}

	a.AddLastCloseListener(customHandler)
}
