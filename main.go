package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/commands"
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

	err = command.Handle(&newApp, args[2:])

	if err != nil {
		panic(err)
	}

	fmt.Println("Done")

	os.Exit(0)
}
