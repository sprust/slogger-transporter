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
	"syscall"
)

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

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	command, err := commands.GetCommand(args[1])

	if err != nil {
		panic(err)
	}

	newApp := app.NewApp(context.Background())

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

	slog.Info("Exit")

	os.Exit(0)
}
