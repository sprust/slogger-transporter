package main

import (
	"github.com/joho/godotenv"
	"os"
	"slogger-transporter/internal/app"
	"slogger-transporter/internal/commands"
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

	newApp := app.NewApp(commands.GetCommands())

	newApp.Start(commandName, commandArgs)
}
