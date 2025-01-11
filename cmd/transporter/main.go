package main

import (
	"flag"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"slices"
	"slogger-transporter/internal/commands"
	"slogger-transporter/internal/config"
	"slogger-transporter/pkg/foundation/app"
	"strconv"
	"strings"
)

var args = os.Args

func init() {
	env := flag.String("env", "", "Specify the environment file to load")

	flag.Parse()

	var err error

	if env == nil || *env == "" {
		err = godotenv.Load()
	} else {
		err = godotenv.Load(*env)

		args = filterArgs()
	}

	if err != nil {
		panic(err)
	}
}

func main() {
	var commandName string
	var commandArgs []string

	if len(args) > 1 {
		commandName = args[1]
	}

	if len(args) > 2 {
		commandArgs = args[2:]
	}

	newApp := app.NewApp(
		commands.GetCommands(),
		initAppConfig(),
	)

	newApp.Start(commandName, commandArgs)
}

func initAppConfig() *app.Config {
	logKeepDays, err := strconv.Atoi(os.Getenv("LOG_KEEP_DAYS"))

	if err != nil {
		fmt.Println("LOG_KEEP_DAYS is not set or invalid, using default value 3")

		logKeepDays = 3
	}

	return app.NewConfig(
		getLogLevels(),
		os.Getenv("LOG_DIR"),
		logKeepDays,
	)
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

func filterArgs() []string {
	var result []string

	for _, arg := range args {
		if !strings.HasPrefix(arg, "--") {
			result = append(result, arg)
		}
	}

	return result
}
